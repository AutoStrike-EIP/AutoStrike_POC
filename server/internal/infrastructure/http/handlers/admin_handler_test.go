package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"autostrike/internal/application"
	"autostrike/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewAdminHandler(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	if handler == nil {
		t.Fatal("NewAdminHandler returned nil")
	}
}

func TestAdminHandler_ListUsers_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	repo.users["user-1"] = &entity.User{
		ID:       "user-1",
		Username: "admin",
		Email:    "admin@example.com",
		Role:     entity.RoleAdmin,
		IsActive: true,
	}
	repo.users["user-2"] = &entity.User{
		ID:       "user-2",
		Username: "viewer",
		Email:    "viewer@example.com",
		Role:     entity.RoleViewer,
		IsActive: true,
	}

	router := gin.New()
	router.GET("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ListUsers(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response ListUsersResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Total != 2 {
		t.Errorf("Expected 2 users, got %d", response.Total)
	}
}

func TestAdminHandler_ListUsers_Forbidden(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.GET("/users", func(c *gin.Context) {
		c.Set("role", "viewer") // Non-admin role
		handler.ListUsers(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAdminHandler_GetUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	repo.users["user-1"] = &entity.User{
		ID:       "user-1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     entity.RoleAdmin,
		IsActive: true,
	}

	router := gin.New()
	router.GET("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.GetUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/user-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_GetUser_NotFound(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.GET("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.GetUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/nonexistent", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestAdminHandler_CreateUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.CreateUser(c)
	})

	body := CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		Role:     "viewer",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_CreateUser_InvalidRole(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.CreateUser(c)
	})

	body := CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		Role:     "invalid_role",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAdminHandler_CreateUser_DuplicateUsername(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	repo.users["user-1"] = &entity.User{
		ID:       "user-1",
		Username: "existing",
		Email:    "existing@example.com",
		Role:     entity.RoleViewer,
		IsActive: true,
	}

	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.CreateUser(c)
	})

	body := CreateUserRequest{
		Username: "existing",
		Email:    "new@example.com",
		Password: "password123",
		Role:     "viewer",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}
}

func TestAdminHandler_UpdateUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "olduser",
		Email:        "old@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
		IsActive:     true,
	}

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{
		Username: "newusername",
		Email:    "new@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_UpdateUserRole_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
		IsActive:     true,
	}

	router := gin.New()
	router.PUT("/users/:id/role", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUserRole(c)
	})

	body := UpdateRoleRequest{Role: "operator"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1/role", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Role != "operator" {
		t.Errorf("Expected role 'operator', got %q", response.Role)
	}
}

