package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"autostrike/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestLoggingMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name       string
		statusCode int
	}{
		{"success request", http.StatusOK},
		{"client error", http.StatusBadRequest},
		{"server error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(LoggingMiddleware(logger))
			router.GET("/test", func(c *gin.Context) {
				c.Status(tt.statusCode)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test?foo=bar", nil)
			req.Header.Set("User-Agent", "test-agent")
			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestLoggingMiddleware_WithErrors(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	router := gin.New()
	router.Use(LoggingMiddleware(logger))
	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(gin.Error{Err: http.ErrAbortHandler, Type: gin.ErrorTypePrivate})
		c.Status(http.StatusBadRequest)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	router := gin.New()
	router.Use(RecoveryMiddleware(logger))
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestRecoveryMiddleware_NoPanic(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	router := gin.New()
	router.Use(RecoveryMiddleware(logger))
	router.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ok", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_NoHeader(t *testing.T) {
	config := &AuthConfig{JWTSecret: "secret"}

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	config := &AuthConfig{JWTSecret: "secret"}

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	config := &AuthConfig{JWTSecret: "secret"}

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secret := "test-secret-key"
	config := &AuthConfig{JWTSecret: secret}

	// Create valid access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"role": "admin",
		"type": "access", // Must be access token
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")
		if userID != "user123" || role != "admin" {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_RefreshTokenRejected(t *testing.T) {
	secret := "test-secret-key"
	config := &AuthConfig{JWTSecret: secret}

	// Create refresh token (should be rejected for API auth)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"type": "refresh",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for refresh token, got %d", w.Code)
	}
}

func TestAuthMiddleware_TokenWithoutType(t *testing.T) {
	secret := "test-secret-key"
	config := &AuthConfig{JWTSecret: secret}

	// Create token without type claim (should be rejected)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"role": "admin",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for token without type, got %d", w.Code)
	}
}

func TestAuthMiddleware_WrongSigningMethod(t *testing.T) {
	config := &AuthConfig{JWTSecret: "secret"}

	// Create token with wrong signing method
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "user123",
	})
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAgentAuthMiddleware_NoKey(t *testing.T) {
	config := &AuthConfig{AgentSecret: "agent-secret"}

	router := gin.New()
	router.Use(AgentAuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAgentAuthMiddleware_InvalidKey(t *testing.T) {
	config := &AuthConfig{AgentSecret: "agent-secret"}

	router := gin.New()
	router.Use(AgentAuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Agent-Key", "wrong-key")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAgentAuthMiddleware_ValidKey(t *testing.T) {
	config := &AuthConfig{AgentSecret: "agent-secret"}

	router := gin.New()
	router.Use(AgentAuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		isAgent, _ := c.Get("is_agent")
		if isAgent != true {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Agent-Key", "agent-secret")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRoleMiddleware_NoRole(t *testing.T) {
	router := gin.New()
	router.Use(RoleMiddleware("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestRoleMiddleware_InvalidRoleFormat(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", 123) // Not a string
		c.Next()
	})
	router.Use(RoleMiddleware("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestRoleMiddleware_InsufficientPermissions(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})
	router.Use(RoleMiddleware("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestRoleMiddleware_HasRole(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	})
	router.Use(RoleMiddleware("admin", "superadmin"))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthConfig_Struct(t *testing.T) {
	config := &AuthConfig{
		JWTSecret:   "jwt-secret",
		AgentSecret: "agent-secret",
	}

	if config.JWTSecret != "jwt-secret" {
		t.Errorf("JWTSecret = %s, want jwt-secret", config.JWTSecret)
	}
	if config.AgentSecret != "agent-secret" {
		t.Errorf("AgentSecret = %s, want agent-secret", config.AgentSecret)
	}
}

func TestPermissionMiddleware_NoRole(t *testing.T) {
	router := gin.New()
	router.Use(PermissionMiddleware(entity.PermissionUsersView))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestPermissionMiddleware_InvalidRoleFormat(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", 123) // Not a string
		c.Next()
	})
	router.Use(PermissionMiddleware(entity.PermissionUsersView))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestPermissionMiddleware_AdminHasAllPermissions(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	})
	router.Use(PermissionMiddleware(entity.PermissionUsersView, entity.PermissionUsersCreate))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestPermissionMiddleware_ViewerLacksPermission(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "viewer")
		c.Next()
	})
	router.Use(PermissionMiddleware(entity.PermissionUsersCreate))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestPermissionMiddleware_ViewerHasViewPermission(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "viewer")
		c.Next()
	})
	router.Use(PermissionMiddleware(entity.PermissionAgentsView))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestPermissionMiddleware_OperatorCanExecute(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "operator")
		c.Next()
	})
	router.Use(PermissionMiddleware(entity.PermissionExecutionsStart))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRequireAnyPermission_NoRole(t *testing.T) {
	router := gin.New()
	router.Use(RequireAnyPermission(entity.PermissionUsersView, entity.PermissionAgentsView))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestRequireAnyPermission_InvalidRoleFormat(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", 123)
		c.Next()
	})
	router.Use(RequireAnyPermission(entity.PermissionUsersView))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestRequireAnyPermission_HasOnePermission(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "viewer")
		c.Next()
	})
	// Viewer has agents:view but not users:view
	router.Use(RequireAnyPermission(entity.PermissionUsersView, entity.PermissionAgentsView))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRequireAnyPermission_HasNoPermissions(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "viewer")
		c.Next()
	})
	// Viewer has neither users:create nor analytics:compare
	router.Use(RequireAnyPermission(entity.PermissionUsersCreate, entity.PermissionAnalyticsCompare))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestRequireAnyPermission_AdminHasAll(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	})
	router.Use(RequireAnyPermission(entity.PermissionUsersCreate, entity.PermissionAnalyticsCompare))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNoAuthMiddleware_SetsDefaultContext(t *testing.T) {
	router := gin.New()
	router.Use(NoAuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")
		if userID != "anonymous" {
			c.Status(http.StatusBadRequest)
			return
		}
		if role != "admin" {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_LowercaseBearer(t *testing.T) {
	secret := "test-secret-key"
	config := &AuthConfig{JWTSecret: secret}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"role": "admin",
		"type": "access",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "bearer "+tokenString) // lowercase "bearer"
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for lowercase 'bearer', got %d", w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	secret := "test-secret-key"
	config := &AuthConfig{JWTSecret: secret}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"role": "admin",
		"type": "access",
		"exp":  time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	})
	tokenString, _ := token.SignedString([]byte(secret))

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for expired token, got %d", w.Code)
	}
}

