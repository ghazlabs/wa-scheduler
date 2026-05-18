package driver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ghazlabs/wa-scheduler/internal/core"
)

// mockService adalah implementasi dari core.Service
type mockService struct{}

func (m *mockService) InitializeService(ctx context.Context) {}

func (m *mockService) GetAllMessages(ctx context.Context, input core.GetAllMessagesInput) ([]core.Message, error) {
	return []core.Message{}, nil
}

func (m *mockService) SendMessage(ctx context.Context, input core.ScheduleMessageInput) error {
	return nil
}

func (m *mockService) RetryMessage(ctx context.Context, input core.RetryMessageInput) error {
	return nil
}

// newTestAPI membuat instance API untuk test
func newTestAPI() *API {
	api, _ := NewAPI(APIConfig{
		Service:            &mockService{},
		ClientUsername:     "admin",
		ClientPassword:     "admin",
		WebClientPublicDir: ".",
	})
	return api
}

// Test validasi endpoint GET /messages dengan status=failed
// Expected result: response HTTP 200 OK
func TestGetMessages_FailedStatus_Returns200(t *testing.T) {
	api := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/messages?status=failed", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	api.serveGetMessages(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// Test validasi endpoint GET /messages dengan status=scheduled
// Expected result: response HTTP 200 OK
func TestGetMessages_ScheduledStatus_Returns200(t *testing.T) {
	api := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/messages?status=scheduled", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	api.serveGetMessages(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// Test validasi endpoint GET /messages dengan status=sent
// Expected result: response HTTP 200 OK
func TestGetMessages_SentStatus_Returns200(t *testing.T) {
	api := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/messages?status=sent", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	api.serveGetMessages(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// Test validasi endpoint GET /messages dengan status tidak valid
// Expected result: response HTTP 400 Bad Request
// dan body response mengandung pesan error
func TestGetMessages_InvalidStatus_Returns400(t *testing.T) {
	api := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/messages?status=invalid", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	api.serveGetMessages(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	if body["err"] == nil {
		t.Error("expected error message in response body")
	}
}

// Test validasi endpoint GET /messages tanpa parameter status
// Expected result: response HTTP 200 OK
func TestGetMessages_NoStatus_Returns200(t *testing.T) {
	api := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/messages", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	api.serveGetMessages(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}