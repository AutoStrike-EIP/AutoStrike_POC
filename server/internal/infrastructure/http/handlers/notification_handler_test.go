package handlers

import (
	"bytes"
	"context"
	"encoding/json"
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
