package middleware

import (
	"net/http"
	"strings"

	"autostrike/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Error messages
const (
	errRoleNotFound            = "role not found"
	errInvalidRoleFormat       = "invalid role format"
	errInsufficientPermissions = "insufficient permissions"
)

// TokenBlacklistChecker checks if a token has been revoked
type TokenBlacklistChecker interface {
	IsRevoked(token string) bool
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	JWTSecret      string
	AgentSecret    string
	TokenBlacklist TokenBlacklistChecker
}

// NoAuthMiddleware creates a middleware that sets default user context when auth is disabled
// This allows handlers that check for user_id to work in development mode
func NoAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set default anonymous user with admin role for full access
		c.Set("user_id", "anonymous")
		c.Set("role", "admin")
		c.Next()
	}
}

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(config *AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		claims, validationErr := validateAccessToken(tokenString, config)
		if validationErr != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": validationErr})
			c.Abort()
			return
		}

		c.Set("user_id", claims["sub"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

// extractBearerToken extracts the token from an Authorization header.
// Returns the token string and an empty error, or empty token and error message.
func extractBearerToken(authHeader string) (string, string) {
	if authHeader == "" {
		return "", "authorization header required"
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", "invalid authorization header format"
	}
	return parts[1], ""
}

// validateAccessToken parses, validates, and checks revocation of a JWT access token.
// Returns the claims on success or an error message string.
func validateAccessToken(tokenString string, config *AuthConfig) (jwt.MapClaims, string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, "invalid token"
	}

	if config.TokenBlacklist != nil && config.TokenBlacklist.IsRevoked(tokenString) {
		return nil, "token has been revoked"
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, "invalid token claims"
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "access" {
		return nil, "invalid token type"
	}

	return claims, ""
}

// AgentAuthMiddleware creates an agent authentication middleware
func AgentAuthMiddleware(config *AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		agentKey := c.GetHeader("X-Agent-Key")

		if agentKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "agent key required"})
			c.Abort()
			return
		}

		// For mTLS, the certificate validation happens at TLS layer
		// This is a secondary check using a pre-shared key
		if agentKey != config.AgentSecret {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid agent key"})
			c.Abort()
			return
		}

		c.Set("is_agent", true)
		c.Next()
	}
}

// RoleMiddleware creates a role-based authorization middleware
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": errRoleNotFound})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": errInvalidRoleFormat})
			c.Abort()
			return
		}

		for _, r := range requiredRoles {
			if r == roleStr {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": errInsufficientPermissions})
		c.Abort()
	}
}

// PermissionMiddleware creates a permission-based authorization middleware
func PermissionMiddleware(requiredPermissions ...entity.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": errRoleNotFound})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": errInvalidRoleFormat})
			c.Abort()
			return
		}

		userRole := entity.UserRole(roleStr)

		// Check if user has ALL required permissions
		for _, perm := range requiredPermissions {
			if !entity.HasPermission(userRole, perm) {
				c.JSON(http.StatusForbidden, gin.H{
					"error":     errInsufficientPermissions,
					"required":  string(perm),
					"user_role": roleStr,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireAnyPermission creates a middleware that requires at least one of the given permissions
func RequireAnyPermission(permissions ...entity.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": errRoleNotFound})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": errInvalidRoleFormat})
			c.Abort()
			return
		}

		userRole := entity.UserRole(roleStr)

		// Check if user has ANY of the required permissions
		for _, perm := range permissions {
			if entity.HasPermission(userRole, perm) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":     errInsufficientPermissions,
			"user_role": roleStr,
		})
		c.Abort()
	}
}
