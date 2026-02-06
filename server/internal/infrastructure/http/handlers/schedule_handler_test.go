package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"autostrike/internal/application"
	"autostrike/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// mockScheduleRepo implements repository.ScheduleRepository for testing
type mockScheduleRepo struct {
	schedules    map[string]*entity.Schedule
	runs         map[string][]*entity.ScheduleRun
	createErr    error
	updateErr    error
	deleteErr    error
	findErr      error
	createRunErr error
	findRunsErr  error
}

func newMockScheduleRepo() *mockScheduleRepo {
	return &mockScheduleRepo{
		schedules: make(map[string]*entity.Schedule),
		runs:      make(map[string][]*entity.ScheduleRun),
	}
}

func (m *mockScheduleRepo) Create(ctx context.Context, schedule *entity.Schedule) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.schedules[schedule.ID] = schedule
	return nil
}

func (m *mockScheduleRepo) Update(ctx context.Context, schedule *entity.Schedule) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.schedules[schedule.ID] = schedule
	return nil
}

func (m *mockScheduleRepo) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.schedules, id)
	return nil
}

func (m *mockScheduleRepo) FindByID(ctx context.Context, id string) (*entity.Schedule, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	if s, ok := m.schedules[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (m *mockScheduleRepo) FindAll(ctx context.Context) ([]*entity.Schedule, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	var result []*entity.Schedule
	for _, s := range m.schedules {
		result = append(result, s)
	}
	return result, nil
}

func (m *mockScheduleRepo) FindByStatus(ctx context.Context, status entity.ScheduleStatus) ([]*entity.Schedule, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	var result []*entity.Schedule
	for _, s := range m.schedules {
		if s.Status == status {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockScheduleRepo) FindActiveSchedulesDue(ctx context.Context, now time.Time) ([]*entity.Schedule, error) {
	return nil, nil
}

func (m *mockScheduleRepo) FindByScenarioID(ctx context.Context, scenarioID string) ([]*entity.Schedule, error) {
	var result []*entity.Schedule
	for _, s := range m.schedules {
		if s.ScenarioID == scenarioID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockScheduleRepo) CreateRun(ctx context.Context, run *entity.ScheduleRun) error {
	if m.createRunErr != nil {
		return m.createRunErr
	}
	m.runs[run.ScheduleID] = append(m.runs[run.ScheduleID], run)
	return nil
}

func (m *mockScheduleRepo) UpdateRun(ctx context.Context, run *entity.ScheduleRun) error {
	return nil
}

func (m *mockScheduleRepo) FindRunsByScheduleID(ctx context.Context, scheduleID string, limit int) ([]*entity.ScheduleRun, error) {
	if m.findRunsErr != nil {
		return nil, m.findRunsErr
	}
	runs := m.runs[scheduleID]
	if len(runs) > limit {
		runs = runs[:limit]
	}
	return runs, nil
}

// setupRealScheduleHandler creates a handler with real service using mock repo
func setupRealScheduleHandler(repo *mockScheduleRepo) (*ScheduleHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	// Create real service with mock repo (nil execution service for most tests)
	service := application.NewScheduleService(repo, nil, logger)
	handler := NewScheduleHandler(service)

	router := gin.New()
	api := router.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	})
	handler.RegisterRoutes(api)

	return handler, router
}

// setupRealScheduleHandlerNoAuth creates a handler without auth middleware
func setupRealScheduleHandlerNoAuth(repo *mockScheduleRepo) (*ScheduleHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	service := application.NewScheduleService(repo, nil, logger)
	handler := NewScheduleHandler(service)

	router := gin.New()
	api := router.Group("/api/v1")
	handler.RegisterRoutes(api)

	return handler, router
}

func TestScheduleHandler_GetAll_Success(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:   "sched-1",
		Name: "Test Schedule",
	}
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var schedules []*entity.Schedule
	if err := json.Unmarshal(w.Body.Bytes(), &schedules); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(schedules) != 1 {
		t.Errorf("len(schedules) = %d, want 1", len(schedules))
	}
}

func TestScheduleHandler_GetAll_Empty(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var schedules []*entity.Schedule
	if err := json.Unmarshal(w.Body.Bytes(), &schedules); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(schedules) != 0 {
		t.Errorf("len(schedules) = %d, want 0", len(schedules))
	}
}

func TestScheduleHandler_GetAll_Unauthorized(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandlerNoAuth(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestScheduleHandler_GetAll_Error(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.findErr = errors.New("database error")
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestScheduleHandler_GetByID_Success(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:   "sched-1",
		Name: "Test Schedule",
	}
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var schedule entity.Schedule
	if err := json.Unmarshal(w.Body.Bytes(), &schedule); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if schedule.Name != "Test Schedule" {
		t.Errorf("Name = %q, want %q", schedule.Name, "Test Schedule")
	}
}

func TestScheduleHandler_GetByID_NotFound(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestScheduleHandler_GetByID_Unauthorized(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandlerNoAuth(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestScheduleHandler_GetByID_Error(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.findErr = errors.New("database error")
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestScheduleHandler_Create_Success(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "New Schedule",
		ScenarioID: "scenario-1",
		Frequency:  "daily",
		SafeMode:   true,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status = %d, want %d. Body: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var schedule entity.Schedule
	if err := json.Unmarshal(w.Body.Bytes(), &schedule); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if schedule.Name != "New Schedule" {
		t.Errorf("Name = %q, want %q", schedule.Name, "New Schedule")
	}
}

func TestScheduleHandler_Create_WithStartAt(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	startAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	body := CreateScheduleRequest{
		Name:       "Scheduled Task",
		ScenarioID: "scenario-1",
		Frequency:  "once",
		StartAt:    startAt,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status = %d, want %d. Body: %s", w.Code, http.StatusCreated, w.Body.String())
	}
}

func TestScheduleHandler_Create_InvalidStartAt(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "Scheduled Task",
		ScenarioID: "scenario-1",
		Frequency:  "once",
		StartAt:    "invalid-date",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestScheduleHandler_Create_InvalidJSON(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestScheduleHandler_Create_MissingFields(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	body := map[string]string{
		"name": "Test",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestScheduleHandler_Create_Unauthorized(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandlerNoAuth(repo)

	body := CreateScheduleRequest{
		Name:       "New Schedule",
		ScenarioID: "scenario-1",
		Frequency:  "daily",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestScheduleHandler_Create_Error(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.createErr = errors.New("database error")
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "New Schedule",
		ScenarioID: "scenario-1",
		Frequency:  "daily",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestScheduleHandler_Update_Success(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:     "sched-1",
		Name:   "Old Name",
		Status: entity.ScheduleStatusActive,
	}
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "Updated Name",
		ScenarioID: "scenario-1",
		Frequency:  "hourly",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/schedules/sched-1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var schedule entity.Schedule
	if err := json.Unmarshal(w.Body.Bytes(), &schedule); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if schedule.Name != "Updated Name" {
		t.Errorf("Name = %q, want %q", schedule.Name, "Updated Name")
	}
}

func TestScheduleHandler_Update_WithStartAt(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:     "sched-1",
		Name:   "Old Name",
		Status: entity.ScheduleStatusActive,
	}
	_, router := setupRealScheduleHandler(repo)

	startAt := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
	body := CreateScheduleRequest{
		Name:       "Updated Name",
		ScenarioID: "scenario-1",
		Frequency:  "once",
		StartAt:    startAt,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/schedules/sched-1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestScheduleHandler_Update_InvalidStartAt(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:     "sched-1",
		Name:   "Old Name",
		Status: entity.ScheduleStatusActive,
	}
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "Updated Name",
		ScenarioID: "scenario-1",
		Frequency:  "once",
		StartAt:    "bad-date-format",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/schedules/sched-1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestScheduleHandler_Update_NotFound(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "Test",
		ScenarioID: "scenario-1",
		Frequency:  "daily",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/schedules/nonexistent", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestScheduleHandler_Update_Unauthorized(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandlerNoAuth(repo)

	body := CreateScheduleRequest{
		Name:       "Test",
		ScenarioID: "scenario-1",
		Frequency:  "daily",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/schedules/sched-1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestScheduleHandler_Update_InvalidJSON(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/schedules/sched-1", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestScheduleHandler_Delete_Success(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{ID: "sched-1"}
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/schedules/sched-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestScheduleHandler_Delete_Unauthorized(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandlerNoAuth(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/schedules/sched-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestScheduleHandler_Delete_Error(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.deleteErr = errors.New("database error")
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/schedules/sched-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestScheduleHandler_Pause_Success(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:     "sched-1",
		Status: entity.ScheduleStatusActive,
	}
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/pause", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var schedule entity.Schedule
	if err := json.Unmarshal(w.Body.Bytes(), &schedule); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if schedule.Status != entity.ScheduleStatusPaused {
		t.Errorf("Status = %q, want %q", schedule.Status, entity.ScheduleStatusPaused)
	}
}

func TestScheduleHandler_Pause_NotFound(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/nonexistent/pause", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestScheduleHandler_Pause_Unauthorized(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandlerNoAuth(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/pause", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestScheduleHandler_Resume_Success(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:     "sched-1",
		Status: entity.ScheduleStatusPaused,
	}
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/resume", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var schedule entity.Schedule
	if err := json.Unmarshal(w.Body.Bytes(), &schedule); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if schedule.Status != entity.ScheduleStatusActive {
		t.Errorf("Status = %q, want %q", schedule.Status, entity.ScheduleStatusActive)
	}
}

func TestScheduleHandler_Resume_NotFound(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/nonexistent/resume", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestScheduleHandler_Resume_Unauthorized(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandlerNoAuth(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/resume", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestScheduleHandler_RunNow_NotFound(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/nonexistent/run", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestScheduleHandler_RunNow_Unauthorized(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandlerNoAuth(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/run", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestScheduleHandler_GetRuns_Success(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.runs["sched-1"] = []*entity.ScheduleRun{
		{ID: "run-1", ScheduleID: "sched-1", Status: "completed"},
		{ID: "run-2", ScheduleID: "sched-1", Status: "failed"},
	}
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1/runs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var runs []*entity.ScheduleRun
	if err := json.Unmarshal(w.Body.Bytes(), &runs); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(runs) != 2 {
		t.Errorf("len(runs) = %d, want 2", len(runs))
	}
}

func TestScheduleHandler_GetRuns_WithLimit(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.runs["sched-1"] = []*entity.ScheduleRun{
		{ID: "run-1", ScheduleID: "sched-1"},
		{ID: "run-2", ScheduleID: "sched-1"},
		{ID: "run-3", ScheduleID: "sched-1"},
	}
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1/runs?limit=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var runs []*entity.ScheduleRun
	if err := json.Unmarshal(w.Body.Bytes(), &runs); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(runs) != 2 {
		t.Errorf("len(runs) = %d, want 2", len(runs))
	}
}

func TestScheduleHandler_GetRuns_InvalidLimit(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	// Invalid limit should use default (20)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1/runs?limit=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestScheduleHandler_GetRuns_Empty(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1/runs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var runs []*entity.ScheduleRun
	if err := json.Unmarshal(w.Body.Bytes(), &runs); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(runs) != 0 {
		t.Errorf("len(runs) = %d, want 0", len(runs))
	}
}

func TestScheduleHandler_GetRuns_Unauthorized(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandlerNoAuth(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1/runs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestNewScheduleHandler(t *testing.T) {
	handler := NewScheduleHandler(nil)
	if handler == nil {
		t.Error("NewScheduleHandler returned nil")
	}
}

func TestScheduleHandler_RegisterRoutes(t *testing.T) {
	handler := NewScheduleHandler(nil)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")

	// Should not panic
	handler.RegisterRoutes(api)
}

func TestScheduleHandler_Create_CronFrequency(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "Cron Schedule",
		ScenarioID: "scenario-1",
		Frequency:  "cron",
		CronExpr:   "0 * * * *",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status = %d, want %d. Body: %s", w.Code, http.StatusCreated, w.Body.String())
	}
}

func TestScheduleHandler_Create_CronFrequency_MissingExpr(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "Cron Schedule",
		ScenarioID: "scenario-1",
		Frequency:  "cron",
		// Missing CronExpr
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestScheduleHandler_Create_CronFrequency_InvalidExpr(t *testing.T) {
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "Cron Schedule",
		ScenarioID: "scenario-1",
		Frequency:  "cron",
		CronExpr:   "invalid cron",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// --- mockFailingScenarioRepo always returns an error from FindByID ---
type mockFailingScenarioRepo struct{}

func (m *mockFailingScenarioRepo) Create(_ context.Context, _ *entity.Scenario) error   { return nil }
func (m *mockFailingScenarioRepo) Update(_ context.Context, _ *entity.Scenario) error   { return nil }
func (m *mockFailingScenarioRepo) Delete(_ context.Context, _ string) error              { return nil }
func (m *mockFailingScenarioRepo) FindByID(_ context.Context, _ string) (*entity.Scenario, error) {
	return nil, errors.New("scenario not found")
}
func (m *mockFailingScenarioRepo) FindAll(_ context.Context) ([]*entity.Scenario, error) {
	return nil, nil
}
func (m *mockFailingScenarioRepo) FindByTag(_ context.Context, _ string) ([]*entity.Scenario, error) {
	return nil, nil
}
func (m *mockFailingScenarioRepo) ImportFromYAML(_ context.Context, _ string) error { return nil }

// setupScheduleHandlerWithExecService creates a handler with a real ScheduleService
// that has an ExecutionService backed by a failing scenario repo. This allows
// RunNow to be tested: the execution will fail (scenario not found) but the
// run record is still saved, so the handler returns 200 OK.
func setupScheduleHandlerWithExecService(repo *mockScheduleRepo) (*ScheduleHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	execService := application.NewExecutionService(nil, &mockFailingScenarioRepo{}, nil, nil, nil, nil)
	service := application.NewScheduleService(repo, execService, logger)
	handler := NewScheduleHandler(service)

	router := gin.New()
	api := router.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Next()
	})
	handler.RegisterRoutes(api)

	return handler, router
}

// ---------------------------------------------------------------
// RunNow - additional coverage
// ---------------------------------------------------------------

func TestScheduleHandler_RunNow_Success(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:         "sched-1",
		Name:       "Run Now Test",
		ScenarioID: "scenario-1",
		Status:     entity.ScheduleStatusActive,
	}
	_, router := setupScheduleHandlerWithExecService(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/run", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var run entity.ScheduleRun
	if err := json.Unmarshal(w.Body.Bytes(), &run); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Execution failed (scenario repo error), so the run should be marked as failed
	if run.Status != "failed" {
		t.Errorf("run.Status = %q, want %q", run.Status, "failed")
	}
	if run.ScheduleID != "sched-1" {
		t.Errorf("run.ScheduleID = %q, want %q", run.ScheduleID, "sched-1")
	}
}

func TestScheduleHandler_RunNow_FindError(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.findErr = errors.New("database error")
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/run", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestScheduleHandler_RunNow_CreateRunError(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:         "sched-1",
		Name:       "Run Now Test",
		ScenarioID: "scenario-1",
		Status:     entity.ScheduleStatusActive,
	}
	repo.createRunErr = errors.New("failed to save run")
	_, router := setupScheduleHandlerWithExecService(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/run", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ---------------------------------------------------------------
// Pause - additional error paths
// ---------------------------------------------------------------

func TestScheduleHandler_Pause_FindError(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.findErr = errors.New("database error")
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/pause", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestScheduleHandler_Pause_UpdateError(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:     "sched-1",
		Status: entity.ScheduleStatusActive,
	}
	repo.updateErr = errors.New("update failed")
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/pause", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ---------------------------------------------------------------
// Resume - additional error paths
// ---------------------------------------------------------------

func TestScheduleHandler_Resume_FindError(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.findErr = errors.New("database error")
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/resume", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestScheduleHandler_Resume_UpdateError(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:     "sched-1",
		Status: entity.ScheduleStatusPaused,
	}
	repo.updateErr = errors.New("update failed")
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules/sched-1/resume", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ---------------------------------------------------------------
// Update - additional error paths
// ---------------------------------------------------------------

func TestScheduleHandler_Update_FindError(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.findErr = errors.New("database error")
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "Updated Name",
		ScenarioID: "scenario-1",
		Frequency:  "daily",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/schedules/sched-1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestScheduleHandler_Update_UpdateError(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.schedules["sched-1"] = &entity.Schedule{
		ID:     "sched-1",
		Name:   "Old Name",
		Status: entity.ScheduleStatusActive,
	}
	repo.updateErr = errors.New("update failed")
	_, router := setupRealScheduleHandler(repo)

	body := CreateScheduleRequest{
		Name:       "Updated Name",
		ScenarioID: "scenario-1",
		Frequency:  "daily",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/schedules/sched-1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ---------------------------------------------------------------
// GetRuns - limit edge cases
// ---------------------------------------------------------------

func TestScheduleHandler_GetRuns_LimitZero(t *testing.T) {
	repo := newMockScheduleRepo()
	for i := 0; i < 5; i++ {
		repo.runs["sched-1"] = append(repo.runs["sched-1"], &entity.ScheduleRun{
			ID:         "run-" + string(rune('1'+i)),
			ScheduleID: "sched-1",
		})
	}
	_, router := setupRealScheduleHandler(repo)

	// limit=0 is not > 0, so the handler should use the default of 20
	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1/runs?limit=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var runs []*entity.ScheduleRun
	if err := json.Unmarshal(w.Body.Bytes(), &runs); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// All 5 runs should be returned (default limit 20 > 5)
	if len(runs) != 5 {
		t.Errorf("len(runs) = %d, want 5", len(runs))
	}
}

func TestScheduleHandler_GetRuns_LimitNegative(t *testing.T) {
	repo := newMockScheduleRepo()
	for i := 0; i < 3; i++ {
		repo.runs["sched-1"] = append(repo.runs["sched-1"], &entity.ScheduleRun{
			ID:         "run-" + string(rune('1'+i)),
			ScheduleID: "sched-1",
		})
	}
	_, router := setupRealScheduleHandler(repo)

	// limit=-1 is not > 0, so the handler should use the default of 20
	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1/runs?limit=-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var runs []*entity.ScheduleRun
	if err := json.Unmarshal(w.Body.Bytes(), &runs); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// All 3 runs should be returned (default limit 20 > 3)
	if len(runs) != 3 {
		t.Errorf("len(runs) = %d, want 3", len(runs))
	}
}

func TestScheduleHandler_GetRuns_LimitExceedsMax(t *testing.T) {
	repo := newMockScheduleRepo()
	for i := 0; i < 3; i++ {
		repo.runs["sched-1"] = append(repo.runs["sched-1"], &entity.ScheduleRun{
			ID:         "run-" + string(rune('1'+i)),
			ScheduleID: "sched-1",
		})
	}
	_, router := setupRealScheduleHandler(repo)

	// limit=200 exceeds max of 100, so the handler should use the default of 20
	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1/runs?limit=200", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	var runs []*entity.ScheduleRun
	if err := json.Unmarshal(w.Body.Bytes(), &runs); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// All 3 runs should be returned (default limit 20 > 3)
	if len(runs) != 3 {
		t.Errorf("len(runs) = %d, want 3", len(runs))
	}
}

// ---------------------------------------------------------------
// Delete - additional coverage
// ---------------------------------------------------------------

func TestScheduleHandler_GetRuns_Error(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.findRunsErr = errors.New("database error")
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1/runs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestScheduleHandler_GetRuns_LimitBoundary100(t *testing.T) {
	repo := newMockScheduleRepo()
	repo.runs["sched-1"] = []*entity.ScheduleRun{
		{ID: "run-1", ScheduleID: "sched-1"},
	}
	_, router := setupRealScheduleHandler(repo)

	// limit=100 is the max valid value
	req := httptest.NewRequest(http.MethodGet, "/api/v1/schedules/sched-1/runs?limit=100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestScheduleHandler_Delete_NotFound_Succeeds(t *testing.T) {
	// Delete with a nonexistent ID does not trigger a "not found" error
	// because the Delete handler calls scheduleService.Delete which calls
	// repo.Delete directly without checking existence first.
	// If repo.Delete returns nil for a missing ID, the handler returns 200.
	repo := newMockScheduleRepo()
	_, router := setupRealScheduleHandler(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/schedules/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
}