func TestAdminHandler_DeactivateUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	repo.users["admin-1"] = &entity.User{
		ID:       "admin-1",
		Username: "admin",
		Email:    "admin@example.com",
		Role:     entity.RoleAdmin,
		IsActive: true,
	}
	repo.users["user-1"] = &entity.User{
		ID:       "user-1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     entity.RoleViewer,
		IsActive: true,
	}

	router := gin.New()
	router.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "admin-1") // Current user is admin-1
		handler.DeactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/user-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_DeactivateUser_CannotDeactivateSelf(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	repo.users["admin-1"] = &entity.User{
		ID:       "admin-1",
		Username: "admin",
		Email:    "admin@example.com",
		Role:     entity.RoleAdmin,
		IsActive: true,
	}

	router := gin.New()
	router.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "admin-1") // Trying to deactivate self
		handler.DeactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/admin-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAdminHandler_DeactivateUser_LastAdmin(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	// Only one admin user
	repo.users["admin-1"] = &entity.User{
		ID:       "admin-1",
		Username: "admin",
		Email:    "admin@example.com",
		Role:     entity.RoleAdmin,
		IsActive: true,
	}

	router := gin.New()
	router.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "other-admin") // Different admin trying to deactivate
		handler.DeactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/admin-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_ReactivateUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	repo.users["user-1"] = &entity.User{
		ID:       "user-1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     entity.RoleViewer,
		IsActive: false, // Deactivated
	}

	router := gin.New()
	router.POST("/users/:id/reactivate", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ReactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/user-1/reactivate", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAdminHandler_ResetPassword_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
		IsActive:     true,
	}

	router := gin.New()
	router.POST("/users/:id/reset-password", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ResetPassword(c)
	})

	body := ResetPasswordRequest{NewPassword: "newpassword123"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/user-1/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_ResetPassword_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users/:id/reset-password", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ResetPassword(c)
	})

	body := ResetPasswordRequest{NewPassword: "newpassword123"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/nonexistent/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestAdminHandler_RegisterRoutes(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	routes := router.Routes()

	// Check we have the expected routes by count and some key paths
	expectedCount := 8 // GET list, GET single, POST create, PUT update, PUT role, DELETE, POST reactivate, POST reset
	if len(routes) < expectedCount {
		t.Errorf("Expected at least %d routes, got %d", expectedCount, len(routes))
	}

	// Check specific routes exist
	checkRoute := func(method, path string) {
		for _, route := range routes {
			if route.Path == path && route.Method == method {
				return
			}
		}
		t.Errorf("Route %s %s not found", method, path)
	}

	checkRoute("GET", "/api/v1/admin/users")
	checkRoute("POST", "/api/v1/admin/users")
	checkRoute("GET", "/api/v1/admin/users/:id")
	checkRoute("PUT", "/api/v1/admin/users/:id")
	checkRoute("DELETE", "/api/v1/admin/users/:id")
}

func TestUserResponse_ToUserResponse(t *testing.T) {
	user := &entity.User{
		ID:       "user-1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     entity.RoleAdmin,
		IsActive: true,
	}

	resp := toUserResponse(user)

	if resp.ID != "user-1" {
		t.Errorf("ID = %q, want 'user-1'", resp.ID)
	}
	if resp.Username != "testuser" {
		t.Errorf("Username = %q, want 'testuser'", resp.Username)
	}
	if resp.Role != "admin" {
		t.Errorf("Role = %q, want 'admin'", resp.Role)
	}
	if resp.RoleDisplay != "Administrator" {
		t.Errorf("RoleDisplay = %q, want 'Administrator'", resp.RoleDisplay)
	}
}

func TestAdminHandler_NoRole(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.GET("/users", handler.ListUsers) // No role set

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAdminHandler_GetUser_Forbidden(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.GET("/users/:id", func(c *gin.Context) {
		c.Set("role", "viewer")
		handler.GetUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/user-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAdminHandler_CreateUser_Forbidden(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		c.Set("role", "viewer")
		handler.CreateUser(c)
	})

	body := CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		Role:     "viewer",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAdminHandler_CreateUser_InvalidRequest(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.CreateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAdminHandler_UpdateUser_Forbidden(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "viewer")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{Username: "newname"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAdminHandler_UpdateUser_NotFound(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{Username: "newname"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/nonexistent", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestAdminHandler_UpdateUser_InvalidRole(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	repo.users["user-1"] = &entity.User{
		ID:       "user-1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     entity.RoleViewer,
	}

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{Role: "invalid_role"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAdminHandler_UpdateUser_InvalidRequest(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAdminHandler_UpdateUserRole_Forbidden(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/users/:id/role", func(c *gin.Context) {
		c.Set("role", "viewer")
		handler.UpdateUserRole(c)
	})

	body := UpdateRoleRequest{Role: "admin"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1/role", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAdminHandler_UpdateUserRole_InvalidRequest(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/users/:id/role", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUserRole(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1/role", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAdminHandler_UpdateUserRole_InvalidRole(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/users/:id/role", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUserRole(c)
	})

	body := UpdateRoleRequest{Role: "invalid_role"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1/role", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAdminHandler_UpdateUserRole_NotFound(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/users/:id/role", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUserRole(c)
	})

	body := UpdateRoleRequest{Role: "viewer"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/nonexistent/role", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestAdminHandler_DeactivateUser_Forbidden(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("role", "viewer")
		handler.DeactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/user-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAdminHandler_DeactivateUser_NotFound(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "admin-1")
		handler.DeactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/nonexistent", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestAdminHandler_ReactivateUser_Forbidden(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users/:id/reactivate", func(c *gin.Context) {
		c.Set("role", "viewer")
		handler.ReactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/user-1/reactivate", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAdminHandler_ReactivateUser_NotFound(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users/:id/reactivate", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ReactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/nonexistent/reactivate", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestAdminHandler_ResetPassword_Forbidden(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users/:id/reset-password", func(c *gin.Context) {
		c.Set("role", "viewer")
		handler.ResetPassword(c)
	})

	body := ResetPasswordRequest{NewPassword: "newpassword123"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/user-1/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAdminHandler_ResetPassword_InvalidRequest(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users/:id/reset-password", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ResetPassword(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/user-1/reset-password", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAdminHandler_GetRoles_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.GET("/roles", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.GetRoles(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/roles", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAdminHandler_GetRoles_Forbidden(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.GET("/roles", func(c *gin.Context) {
		c.Set("role", "viewer")
		handler.GetRoles(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/roles", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestUserResponse_ToUserResponse_WithLastLogin(t *testing.T) {
	loginTime := time.Now()
	user := &entity.User{
		ID:          "user-1",
		Username:    "testuser",
		Email:       "test@example.com",
		Role:        entity.RoleAdmin,
		IsActive:    true,
		LastLoginAt: &loginTime,
	}

	resp := toUserResponse(user)

	if resp.LastLoginAt == nil {
		t.Error("LastLoginAt should not be nil")
	}
}

func TestAdminHandler_isAdmin_InvalidRoleType(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.GET("/users", func(c *gin.Context) {
		c.Set("role", 123) // Invalid type (int instead of string)
		handler.ListUsers(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestAdminHandler_UpdateUser_DuplicateUsername(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "user1",
		Email:        "user1@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
	}
	repo.users["user-2"] = &entity.User{
		ID:           "user-2",
		Username:     "user2",
		Email:        "user2@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
	}

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{Username: "user2"} // Duplicate username
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_GetUser_EmptyID(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	// Route without :id param so c.Param("id") returns ""
	router := gin.New()
	router.GET("/users/get", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.GetUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/get", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_UpdateUser_EmptyID(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	// Route without :id param so c.Param("id") returns ""
	router := gin.New()
	router.PUT("/users/update", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{Username: "newname"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/update", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_UpdateUserRole_EmptyID(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	// Route without :id param so c.Param("id") returns ""
	router := gin.New()
	router.PUT("/users/update-role", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUserRole(c)
	})

	body := UpdateRoleRequest{Role: "viewer"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/update-role", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_DeactivateUser_EmptyID(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	// Route without :id param so c.Param("id") returns ""
	router := gin.New()
	router.DELETE("/users/deactivate", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "admin-1")
		handler.DeactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/deactivate", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_ReactivateUser_EmptyID(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	// Route without :id param so c.Param("id") returns ""
	router := gin.New()
	router.POST("/users/reactivate-user", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ReactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/reactivate-user", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_ResetPassword_EmptyID(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	// Route without :id param so c.Param("id") returns ""
	router := gin.New()
	router.POST("/users/reset-pass", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ResetPassword(c)
	})

	body := ResetPasswordRequest{NewPassword: "newpassword123"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/reset-pass", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Error-returning mock for admin handler to test generic service error paths ---

// errorUserRepo is a mock user repo that can return configurable errors per operation
type errorUserRepo struct {
	users          map[string]*entity.User
	findAllErr     error
	findByIDErr    error
	findByUserErr  error
	findByEmailErr error
	createErr      error
	updateErr      error
	deactivateErr  error
	reactivateErr  error
}

func newErrorUserRepo() *errorUserRepo {
	return &errorUserRepo{users: make(map[string]*entity.User)}
}

func (m *errorUserRepo) Create(_ context.Context, user *entity.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users[user.ID] = user
	return nil
}

func (m *errorUserRepo) Update(_ context.Context, user *entity.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.users[user.ID] = user
	return nil
}

func (m *errorUserRepo) Delete(_ context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func (m *errorUserRepo) FindByID(_ context.Context, id string) (*entity.User, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	user, ok := m.users[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (m *errorUserRepo) FindByUsername(_ context.Context, username string) (*entity.User, error) {
	if m.findByUserErr != nil {
		return nil, m.findByUserErr
	}
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *errorUserRepo) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *errorUserRepo) FindAll(_ context.Context) ([]*entity.User, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	result := make([]*entity.User, 0, len(m.users))
	for _, user := range m.users {
		result = append(result, user)
	}
	return result, nil
}

func (m *errorUserRepo) FindActive(_ context.Context) ([]*entity.User, error) {
	result := make([]*entity.User, 0)
	for _, user := range m.users {
		if user.IsActive {
			result = append(result, user)
		}
	}
	return result, nil
}

func (m *errorUserRepo) UpdateLastLogin(_ context.Context, id string) error {
	return nil
}

func (m *errorUserRepo) Deactivate(_ context.Context, id string) error {
	if m.deactivateErr != nil {
		return m.deactivateErr
	}
	if user, ok := m.users[id]; ok {
		user.IsActive = false
	}
	return nil
}

func (m *errorUserRepo) Reactivate(_ context.Context, id string) error {
	if m.reactivateErr != nil {
		return m.reactivateErr
	}
	if user, ok := m.users[id]; ok {
		user.IsActive = true
	}
	return nil
}

func (m *errorUserRepo) CountByRole(_ context.Context, role entity.UserRole) (int, error) {
	count := 0
	for _, user := range m.users {
		if user.Role == role && user.IsActive {
			count++
		}
	}
	return count, nil
}

func (m *errorUserRepo) DeactivateAdminIfNotLast(_ context.Context, id string) error {
	if m.deactivateErr != nil {
		return m.deactivateErr
	}
	user, ok := m.users[id]
	if !ok {
		return errors.New("user not found")
	}
	user.IsActive = false
	return nil
}

// --- Tests for generic error paths in admin handler ---

func TestAdminHandler_ListUsers_ServiceError(t *testing.T) {
	repo := newErrorUserRepo()
	repo.findAllErr = errors.New("database connection lost")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.GET("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ListUsers(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to get users" {
		t.Errorf("Expected error 'failed to get users', got %q", response["error"])
	}
}

func TestAdminHandler_GetUser_GenericServiceError(t *testing.T) {
	repo := newErrorUserRepo()
	repo.findByIDErr = errors.New("database connection lost")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.GET("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.GetUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/user-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to get user" {
		t.Errorf("Expected error 'failed to get user', got %q", response["error"])
	}
}

func TestAdminHandler_CreateUser_GenericServiceError(t *testing.T) {
	repo := newErrorUserRepo()
	repo.createErr = errors.New("database write error")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.CreateUser(c)
	})

	body := CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		Role:     "viewer",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to create user" {
		t.Errorf("Expected error 'failed to create user', got %q", response["error"])
	}
}

func TestAdminHandler_UpdateUser_GenericGetError(t *testing.T) {
	// Tests the generic error path in UpdateUser when GetUser fails with non-ErrUserNotFound
	repo := newErrorUserRepo()
	repo.findByIDErr = errors.New("database connection lost")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{Username: "newname"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to get user" {
		t.Errorf("Expected error 'failed to get user', got %q", response["error"])
	}
}

func TestAdminHandler_UpdateUser_GenericUpdateError(t *testing.T) {
	// Tests the generic error path in UpdateUser when authService.UpdateUser fails with generic error
	repo := newErrorUserRepo()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
		IsActive:     true,
	}
	repo.updateErr = errors.New("database write error")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{Email: "newemail@example.com"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to update user" {
		t.Errorf("Expected error 'failed to update user', got %q", response["error"])
	}
}

func TestAdminHandler_UpdateUserRole_GenericServiceError(t *testing.T) {
	// Tests generic error path in UpdateUserRole when service returns non-ErrUserNotFound
	repo := newErrorUserRepo()
	repo.findByIDErr = errors.New("database connection lost")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/users/:id/role", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUserRole(c)
	})

	body := UpdateRoleRequest{Role: "viewer"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1/role", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to update user role" {
		t.Errorf("Expected error 'failed to update user role', got %q", response["error"])
	}
}

func TestAdminHandler_DeactivateUser_GenericServiceError(t *testing.T) {
	// Tests generic error path in DeactivateUser when service returns unknown error
	repo := newErrorUserRepo()
	repo.users["user-1"] = &entity.User{
		ID:       "user-1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     entity.RoleViewer,
		IsActive: true,
	}
	repo.deactivateErr = errors.New("database connection lost")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "admin-1") // Different user
		handler.DeactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/user-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to deactivate user" {
		t.Errorf("Expected error 'failed to deactivate user', got %q", response["error"])
	}
}

func TestAdminHandler_ReactivateUser_GenericServiceError(t *testing.T) {
	// Tests generic error path in ReactivateUser when service returns non-ErrUserNotFound
	repo := newErrorUserRepo()
	repo.findByIDErr = errors.New("database connection lost")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users/:id/reactivate", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ReactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/user-1/reactivate", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to reactivate user" {
		t.Errorf("Expected error 'failed to reactivate user', got %q", response["error"])
	}
}

func TestAdminHandler_ResetPassword_GenericServiceError(t *testing.T) {
	// Tests generic error path in ResetPassword when service returns non-ErrUserNotFound
	repo := newErrorUserRepo()
	repo.findByIDErr = errors.New("database connection lost")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users/:id/reset-password", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ResetPassword(c)
	})

	body := ResetPasswordRequest{NewPassword: "newpassword123"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/user-1/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to reset password" {
		t.Errorf("Expected error 'failed to reset password', got %q", response["error"])
	}
}

func TestAdminHandler_CreateUser_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	repo.users["user-1"] = &entity.User{
		ID:       "user-1",
		Username: "existinguser",
		Email:    "existing@example.com",
		Role:     entity.RoleViewer,
		IsActive: true,
	}

	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.CreateUser(c)
	})

	body := CreateUserRequest{
		Username: "newuser",
		Email:    "existing@example.com", // Duplicate email
		Password: "password123",
		Role:     "viewer",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate email, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_CreateUser_ShortPassword(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.CreateUser(c)
	})

	body := CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "short", // Less than 8 chars
		Role:     "viewer",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for short password, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_CreateUser_ShortUsername(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.CreateUser(c)
	})

	body := CreateUserRequest{
		Username: "ab", // Less than 3 chars
		Email:    "new@example.com",
		Password: "password123",
		Role:     "viewer",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for short username, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_CreateUser_InvalidEmail(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.CreateUser(c)
	})

	body := CreateUserRequest{
		Username: "newuser",
		Email:    "not-an-email",
		Password: "password123",
		Role:     "viewer",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid email, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_UpdateUser_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "user1",
		Email:        "user1@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
	}
	repo.users["user-2"] = &entity.User{
		ID:           "user-2",
		Username:     "user2",
		Email:        "user2@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
	}

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{Email: "user2@example.com"} // Duplicate email
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate email, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_DeactivateUser_NoUserIDInContext(t *testing.T) {
	// Tests the edge case where user_id is not set in context for DeactivateUser
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	repo.users["user-1"] = &entity.User{
		ID:       "user-1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     entity.RoleViewer,
		IsActive: true,
	}

	router := gin.New()
	router.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		// No user_id set - currentUserID will be empty string
		handler.DeactivateUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/user-1", nil)
	router.ServeHTTP(w, req)

	// Should succeed because empty string != "user-1"
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_ResetPassword_WeakPassword(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
		IsActive:     true,
	}

	router := gin.New()
	router.POST("/users/:id/reset-password", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.ResetPassword(c)
	})

	// Password less than 8 characters should fail validation
	body := ResetPasswordRequest{NewPassword: "short"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/user-1/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for weak password, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_UpdateUser_InvalidEmailFormat(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
	}

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{Email: "not-an-email"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid email format, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_UpdateUser_WithValidRole(t *testing.T) {
	// Tests the mergeUserFields path where a valid non-empty role is provided
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
		IsActive:     true,
	}

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	// Provide a valid role in the update request to exercise the mergeUserFields valid-role path
	body := UpdateUserRequest{Role: "operator"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response.Role != "operator" {
		t.Errorf("Expected role 'operator', got %q", response.Role)
	}
}

func TestAdminHandler_UpdateUser_ShortUsername(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAdminHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleViewer,
	}

	router := gin.New()
	router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("role", "admin")
		handler.UpdateUser(c)
	})

	body := UpdateUserRequest{Username: "ab"} // Less than 3 chars
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for short username, got %d: %s", w.Code, w.Body.String())
	}
}
