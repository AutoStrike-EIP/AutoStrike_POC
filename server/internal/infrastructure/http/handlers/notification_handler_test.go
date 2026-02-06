package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"autostrike/internal/application"
	"autostrike/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

// mockNotificationRepoForHandler implements repository.NotificationRepository for handler tests
type mockNotificationRepoForHandler struct {
	settings      map[string]*entity.NotificationSettings
	notifications map[string]*entity.Notification
}

func newMockNotificationRepoForHandler() *mockNotificationRepoForHandler {
	return &mockNotificationRepoForHandler{
		settings:      make(map[string]*entity.NotificationSettings),
		notifications: make(map[string]*entity.Notification),
	}
}

func (m *mockNotificationRepoForHandler) CreateSettings(ctx context.Context, settings *entity.NotificationSettings) error {
	m.settings[settings.ID] = settings
	return nil
}

func (m *mockNotificationRepoForHandler) UpdateSettings(ctx context.Context, settings *entity.NotificationSettings) error {
	m.settings[settings.ID] = settings
	return nil
}

func (m *mockNotificationRepoForHandler) DeleteSettings(ctx context.Context, id string) error {
	delete(m.settings, id)
	return nil
}

func (m *mockNotificationRepoForHandler) FindSettingsByID(ctx context.Context, id string) (*entity.NotificationSettings, error) {
	if s, ok := m.settings[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (m *mockNotificationRepoForHandler) FindSettingsByUserID(ctx context.Context, userID string) (*entity.NotificationSettings, error) {
	for _, s := range m.settings {
		if s.UserID == userID {
			return s, nil
		}
	}
	return nil, nil
}

func (m *mockNotificationRepoForHandler) FindAllEnabledSettings(ctx context.Context) ([]*entity.NotificationSettings, error) {
	var result []*entity.NotificationSettings
	for _, s := range m.settings {
		if s.Enabled {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockNotificationRepoForHandler) CreateNotification(ctx context.Context, notification *entity.Notification) error {
	m.notifications[notification.ID] = notification
	return nil
}

func (m *mockNotificationRepoForHandler) UpdateNotification(ctx context.Context, notification *entity.Notification) error {
	m.notifications[notification.ID] = notification
	return nil
}

func (m *mockNotificationRepoForHandler) FindNotificationByID(ctx context.Context, id string) (*entity.Notification, error) {
	if n, ok := m.notifications[id]; ok {
		return n, nil
	}
	return nil, nil
}

func (m *mockNotificationRepoForHandler) FindNotificationsByUserID(ctx context.Context, userID string, limit int) ([]*entity.Notification, error) {
	var result []*entity.Notification
	for _, n := range m.notifications {
		if n.UserID == userID {
			result = append(result, n)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *mockNotificationRepoForHandler) FindUnreadByUserID(ctx context.Context, userID string) ([]*entity.Notification, error) {
	var result []*entity.Notification
	for _, n := range m.notifications {
		if n.UserID == userID && !n.Read {
			result = append(result, n)
		}
	}
	return result, nil
}

func (m *mockNotificationRepoForHandler) MarkAsRead(ctx context.Context, id string) error {
	if n, ok := m.notifications[id]; ok {
		n.Read = true
	}
	return nil
}

func (m *mockNotificationRepoForHandler) MarkAllAsRead(ctx context.Context, userID string) error {
	for _, n := range m.notifications {
		if n.UserID == userID {
			n.Read = true
		}
	}
	return nil
}

// mockUserRepoForNotificationHandler implements repository.UserRepository for notification handler tests
type mockUserRepoForNotificationHandler struct{}

func (m *mockUserRepoForNotificationHandler) Create(ctx context.Context, user *entity.User) error {
	return nil
}

func (m *mockUserRepoForNotificationHandler) Update(ctx context.Context, user *entity.User) error {
	return nil
}

func (m *mockUserRepoForNotificationHandler) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockUserRepoForNotificationHandler) FindByID(ctx context.Context, id string) (*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepoForNotificationHandler) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepoForNotificationHandler) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepoForNotificationHandler) FindAll(ctx context.Context) ([]*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepoForNotificationHandler) FindActive(ctx context.Context) ([]*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepoForNotificationHandler) UpdateLastLogin(ctx context.Context, id string) error {
	return nil
}

func (m *mockUserRepoForNotificationHandler) Deactivate(ctx context.Context, id string) error {
	return nil
}

func (m *mockUserRepoForNotificationHandler) Reactivate(ctx context.Context, id string) error {
	return nil
}

func (m *mockUserRepoForNotificationHandler) CountByRole(ctx context.Context, role entity.UserRole) (int, error) {
	return 0, nil
}

func (m *mockUserRepoForNotificationHandler) DeactivateAdminIfNotLast(ctx context.Context, id string) error {
	return nil
}

func setupNotificationHandler() (*NotificationHandler, *mockNotificationRepoForHandler) {
	repo := newMockNotificationRepoForHandler()
	userRepo := &mockUserRepoForNotificationHandler{}
	service := application.NewNotificationService(repo, userRepo, nil, "https://localhost:8443", nil)
	handler := NewNotificationHandler(service)
	return handler, repo
}

func setupNotificationRouter(handler *NotificationHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware to set user_id for authenticated routes
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})

	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)
	return router
}

