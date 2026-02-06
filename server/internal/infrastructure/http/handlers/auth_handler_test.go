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
	"autostrike/internal/infrastructure/http/middleware"
	"autostrike/internal/infrastructure/persistence/sqlite"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo implements repository.UserRepository for testing
type mockUserRepo struct {
	users     map[string]*entity.User
	findErr   error
	createErr error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*entity.User)}
}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	user, ok := m.users[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockUserRepo) FindAll(ctx context.Context) ([]*entity.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	result := make([]*entity.User, 0, len(m.users))
	for _, user := range m.users {
		result = append(result, user)
	}
	return result, nil
}

func (m *mockUserRepo) FindActive(ctx context.Context) ([]*entity.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	result := make([]*entity.User, 0)
	for _, user := range m.users {
		if user.IsActive {
			result = append(result, user)
		}
	}
	return result, nil
}

func (m *mockUserRepo) UpdateLastLogin(ctx context.Context, id string) error {
	user, ok := m.users[id]
	if !ok {
		return sql.ErrNoRows
	}
	now := time.Now()
	user.LastLoginAt = &now
	return nil
}

func (m *mockUserRepo) Deactivate(ctx context.Context, id string) error {
	user, ok := m.users[id]
	if !ok {
		return sql.ErrNoRows
	}
	user.IsActive = false
	return nil
}

func (m *mockUserRepo) Reactivate(ctx context.Context, id string) error {
	user, ok := m.users[id]
	if !ok {
		return sql.ErrNoRows
	}
	user.IsActive = true
	return nil
}

func (m *mockUserRepo) CountByRole(ctx context.Context, role entity.UserRole) (int, error) {
	count := 0
	for _, user := range m.users {
		if user.Role == role && user.IsActive {
			count++
		}
	}
	return count, nil
}

func (m *mockUserRepo) DeactivateAdminIfNotLast(ctx context.Context, id string) error {
	user, ok := m.users[id]
	if !ok {
		return sqlite.ErrUserNotFound
	}
	// If admin, check if last admin
	if user.Role == entity.RoleAdmin {
		count := 0
		for _, u := range m.users {
			if u.Role == entity.RoleAdmin && u.IsActive && u.ID != id {
				count++
			}
		}
		if count == 0 {
			return sqlite.ErrLastAdmin
		}
	}
	user.IsActive = false
	return nil
}

func TestNewAuthHandler(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	if handler == nil {
		t.Fatal("NewAuthHandler returned nil")
	}
}

func TestAuthHandler_RegisterRoutes(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	handler.RegisterRoutes(router)

	// Check public routes are registered (logout is no longer public)
	routes := router.Routes()
	expectedPaths := map[string]string{
		"/api/v1/auth/login":   "POST",
		"/api/v1/auth/refresh": "POST",
	}

	for path, method := range expectedPaths {
		found := false
		for _, route := range routes {
			if route.Path == path && route.Method == method {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Route %s %s not found", method, path)
		}
	}

	// Verify logout is NOT in public routes
	for _, route := range routes {
		if route.Path == "/api/v1/auth/logout" {
			t.Error("Logout should not be in public routes")
		}
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	// Create user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleAdmin,
		IsActive:     true,
	}

	router := gin.New()
	router.POST("/login", handler.Login)

	body := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response application.TokenResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("AccessToken should not be empty")
	}
	if response.RefreshToken == "" {
		t.Error("RefreshToken should not be empty")
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}

	router := gin.New()
	router.POST("/login", handler.Login)

	body := LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/login", handler.Login)

	body := LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Login_BadRequest(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/login", handler.Login)

	// Missing password
	body := `{"username": "testuser"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/login", handler.Login)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Refresh_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleAdmin,
		IsActive:     true,
	}

	// Login first to get tokens
	tokens, _ := service.Login(context.Background(), "testuser", "password123")

	router := gin.New()
	router.POST("/refresh", handler.Refresh)

	body := RefreshRequest{
		RefreshToken: tokens.RefreshToken,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_Refresh_InvalidToken(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/refresh", handler.Refresh)

	body := RefreshRequest{
		RefreshToken: "invalid-token",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Refresh_BadRequest(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/refresh", handler.Refresh)

	// Missing refresh_token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", "user-1") // Simulate auth middleware
		handler.Logout(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthHandler_Logout_Unauthenticated(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/logout", handler.Logout)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for unauthenticated logout, got %d", w.Code)
	}
}

func TestAuthHandler_Me_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	repo.users["user-1"] = &entity.User{
		ID:       "user-1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     entity.RoleAdmin,
	}

	router := gin.New()
	router.GET("/me", func(c *gin.Context) {
		// Simulate auth middleware setting user_id
		c.Set("user_id", "user-1")
		handler.Me(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var user entity.User
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Username = %q, want 'testuser'", user.Username)
	}
}

func TestAuthHandler_Me_NotAuthenticated(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.GET("/me", handler.Me)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Me_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.GET("/me", func(c *gin.Context) {
		c.Set("user_id", "nonexistent")
		handler.Me(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestAuthHandler_Me_InvalidUserID(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.GET("/me", func(c *gin.Context) {
		// Set invalid type for user_id
		c.Set("user_id", 123)
		handler.Me(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestAuthHandler_RegisterProtectedRoutes(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterProtectedRoutes(api)

	routes := router.Routes()
	found := false
	for _, route := range routes {
		if route.Path == "/api/v1/auth/me" && route.Method == "GET" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Route GET /api/v1/auth/me not found")
	}
}

func TestLoginRequest_Struct(t *testing.T) {
	req := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	if req.Username != "testuser" {
		t.Errorf("Username = %q, want 'testuser'", req.Username)
	}
	if req.Password != "password123" {
		t.Errorf("Password = %q, want 'password123'", req.Password)
	}
}

func TestRefreshRequest_Struct(t *testing.T) {
	req := RefreshRequest{
		RefreshToken: "some-token",
	}

	if req.RefreshToken != "some-token" {
		t.Errorf("RefreshToken = %q, want 'some-token'", req.RefreshToken)
	}
}

// --- Edge case tests ---

func TestAuthHandler_Login_MalformedJSON(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/login", handler.Login)

	// Completely malformed JSON (truncated)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(`{"username": "test", "password`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for malformed JSON, got %d", w.Code)
	}
}