func TestAuthMiddleware_OnlyTokenPrefix(t *testing.T) {
	config := &AuthConfig{JWTSecret: "secret"}

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Token abc123") // Wrong prefix
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for 'Token' prefix, got %d", w.Code)
	}
}

// --- Rate Limiter Tests ---

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(5, time.Minute)
	defer rl.Close()
	if rl == nil {
		t.Fatal("Expected non-nil rate limiter")
	}
	if rl.limit != 5 {
		t.Errorf("Expected limit 5, got %d", rl.limit)
	}
	if rl.window != time.Minute {
		t.Errorf("Expected window 1m, got %v", rl.window)
	}
}

func TestRateLimiter_Close_Idempotent(t *testing.T) {
	rl := NewRateLimiter(5, time.Minute)
	rl.Close()
	rl.Close() // must not panic
}

func TestRateLimiter_Allow_WithinLimit(t *testing.T) {
	rl := &RateLimiter{
		ips:    make(map[string]*ipEntry),
		limit:  3,
		window: time.Minute,
	}

	for i := 0; i < 3; i++ {
		if !rl.allow("192.168.1.1") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}
}

func TestRateLimiter_Allow_ExceedsLimit(t *testing.T) {
	rl := &RateLimiter{
		ips:    make(map[string]*ipEntry),
		limit:  2,
		window: time.Minute,
	}

	rl.allow("192.168.1.1")
	rl.allow("192.168.1.1")

	if rl.allow("192.168.1.1") {
		t.Error("Third request should be denied")
	}
}

func TestRateLimiter_Allow_DifferentIPs(t *testing.T) {
	rl := &RateLimiter{
		ips:    make(map[string]*ipEntry),
		limit:  1,
		window: time.Minute,
	}

	if !rl.allow("192.168.1.1") {
		t.Error("First IP should be allowed")
	}
	if !rl.allow("192.168.1.2") {
		t.Error("Second IP should be allowed (different IP)")
	}
}

func TestRateLimiter_Allow_WindowReset(t *testing.T) {
	rl := &RateLimiter{
		ips:    make(map[string]*ipEntry),
		limit:  1,
		window: time.Millisecond,
	}

	rl.allow("192.168.1.1")

	// Wait for window to expire
	time.Sleep(5 * time.Millisecond)

	if !rl.allow("192.168.1.1") {
		t.Error("Request should be allowed after window reset")
	}
}

func TestRateLimiter_Cleanup(t *testing.T) {
	rl := &RateLimiter{
		ips:    make(map[string]*ipEntry),
		limit:  5,
		window: time.Minute,
	}

	// Add expired entry
	rl.ips["expired"] = &ipEntry{count: 3, resetAt: time.Now().Add(-time.Hour)}
	// Add active entry
	rl.ips["active"] = &ipEntry{count: 1, resetAt: time.Now().Add(time.Hour)}

	rl.cleanup()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, exists := rl.ips["expired"]; exists {
		t.Error("Expected expired entry to be cleaned up")
	}
	if _, exists := rl.ips["active"]; !exists {
		t.Error("Expected active entry to remain")
	}
}

func TestRateLimiter_Cleanup_EmptyMap(t *testing.T) {
	rl := &RateLimiter{
		ips:    make(map[string]*ipEntry),
		limit:  5,
		window: time.Minute,
	}
	rl.cleanup() // Should not panic
}

