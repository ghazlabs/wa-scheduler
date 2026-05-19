package driver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ghazlabs/wa-scheduler/internal/core"
	"github.com/stretchr/testify/assert"
)

type mockService struct{}

func (m *mockService) InitializeService(ctx context.Context) {}

func (m *mockService) GetAllMessages(ctx context.Context, input core.GetAllMessagesInput) ([]core.Message, error) {
	if input.Status == "" {
		return []core.Message{
			{ID: "test-1", Status: core.MessageStatusFailed},
			{ID: "test-2", Status: core.MessageStatusSent},
			{ID: "test-3", Status: core.MessageStatusScheduled},
		}, nil
	}
	return []core.Message{
		{ID: "test-test-1", Status: input.Status},
	}, nil
}

func (m *mockService) SendMessage(ctx context.Context, input core.ScheduleMessageInput) error {
	return nil
}

func (m *mockService) RetryMessage(ctx context.Context, input core.RetryMessageInput) error {
	return nil
}

func newTestAPI() *API {
	api, _ := NewAPI(APIConfig{
		Service:            &mockService{},
		ClientUsername:     "admin",
		ClientPassword:     "admin",
		WebClientPublicDir: ".",
	})
	return api
}

func parseBody(w *httptest.ResponseRecorder) map[string]interface{} {
	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	return body
}

func TestGetMessages_FailedStatus_Returns200WithFailedMessages(t *testing.T) {
	api := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/messages?status=failed", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	api.serveGetMessages(w, req)

	body := parseBody(w)
	data := body["data"].([]interface{})
	firstMessage := data[0].(map[string]interface{})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, body["ok"].(bool))
	assert.NotEmpty(t, data)
	assert.Equal(t, string(core.MessageStatusFailed), firstMessage["status"])
}

func TestGetMessages_ScheduledStatus_Returns200(t *testing.T) {
	api := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/messages?status=scheduled", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	api.serveGetMessages(w, req)

	body := parseBody(w)
	data := body["data"].([]interface{})
	firstMessage := data[0].(map[string]interface{})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, data)
	assert.Equal(t, string(core.MessageStatusScheduled), firstMessage["status"])
}

func TestGetMessages_SentStatus_Returns200(t *testing.T) {
	api := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/messages?status=sent", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	api.serveGetMessages(w, req)

	body := parseBody(w)
	data := body["data"].([]interface{})
	firstMessage := data[0].(map[string]interface{})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, data)
	assert.Equal(t, string(core.MessageStatusSent), firstMessage["status"])
}

func TestGetMessages_InvalidStatus_Returns400(t *testing.T) {
	api := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/messages?status=invalid", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	api.serveGetMessages(w, req)

	body := parseBody(w)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, body["ok"].(bool))
	assert.Equal(t, "invalid status", body["msg"])
}

func TestGetMessages_NoStatus_Returns200(t *testing.T) {
	api := newTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/messages", nil)
	req.SetBasicAuth("admin", "admin")
	w := httptest.NewRecorder()

	api.serveGetMessages(w, req)
	body := parseBody(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, body["ok"].(bool))
}