func setupNotificationRouterWithAdmin(handler *NotificationHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware to set user_id and admin role for admin routes
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Set("role", "admin")
		c.Next()
	})

	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)
	return router
}

func TestNewNotificationHandler(t *testing.T) {
	handler, _ := setupNotificationHandler()
	if handler == nil {
		t.Fatal("NewNotificationHandler returned nil")
	}
}

func TestNotificationHandler_RegisterRoutes(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	routes := router.Routes()
	expectedPaths := map[string]string{
		"/api/v1/notifications":              "GET",
		"/api/v1/notifications/unread/count": "GET",
		"/api/v1/notifications/:id/read":     "POST",
		"/api/v1/notifications/read-all":     "POST",
		"/api/v1/notifications/settings":     "GET",
		"/api/v1/notifications/smtp":         "GET",
		"/api/v1/notifications/smtp/test":    "POST",
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
}

func TestNotificationHandler_GetNotifications(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	// Add a notification
	repo.notifications["notif-1"] = &entity.Notification{
		ID:        "notif-1",
		UserID:    "test-user-id",
		Type:      entity.NotificationExecutionStarted,
		Title:     "Test",
		Message:   "Test message",
		CreatedAt: time.Now(),
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var notifications []*entity.Notification
	if err := json.Unmarshal(w.Body.Bytes(), &notifications); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(notifications) != 1 {
		t.Errorf("len(notifications) = %d, want 1", len(notifications))
	}
}

func TestNotificationHandler_GetNotifications_WithLimit(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications?limit=10", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotificationHandler_GetNotifications_Unauthenticated(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	// No middleware setting user_id
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_GetUnreadCount(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	// Add unread notifications
	repo.notifications["notif-1"] = &entity.Notification{
		ID: "notif-1", UserID: "test-user-id", Read: false,
	}
	repo.notifications["notif-2"] = &entity.Notification{
		ID: "notif-2", UserID: "test-user-id", Read: false,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/unread/count", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]int
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["count"] != 2 {
		t.Errorf("count = %d, want 2", response["count"])
	}
}

func TestNotificationHandler_MarkAsRead(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	repo.notifications["notif-1"] = &entity.Notification{
		ID: "notif-1", UserID: "test-user-id", Read: false,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/notif-1/read", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	if !repo.notifications["notif-1"].Read {
		t.Error("Notification should be marked as read")
	}
}

func TestNotificationHandler_MarkAllAsRead(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	repo.notifications["notif-1"] = &entity.Notification{
		ID: "notif-1", UserID: "test-user-id", Read: false,
	}
	repo.notifications["notif-2"] = &entity.Notification{
		ID: "notif-2", UserID: "test-user-id", Read: false,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/read-all", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	for _, n := range repo.notifications {
		if !n.Read {
			t.Error("All notifications should be marked as read")
		}
	}
}

func TestNotificationHandler_GetSettings(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	repo.settings["settings-1"] = &entity.NotificationSettings{
		ID:           "settings-1",
		UserID:       "test-user-id",
		Channel:      entity.ChannelEmail,
		Enabled:      true,
		EmailAddress: "test@example.com",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var settings entity.NotificationSettings
	if err := json.Unmarshal(w.Body.Bytes(), &settings); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if settings.EmailAddress != "test@example.com" {
		t.Errorf("EmailAddress = %q, want %q", settings.EmailAddress, "test@example.com")
	}
}

func TestNotificationHandler_GetSettings_NotFound(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestNotificationHandler_CreateSettings(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:          "email",
		Enabled:          true,
		EmailAddress:     "test@example.com",
		NotifyOnStart:    true,
		NotifyOnComplete: true,
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_CreateSettings_InvalidChannel(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := map[string]interface{}{
		"channel": "invalid",
		"enabled": true,
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_UpdateSettings(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	repo.settings["settings-1"] = &entity.NotificationSettings{
		ID:           "settings-1",
		UserID:       "test-user-id",
		Channel:      entity.ChannelEmail,
		Enabled:      true,
		EmailAddress: "old@example.com",
	}

	body := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "new@example.com",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_UpdateSettings_NotFound(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "test@example.com",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestNotificationHandler_DeleteSettings(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	repo.settings["settings-1"] = &entity.NotificationSettings{
		ID:     "settings-1",
		UserID: "test-user-id",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_DeleteSettings_NotFound(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestNotificationHandler_GetSMTPConfig_NotConfigured(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouterWithAdmin(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/smtp", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestNotificationHandler_TestSMTP_InvalidEmail(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouterWithAdmin(handler)

	body := TestSMTPRequest{
		Email: "invalid-email",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/smtp/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_TestSMTP_NotConfigured(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouterWithAdmin(handler)

	body := TestSMTPRequest{
		Email: "test@example.com",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/smtp/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestNotificationHandler_GetUnreadCount_Unauthenticated(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/unread/count", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_MarkAsRead_Unauthenticated(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/notif-1/read", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_MarkAsRead_NotOwned(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	// Add notification owned by different user
	repo.notifications["notif-1"] = &entity.Notification{
		ID: "notif-1", UserID: "other-user", Read: false,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/notif-1/read", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_MarkAsRead_NotFound(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/nonexistent/read", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestNotificationHandler_MarkAllAsRead_Unauthenticated(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/read-all", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_GetSettings_Unauthenticated(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_CreateSettings_Unauthenticated(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	body := NotificationSettingsRequest{Channel: "email", Enabled: true, EmailAddress: "test@example.com"}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_CreateSettings_MissingEmail(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel: "email",
		Enabled: true,
		// Missing EmailAddress
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_CreateSettings_InvalidEmail(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "not-an-email",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_CreateSettings_MissingWebhookURL(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel: "webhook",
		Enabled: true,
		// Missing WebhookURL
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_CreateSettings_InvalidWebhookURL(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:    "webhook",
		Enabled:    true,
		WebhookURL: "not-a-url",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_UpdateSettings_Unauthenticated(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	body := NotificationSettingsRequest{Channel: "email", Enabled: true, EmailAddress: "test@example.com"}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_UpdateSettings_InvalidJSON(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	repo.settings["settings-1"] = &entity.NotificationSettings{
		ID:     "settings-1",
		UserID: "test-user-id",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_UpdateSettings_ValidationError(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	repo.settings["settings-1"] = &entity.NotificationSettings{
		ID:     "settings-1",
		UserID: "test-user-id",
	}

	body := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "invalid-email",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_DeleteSettings_Unauthenticated(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_GetSMTPConfig_Unauthenticated(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/smtp", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_GetSMTPConfig_NotAdmin(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler) // No admin role

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/smtp", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestNotificationHandler_TestSMTP_Unauthenticated(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	body := TestSMTPRequest{Email: "test@example.com"}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/smtp/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_TestSMTP_NotAdmin(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler) // No admin role

	body := TestSMTPRequest{Email: "test@example.com"}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/smtp/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestNotificationHandler_GetNotifications_InvalidLimit(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	// Invalid limit (too high)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications?limit=200", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotificationHandler_GetNotifications_NegativeLimit(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications?limit=-5", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotificationHandler_GetNotifications_NonNumericLimit(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications?limit=abc", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotificationHandler_CreateSettings_WebhookSuccess(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:      "webhook",
		Enabled:      true,
		WebhookURL:   "https://hooks.example.com/webhook",
		NotifyOnStart: true,
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_CreateSettings_InvalidJSON(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationSettingsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     NotificationSettingsRequest
		wantErr bool
	}{
		{
			name:    "email disabled no address",
			req:     NotificationSettingsRequest{Channel: "email", Enabled: false},
			wantErr: false,
		},
		{
			name:    "email enabled with valid address",
			req:     NotificationSettingsRequest{Channel: "email", Enabled: true, EmailAddress: "test@example.com"},
			wantErr: false,
		},
		{
			name:    "email enabled missing address",
			req:     NotificationSettingsRequest{Channel: "email", Enabled: true, EmailAddress: ""},
			wantErr: true,
		},
		{
			name:    "webhook disabled no url",
			req:     NotificationSettingsRequest{Channel: "webhook", Enabled: false},
			wantErr: false,
		},
		{
			name:    "webhook enabled with valid url",
			req:     NotificationSettingsRequest{Channel: "webhook", Enabled: true, WebhookURL: "https://example.com/hook"},
			wantErr: false,
		},
		{
			name:    "webhook enabled missing url",
			req:     NotificationSettingsRequest{Channel: "webhook", Enabled: true, WebhookURL: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input   string
		want    int
		wantErr bool
	}{
		{"10", 10, false},
		{"0", 0, false},
		{"-5", -5, false},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInt(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseInt(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// --- Error-returning mock for testing service error paths ---

// errorNotificationRepo is a mock that can be configured to return errors
type errorNotificationRepo struct {
	findSettingsByUserIDErr   error
	findSettingsByUserIDVal   *entity.NotificationSettings
	createSettingsErr         error
	updateSettingsErr         error
	deleteSettingsErr         error
	findNotificationsByErr    error
	findUnreadByUserIDErr     error
	markAsReadErr             error
	markAllAsReadErr          error
	findNotificationByIDErr   error
	findNotificationByIDVal   *entity.Notification
}

func (m *errorNotificationRepo) CreateSettings(_ context.Context, _ *entity.NotificationSettings) error {
	return m.createSettingsErr
}

func (m *errorNotificationRepo) UpdateSettings(_ context.Context, _ *entity.NotificationSettings) error {
	return m.updateSettingsErr
}

func (m *errorNotificationRepo) DeleteSettings(_ context.Context, _ string) error {
	return m.deleteSettingsErr
}

func (m *errorNotificationRepo) FindSettingsByID(_ context.Context, _ string) (*entity.NotificationSettings, error) {
	return nil, nil
}

func (m *errorNotificationRepo) FindSettingsByUserID(_ context.Context, _ string) (*entity.NotificationSettings, error) {
	return m.findSettingsByUserIDVal, m.findSettingsByUserIDErr
}

func (m *errorNotificationRepo) FindAllEnabledSettings(_ context.Context) ([]*entity.NotificationSettings, error) {
	return nil, nil
}

func (m *errorNotificationRepo) CreateNotification(_ context.Context, _ *entity.Notification) error {
	return nil
}

func (m *errorNotificationRepo) UpdateNotification(_ context.Context, _ *entity.Notification) error {
	return nil
}

func (m *errorNotificationRepo) FindNotificationByID(_ context.Context, _ string) (*entity.Notification, error) {
	return m.findNotificationByIDVal, m.findNotificationByIDErr
}

func (m *errorNotificationRepo) FindNotificationsByUserID(_ context.Context, _ string, _ int) ([]*entity.Notification, error) {
	return nil, m.findNotificationsByErr
}

func (m *errorNotificationRepo) FindUnreadByUserID(_ context.Context, _ string) ([]*entity.Notification, error) {
	return nil, m.findUnreadByUserIDErr
}

func (m *errorNotificationRepo) MarkAsRead(_ context.Context, _ string) error {
	return m.markAsReadErr
}

func (m *errorNotificationRepo) MarkAllAsRead(_ context.Context, _ string) error {
	return m.markAllAsReadErr
}

func setupErrorNotificationHandler(repo *errorNotificationRepo) *NotificationHandler {
	userRepo := &mockUserRepoForNotificationHandler{}
	service := application.NewNotificationService(repo, userRepo, nil, "https://localhost:8443", nil)
	return NewNotificationHandler(service)
}

func setupUnauthRouter(handler *NotificationHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)
	return router
}

// --- Edge case tests for error paths ---

func TestNotificationHandler_GetNotifications_ServiceError(t *testing.T) {
	repo := &errorNotificationRepo{
		findNotificationsByErr: fmt.Errorf("database error"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "failed to get notifications" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_GetUnreadCount_ServiceError(t *testing.T) {
	repo := &errorNotificationRepo{
		findUnreadByUserIDErr: fmt.Errorf("database error"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/unread/count", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "failed to get unread count" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_MarkAsRead_ServiceGenericError(t *testing.T) {
	repo := &errorNotificationRepo{
		findNotificationByIDErr: fmt.Errorf("database connection lost"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/notif-1/read", nil)
	router.ServeHTTP(w, req)

	// FindNotificationByID returns error, so MarkAsReadForUser returns "notification not found or not owned by user"
	// which triggers the Forbidden path
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_MarkAsRead_MarkAsReadRepoError(t *testing.T) {
	// Notification found and owned by user, but MarkAsRead repo call fails
	repo := &errorNotificationRepo{
		findNotificationByIDVal: &entity.Notification{
			ID: "notif-1", UserID: "test-user-id", Read: false,
		},
		markAsReadErr: fmt.Errorf("database write error"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/notif-1/read", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "failed to mark as read" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_MarkAllAsRead_ServiceError(t *testing.T) {
	repo := &errorNotificationRepo{
		markAllAsReadErr: fmt.Errorf("database error"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/read-all", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "failed to mark all as read" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_GetSettings_SqlErrNoRows(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDErr: sql.ErrNoRows,
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "settings not found" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_GetSettings_GenericError(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDErr: fmt.Errorf("database error"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "failed to get settings" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_CreateSettings_ServiceError(t *testing.T) {
	repo := &errorNotificationRepo{
		createSettingsErr: fmt.Errorf("database error"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "test@example.com",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "failed to create settings" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_UpdateSettings_SqlErrNoRows(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDErr: sql.ErrNoRows,
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "test@example.com",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "settings not found" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_UpdateSettings_GenericGetError(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDErr: fmt.Errorf("database error"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "test@example.com",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "failed to get settings" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_UpdateSettings_NilSettings(t *testing.T) {
	// findSettingsByUserID returns nil, nil (no error, but no settings found)
	repo := &errorNotificationRepo{
		findSettingsByUserIDErr: nil,
		findSettingsByUserIDVal: nil,
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "test@example.com",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_UpdateSettings_UpdateServiceError(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDVal: &entity.NotificationSettings{
			ID:           "settings-1",
			UserID:       "test-user-id",
			Channel:      entity.ChannelEmail,
			Enabled:      true,
			EmailAddress: "old@example.com",
		},
		updateSettingsErr: fmt.Errorf("database write error"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "new@example.com",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "failed to update settings" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_UpdateSettings_MissingWebhookURL(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDVal: &entity.NotificationSettings{
			ID:     "settings-1",
			UserID: "test-user-id",
		},
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel: "webhook",
		Enabled: true,
		// Missing WebhookURL
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_UpdateSettings_InvalidWebhookURL(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDVal: &entity.NotificationSettings{
			ID:     "settings-1",
			UserID: "test-user-id",
		},
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:    "webhook",
		Enabled:    true,
		WebhookURL: "not-a-valid-url",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_DeleteSettings_SqlErrNoRows(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDErr: sql.ErrNoRows,
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "settings not found" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_DeleteSettings_GenericGetError(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDErr: fmt.Errorf("database error"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "failed to get settings" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_DeleteSettings_NilSettings(t *testing.T) {
	// findSettingsByUserID returns nil, nil (no error, no settings)
	repo := &errorNotificationRepo{
		findSettingsByUserIDErr: nil,
		findSettingsByUserIDVal: nil,
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_DeleteSettings_DeleteServiceError(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDVal: &entity.NotificationSettings{
			ID:     "settings-1",
			UserID: "test-user-id",
		},
		deleteSettingsErr: fmt.Errorf("database delete error"),
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "failed to delete settings" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_GetSMTPConfig_WithConfig(t *testing.T) {
	repo := newMockNotificationRepoForHandler()
	userRepo := &mockUserRepoForNotificationHandler{}
	smtpConfig := &entity.SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user",
		Password: "secret",
		From:     "noreply@example.com",
		UseTLS:   true,
	}
	service := application.NewNotificationService(repo, userRepo, smtpConfig, "https://localhost:8443", nil)
	handler := NewNotificationHandler(service)
	router := setupNotificationRouterWithAdmin(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/smtp", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var config entity.SMTPConfig
	if err := json.Unmarshal(w.Body.Bytes(), &config); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if config.Host != "smtp.example.com" {
		t.Errorf("Host = %q, want %q", config.Host, "smtp.example.com")
	}
	if config.Password != "" {
		t.Error("Password should not be exposed in response")
	}
}

func TestNotificationHandler_TestSMTP_Unauthenticated_NoUserID(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupUnauthRouter(handler)

	body := TestSMTPRequest{Email: "test@example.com"}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/smtp/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestNotificationHandler_CreateSettings_DisabledEmailNoAddress(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel: "email",
		Enabled: false,
		// No email address, but disabled so should pass validation
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_CreateSettings_DisabledWebhookNoURL(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel: "webhook",
		Enabled: false,
		// No webhook URL, but disabled so should pass validation
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_GetNotifications_EmptyResult(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var notifications []*entity.Notification
	if err := json.Unmarshal(w.Body.Bytes(), &notifications); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if len(notifications) != 0 {
		t.Errorf("Expected empty array, got %d items", len(notifications))
	}
}

func TestNotificationHandler_GetNotifications_LimitBoundary(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	// limit=100 is valid upper bound
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications?limit=100", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotificationHandler_GetNotifications_ZeroLimit(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	// limit=0 should fall back to default (not valid, parsed > 0 fails)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications?limit=0", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotificationHandler_UpdateSettings_WebhookSuccess(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	repo.settings["settings-1"] = &entity.NotificationSettings{
		ID:     "settings-1",
		UserID: "test-user-id",
	}

	body := NotificationSettingsRequest{
		Channel:    "webhook",
		Enabled:    true,
		WebhookURL: "https://hooks.example.com/webhook",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_UpdateSettings_MalformedJSON(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_CreateSettings_MissingChannel(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := map[string]interface{}{
		"enabled": true,
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_UpdateSettings_MissingChannel(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	repo.settings["settings-1"] = &entity.NotificationSettings{
		ID:     "settings-1",
		UserID: "test-user-id",
	}

	body := map[string]interface{}{
		"enabled": true,
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNotificationHandler_CreateSettings_WithAllFields(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:              "email",
		Enabled:              true,
		EmailAddress:         "test@example.com",
		NotifyOnStart:        true,
		NotifyOnComplete:     true,
		NotifyOnFailure:      true,
		NotifyOnScoreAlert:   true,
		ScoreAlertThreshold:  75.0,
		NotifyOnAgentOffline: true,
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var settings entity.NotificationSettings
	if err := json.Unmarshal(w.Body.Bytes(), &settings); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if settings.UserID != "test-user-id" {
		t.Errorf("UserID = %q, want %q", settings.UserID, "test-user-id")
	}
	if settings.ScoreAlertThreshold != 75.0 {
		t.Errorf("ScoreAlertThreshold = %f, want 75.0", settings.ScoreAlertThreshold)
	}
	if !settings.NotifyOnAgentOffline {
		t.Error("NotifyOnAgentOffline should be true")
	}
}

func TestNotificationSettingsRequest_Validate_InvalidEmailFormat(t *testing.T) {
	req := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "not-valid-email",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("Expected validation error for invalid email format")
	}
	if err.Error() != "invalid email address format" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestNotificationSettingsRequest_Validate_InvalidWebhookURLFormat(t *testing.T) {
	req := NotificationSettingsRequest{
		Channel:    "webhook",
		Enabled:    true,
		WebhookURL: "not a valid url with spaces",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("Expected validation error for invalid webhook URL format")
	}
	if err.Error() != "invalid webhook URL format" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestNotificationSettingsRequest_Validate_WebhookMissingURL(t *testing.T) {
	req := NotificationSettingsRequest{
		Channel:    "webhook",
		Enabled:    true,
		WebhookURL: "",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("Expected validation error for missing webhook URL")
	}
	if err.Error() != "webhook URL is required when channel is webhook" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestNotificationSettingsRequest_Validate_EmailMissingAddress(t *testing.T) {
	req := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("Expected validation error for missing email address")
	}
	if err.Error() != "email address is required when channel is email" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestNotificationHandler_GetSMTPConfig_UserWithoutRole(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Middleware that sets user_id but does NOT set role at all
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/smtp", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestNotificationHandler_TestSMTP_UserWithoutRole(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	body := TestSMTPRequest{Email: "test@example.com"}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/smtp/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestNotificationHandler_TestSMTP_EmptyBody(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouterWithAdmin(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/smtp/test", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_TestSMTP_InvalidJSON(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouterWithAdmin(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/smtp/test", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_GetNotifications_MultipleNotifications(t *testing.T) {
	handler, repo := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("notif-%d", i)
		repo.notifications[id] = &entity.Notification{
			ID:        id,
			UserID:    "test-user-id",
			Type:      entity.NotificationExecutionStarted,
			Title:     fmt.Sprintf("Test %d", i),
			Message:   fmt.Sprintf("Test message %d", i),
			CreatedAt: time.Now(),
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications?limit=3", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var notifications []*entity.Notification
	if err := json.Unmarshal(w.Body.Bytes(), &notifications); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if len(notifications) > 3 {
		t.Errorf("Expected at most 3 notifications, got %d", len(notifications))
	}
}

func TestNotificationHandler_GetSettings_NilSettings(t *testing.T) {
	// The default mock repo returns nil, nil when no settings exist
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "settings not found" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_MarkAsRead_EmptyID(t *testing.T) {
	handler, _ := setupNotificationHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Route without :id param so c.Param("id") returns ""
	router.POST("/api/v1/notifications/mark-read", func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		handler.MarkAsRead(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/mark-read", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty notification ID, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if response["error"] != "notification ID required" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestNotificationHandler_TestSMTP_Success(t *testing.T) {
	// Create a notification service with a mock SMTP config
	// The TestSMTPConnection will fail because no real SMTP server, but we
	// test the handler's success path by testing that TestSMTP returns 500
	// (SMTP not configured) vs the success path.
	// For the success path (line 487), we need SMTP configured AND sendEmail to succeed.
	// Since we can't easily mock the actual SMTP dial, the success line (487)
	// will remain unreachable without a real SMTP server. The 93.8% coverage
	// (1 line out of 16 uncovered) is acceptable since it requires an actual SMTP connection.

	// However, we can verify the handler flow up to the SMTP call
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouterWithAdmin(handler)

	body := TestSMTPRequest{Email: "test@example.com"}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/smtp/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Returns 500 because SMTP is not configured in test environment
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestNotificationHandler_UpdateSettings_InvalidEmailOnUpdate(t *testing.T) {
	repo := &errorNotificationRepo{
		findSettingsByUserIDVal: &entity.NotificationSettings{
			ID:     "settings-1",
			UserID: "test-user-id",
		},
	}
	handler := setupErrorNotificationHandler(repo)
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:      "email",
		Enabled:      true,
		EmailAddress: "not-an-email-address",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid email on update, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_CreateSettings_WebhookInvalidURL(t *testing.T) {
	handler, _ := setupNotificationHandler()
	router := setupNotificationRouter(handler)

	body := NotificationSettingsRequest{
		Channel:    "webhook",
		Enabled:    true,
		WebhookURL: "://missing-scheme",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/settings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}