func TestRateLimitMiddleware_Allowed(t *testing.T) {
	rl := &RateLimiter{
		ips:    make(map[string]*ipEntry),
		limit:  5,
		window: time.Minute,
	}

	router := gin.New()
	router.Use(RateLimitMiddleware(rl))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRateLimitMiddleware_Blocked(t *testing.T) {
	rl := &RateLimiter{
		ips:    make(map[string]*ipEntry),
		limit:  1,
		window: time.Minute,
	}

	router := gin.New()
	router.Use(RateLimitMiddleware(rl))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// First request - allowed
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("First request: expected status 200, got %d", w.Code)
	}

	// Second request - blocked
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Second request: expected status 429, got %d", w.Code)
	}
}

// --- Security Headers Tests ---

func TestSecurityHeadersMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(SecurityHeadersMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedHeaders := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":          "DENY",
		"X-XSS-Protection":         "1; mode=block",
		"Referrer-Policy":          "strict-origin-when-cross-origin",
		"Permissions-Policy":       "camera=(), microphone=(), geolocation=()",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}

	for header, expected := range expectedHeaders {
		actual := w.Header().Get(header)
		if actual != expected {
			t.Errorf("Header %s: expected %q, got %q", header, expected, actual)
		}
	}

	// Check CSP contains key directives
	csp := w.Header().Get("Content-Security-Policy")
	if csp == "" {
		t.Error("Expected Content-Security-Policy header to be set")
	}
	for _, directive := range []string{"default-src 'self'", "script-src", "style-src", "connect-src"} {
		if !strings.Contains(csp, directive) {
			t.Errorf("CSP missing directive: %s", directive)
		}
	}
}

// --- Body Size Limit Middleware Tests ---

func TestBodySizeLimitMiddleware_SmallBody(t *testing.T) {
	router := gin.New()
	router.Use(BodySizeLimitMiddleware(1024))
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	body := strings.NewReader(`{"key":"value"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", body)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for small body, got %d", w.Code)
	}
}

func TestBodySizeLimitMiddleware_LargeContentLength(t *testing.T) {
	router := gin.New()
	router.Use(BodySizeLimitMiddleware(1024))
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	body := strings.NewReader(`{"key":"value"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", body)
	req.ContentLength = 2048 // Declared larger than limit
	router.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected status 413 for large Content-Length, got %d", w.Code)
	}
}

func TestBodySizeLimitMiddleware_NilBody(t *testing.T) {
	router := gin.New()
	router.Use(BodySizeLimitMiddleware(1024))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for nil body, got %d", w.Code)
	}
}

func TestBodySizeLimitMiddleware_ExactLimit(t *testing.T) {
	router := gin.New()
	router.Use(BodySizeLimitMiddleware(10))
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	body := strings.NewReader("0123456789") // exactly 10 bytes
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", body)
	req.Header.Set("Content-Type", "application/octet-stream")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for body at exact limit, got %d", w.Code)
	}
}

func TestBodySizeLimitMiddleware_StreamingOverflow(t *testing.T) {
	router := gin.New()
	router.Use(BodySizeLimitMiddleware(16))
	router.POST("/test", func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	// Body is 30+ bytes but Content-Length is not set (simulating chunked)
	largeBody := `{"key":"` + strings.Repeat("x", 100) + `"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(largeBody))
	req.ContentLength = -1 // Unknown length
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// MaxBytesReader should cause ShouldBindJSON to fail
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for streaming overflow, got %d", w.Code)
	}
}

// --- Auth Middleware Blacklist Tests ---

type mockBlacklist struct {
	revoked map[string]bool
}

func (m *mockBlacklist) IsRevoked(token string) bool {
	return m.revoked[token]
}

func TestAuthMiddleware_RevokedToken(t *testing.T) {
	secret := "test-secret-key"
	bl := &mockBlacklist{revoked: make(map[string]bool)}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"role": "admin",
		"type": "access",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	// Mark token as revoked
	bl.revoked[tokenString] = true

	config := &AuthConfig{
		JWTSecret:      secret,
		TokenBlacklist: bl,
	}

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for revoked token, got %d", w.Code)
	}
}

func TestAuthMiddleware_NonRevokedTokenWithBlacklist(t *testing.T) {
	secret := "test-secret-key"
	bl := &mockBlacklist{revoked: make(map[string]bool)}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"role": "admin",
		"type": "access",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	config := &AuthConfig{
		JWTSecret:      secret,
		TokenBlacklist: bl,
	}

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for non-revoked token, got %d", w.Code)
	}
}

func TestAuthMiddleware_NilBlacklist(t *testing.T) {
	secret := "test-secret-key"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"role": "admin",
		"type": "access",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	config := &AuthConfig{
		JWTSecret:      secret,
		TokenBlacklist: nil, // No blacklist
	}

	router := gin.New()
	router.Use(AuthMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 with nil blacklist, got %d", w.Code)
	}
}
