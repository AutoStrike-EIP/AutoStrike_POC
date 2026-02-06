package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"autostrike/internal/application"
	"autostrike/internal/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	service        *application.AuthService
	tokenBlacklist *application.TokenBlacklist
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// NewAuthHandlerWithBlacklist creates a new auth handler with token revocation support
func NewAuthHandlerWithBlacklist(service *application.AuthService, blacklist *application.TokenBlacklist) *AuthHandler {
	return &AuthHandler{service: service, tokenBlacklist: blacklist}
}

// RegisterRoutes registers public auth routes (no auth middleware)
func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
	}
}

// RegisterRoutesWithRateLimit registers public auth routes with rate limiting
func (h *AuthHandler) RegisterRoutesWithRateLimit(r *gin.Engine, loginLimiter, refreshLimiter *middleware.RateLimiter) {
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", middleware.RateLimitMiddleware(loginLimiter), h.Login)
		auth.POST("/refresh", middleware.RateLimitMiddleware(refreshLimiter), h.Refresh)
	}
}

// RegisterProtectedRoutes registers routes that require authentication
func (h *AuthHandler) RegisterProtectedRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.GET("/me", h.Me)
	}
}

// RegisterLogoutRoute registers the logout route under an authenticated group with rate limiting.
// Logout requires authentication to prevent unauthenticated blacklist abuse.
func (h *AuthHandler) RegisterLogoutRoute(r *gin.RouterGroup, logoutLimiter *middleware.RateLimiter) {
	auth := r.Group("/auth")
	{
		auth.POST("/logout", middleware.RateLimitMiddleware(logoutLimiter), h.Logout)
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,max=72"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, application.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// RefreshRequest represents the refresh token request body
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	tokens, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, application.ErrInvalidToken) || errors.Is(err, application.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
			return
		}
		if errors.Is(err, application.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token refresh failed"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout handles user logout by revoking the current access token.
// Requires authentication — the auth middleware must set user_id in context.
func (h *AuthHandler) Logout(c *gin.Context) {
	// Require authentication to prevent unauthenticated blacklist abuse
	if _, exists := c.Get("user_id"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	if h.tokenBlacklist != nil {
		h.revokeTokenFromHeader(c.GetHeader("Authorization"))
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// revokeTokenFromHeader extracts a Bearer token, validates it, and adds it to the blacklist.
// Invalid or already-expired tokens are silently skipped.
func (h *AuthHandler) revokeTokenFromHeader(authHeader string) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return
	}

	tokenString := parts[1]

	// Validate JWT structure and extract expiry
	expiry, valid := getValidTokenExpiry(tokenString)
	if !valid {
		return // Don't revoke structurally invalid tokens
	}

	// Don't revoke already-expired tokens
	if time.Now().After(expiry) {
		return
	}

	h.tokenBlacklist.Revoke(tokenString, expiry)
}

// getValidTokenExpiry extracts and validates the exp claim from a JWT by decoding the payload.
// Returns the expiry time and true if the token has valid structure and exp claim.
// Returns zero time and false for invalid tokens — they should not be blacklisted.
func getValidTokenExpiry(tokenString string) (time.Time, bool) {
	// JWT format: header.payload.signature
	parts := strings.SplitN(tokenString, ".", 3)
	if len(parts) != 3 {
		return time.Time{}, false
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, false
	}

	var claims struct {
		Exp *float64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil || claims.Exp == nil {
		return time.Time{}, false
	}

	return time.Unix(int64(*claims.Exp), 0), true
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(c *gin.Context) {
	// User ID is set by the auth middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.service.GetCurrentUser(c.Request.Context(), userIDStr)
	if err != nil {
		if errors.Is(err, application.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