func TestAuthHandler_Login_EmptyBody(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/login", handler.Login)

	// Empty body - no JSON at all
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty body, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "username and password are required" {
		t.Errorf("Expected error 'username and password are required', got %q", response["error"])
	}
}

func TestAuthHandler_Login_InternalServerError(t *testing.T) {
	repo := newMockUserRepo()
	// Set a generic find error that is NOT sql.ErrNoRows
	repo.findErr = errors.New("database connection lost")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/login", handler.Login)

	body := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 for generic repo error, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "authentication failed" {
		t.Errorf("Expected error 'authentication failed', got %q", response["error"])
	}
}

func TestAuthHandler_Login_InactiveUser(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	// Create an inactive user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "inactiveuser",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleOperator,
		IsActive:     false,
	}

	router := gin.New()
	router.POST("/login", handler.Login)

	body := LoginRequest{
		Username: "inactiveuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// ErrUserInactive is not ErrInvalidCredentials, so it falls through to 500
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 for inactive user, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "authentication failed" {
		t.Errorf("Expected error 'authentication failed', got %q", response["error"])
	}
}

func TestAuthHandler_Login_MissingUsername(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/login", handler.Login)

	// Has password but missing username
	body := `{"password": "password123"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing username, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "username and password are required" {
		t.Errorf("Expected error 'username and password are required', got %q", response["error"])
	}
}

func TestAuthHandler_Refresh_MalformedJSON(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/refresh", handler.Refresh)

	// Completely malformed JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBufferString(`not json at all`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for malformed JSON, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "refresh_token is required" {
		t.Errorf("Expected error 'refresh_token is required', got %q", response["error"])
	}
}

func TestAuthHandler_Refresh_EmptyBody(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/refresh", handler.Refresh)

	// Empty body
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty body, got %d", w.Code)
	}
}

func TestAuthHandler_Refresh_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	// Create user, login to get tokens, then delete the user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleAdmin,
		IsActive:     true,
	}

	tokens, err := service.Login(context.Background(), "testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Remove the user so Refresh will find a valid token but no user
	delete(repo.users, "user-1")

	router := gin.New()
	router.POST("/refresh", handler.Refresh)

	body := RefreshRequest{
		RefreshToken: tokens.RefreshToken,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for user not found during refresh, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "user not found" {
		t.Errorf("Expected error 'user not found', got %q", response["error"])
	}
}

func TestAuthHandler_Refresh_InternalServerError(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	// Create user, login to get tokens, then set a generic repo error
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	repo.users["user-1"] = &entity.User{
		ID:           "user-1",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleAdmin,
		IsActive:     true,
	}

	tokens, err := service.Login(context.Background(), "testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Set a generic error on FindByID so refresh fails with a non-specific error
	repo.findErr = errors.New("database connection lost")

	router := gin.New()
	router.POST("/refresh", handler.Refresh)

	body := RefreshRequest{
		RefreshToken: tokens.RefreshToken,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 for generic repo error during refresh, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "token refresh failed" {
		t.Errorf("Expected error 'token refresh failed', got %q", response["error"])
	}
}

func TestAuthHandler_Me_RepoError(t *testing.T) {
	repo := newMockUserRepo()
	// Set a generic error that is NOT sql.ErrNoRows
	repo.findErr = errors.New("database connection lost")
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.GET("/me", func(c *gin.Context) {
		c.Set("user_id", "user-1")
		handler.Me(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 for generic repo error, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "failed to get user" {
		t.Errorf("Expected error 'failed to get user', got %q", response["error"])
	}
}

func TestAuthHandler_Me_EmptyStringUserID(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.GET("/me", func(c *gin.Context) {
		// Set user_id to an empty string - it exists and is a string,
		// but no user will match it
		c.Set("user_id", "")
		handler.Me(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	router.ServeHTTP(w, req)

	// Empty string user_id passes the type assertion but FindByID returns sql.ErrNoRows
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for empty string user_id, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "user not found" {
		t.Errorf("Expected error 'user not found', got %q", response["error"])
	}
}

func TestAuthHandler_Me_NilUserID(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.GET("/me", func(c *gin.Context) {
		// Set user_id to nil - c.Get returns (nil, true), but type assertion to string fails
		c.Set("user_id", nil)
		handler.Me(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 for nil user_id, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "invalid user id" {
		t.Errorf("Expected error 'invalid user id', got %q", response["error"])
	}
}

func TestAuthHandler_Me_BoolUserID(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.GET("/me", func(c *gin.Context) {
		// Set user_id to a boolean - type assertion to string fails
		c.Set("user_id", true)
		handler.Me(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 for bool user_id, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["error"] != "invalid user id" {
		t.Errorf("Expected error 'invalid user id', got %q", response["error"])
	}
}

func TestAuthHandler_Logout_ResponseBody(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", "user-1")
		handler.Logout(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["message"] != "logged out successfully" {
		t.Errorf("Expected message 'logged out successfully', got %q", response["message"])
	}
}

// --- Token Blacklist / Revocation Tests ---

func TestNewAuthHandlerWithBlacklist(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	bl := application.NewTokenBlacklist()

	handler := NewAuthHandlerWithBlacklist(service, bl)
	if handler == nil {
		t.Fatal("Expected non-nil handler")
	}
}

func TestAuthHandler_RegisterRoutesWithRateLimit(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	loginLimiter := middleware.NewRateLimiter(5, time.Minute)
	refreshLimiter := middleware.NewRateLimiter(10, time.Minute)

	router := gin.New()
	handler.RegisterRoutesWithRateLimit(router, loginLimiter, refreshLimiter)

	routes := router.Routes()
	expectedPaths := map[string]string{
		"/api/v1/auth/login":   "POST",
		"/api/v1/auth/refresh": "POST",
	}

	for path, method := range expectedPaths {
		found := false
		for _, route := range routes {
			if route.Path == path && route.Method == method {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Route %s %s not found", method, path)
		}
	}

	// Verify logout is NOT in public rate-limited routes
	for _, route := range routes {
		if route.Path == "/api/v1/auth/logout" {
			t.Error("Logout should not be in public rate-limited routes")
		}
	}
}

func TestAuthHandler_RegisterLogoutRoute(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	bl := application.NewTokenBlacklist()
	handler := NewAuthHandlerWithBlacklist(service, bl)

	logoutLimiter := middleware.NewRateLimiter(10, time.Minute)

	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterLogoutRoute(api, logoutLimiter)

	routes := router.Routes()
	found := false
	for _, route := range routes {
		if route.Path == "/api/v1/auth/logout" && route.Method == "POST" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Route POST /api/v1/auth/logout not found")
	}
}

func TestAuthHandler_Logout_WithBlacklist(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	bl := application.NewTokenBlacklist()

	handler := NewAuthHandlerWithBlacklist(service, bl)

	// Create a valid JWT with future expiry
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"type": "access",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("test-secret"))

	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", "user123") // Simulate auth middleware
		handler.Logout(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !bl.IsRevoked(tokenString) {
		t.Error("Expected token to be revoked after logout")
	}
}

func TestAuthHandler_Logout_NoAuthHeader_WithBlacklist(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	bl := application.NewTokenBlacklist()

	handler := NewAuthHandlerWithBlacklist(service, bl)

	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", "user-1") // Authenticated but no Bearer header
		handler.Logout(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	router.ServeHTTP(w, req)

	// Should succeed even without Bearer header (token just won't be revoked)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthHandler_Logout_InvalidAuthHeader_WithBlacklist(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	bl := application.NewTokenBlacklist()

	handler := NewAuthHandlerWithBlacklist(service, bl)

	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", "user-1")
		handler.Logout(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthHandler_Logout_MalformedJWT_WithBlacklist(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	bl := application.NewTokenBlacklist()

	handler := NewAuthHandlerWithBlacklist(service, bl)

	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", "user-1")
		handler.Logout(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-jwt")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Malformed JWT (not 3 parts) is NOT revoked — only valid JWTs with exp are blacklisted
	if bl.IsRevoked("not-a-valid-jwt") {
		t.Error("Expected malformed token to NOT be added to blacklist")
	}
}

func TestAuthHandler_Logout_ExpiredToken_WithBlacklist(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	bl := application.NewTokenBlacklist()

	handler := NewAuthHandlerWithBlacklist(service, bl)

	// Create a JWT with an expiry in the past (already expired)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"type": "access",
		"exp":  time.Now().Add(-time.Hour).Unix(), // expired 1 hour ago
	})
	tokenString, _ := token.SignedString([]byte("test-secret"))

	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", "user123")
		handler.Logout(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Already-expired token should NOT be added to the blacklist
	if bl.IsRevoked(tokenString) {
		t.Error("Expected already-expired token to NOT be revoked")
	}
}

func TestAuthHandler_Logout_InvalidBase64Payload_WithBlacklist(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	bl := application.NewTokenBlacklist()

	handler := NewAuthHandlerWithBlacklist(service, bl)

	// Construct a token-like string with 3 parts but invalid base64 in the payload
	tokenString := "eyJhbGciOiJIUzI1NiJ9.!!!invalid-base64!!!.signature"

	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", "user123")
		handler.Logout(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Token with invalid base64 payload should NOT be revoked
	if bl.IsRevoked(tokenString) {
		t.Error("Expected token with invalid base64 payload to NOT be revoked")
	}
}

func TestGetValidTokenExpiry_InvalidBase64(t *testing.T) {
	// Test getValidTokenExpiry directly with invalid base64 in payload
	tokenString := "header.!!!not-valid-base64!!!.signature"
	_, valid := getValidTokenExpiry(tokenString)
	if valid {
		t.Error("Expected getValidTokenExpiry to return false for invalid base64 payload")
	}
}

func TestGetValidTokenExpiry_InvalidJSON(t *testing.T) {
	// Test getValidTokenExpiry directly with valid base64 but not JSON
	// "not json" base64-encoded = "bm90IGpzb24"
	tokenString := "header.bm90IGpzb24.signature"
	_, valid := getValidTokenExpiry(tokenString)
	if valid {
		t.Error("Expected getValidTokenExpiry to return false for non-JSON payload")
	}
}

func TestGetValidTokenExpiry_MissingExpClaim(t *testing.T) {
	// Test getValidTokenExpiry directly with valid JSON but no exp
	// {"sub":"user123"} base64url-encoded
	tokenString := "header.eyJzdWIiOiJ1c2VyMTIzIn0.signature"
	_, valid := getValidTokenExpiry(tokenString)
	if valid {
		t.Error("Expected getValidTokenExpiry to return false for missing exp claim")
	}
}

func TestGetValidTokenExpiry_ValidToken(t *testing.T) {
	// Create a real JWT to test the happy path
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("test-secret"))

	expiry, valid := getValidTokenExpiry(tokenString)
	if !valid {
		t.Error("Expected getValidTokenExpiry to return true for valid token")
	}
	if expiry.IsZero() {
		t.Error("Expected non-zero expiry time")
	}
}

func TestGetValidTokenExpiry_TooFewParts(t *testing.T) {
	// Test with only 2 parts instead of 3
	_, valid := getValidTokenExpiry("header.payload")
	if valid {
		t.Error("Expected getValidTokenExpiry to return false for 2-part token")
	}
}

func TestAuthHandler_Logout_TokenWithoutExp_WithBlacklist(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewAuthService(repo, "test-secret")
	bl := application.NewTokenBlacklist()

	handler := NewAuthHandlerWithBlacklist(service, bl)

	// Create JWT without exp claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "user123",
		"type": "access",
	})
	tokenString, _ := token.SignedString([]byte("test-secret"))

	router := gin.New()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", "user123")
		handler.Logout(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Token without exp claim is NOT revoked — getValidTokenExpiry returns false
	if bl.IsRevoked(tokenString) {
		t.Error("Expected token without exp to NOT be revoked")
	}
}
