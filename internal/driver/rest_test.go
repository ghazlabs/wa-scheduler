package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ghazlabs/wa-scheduler/internal/core"
	"github.com/go-chi/chi/v5"
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
		{ID: "test-1", Status: input.Status},
		{ID: "test-2", Status: input.Status},
	}, nil
}

func (m *mockService) SendMessage(ctx context.Context, input core.ScheduleMessageInput) error {
	return nil
}

func (m *mockService) RetryMessage(ctx context.Context, input core.RetryMessageInput) error {
	if input.ID == "200" {
		return core.ErrMessageNotFound
	}
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

func TestGetAllMessages(t *testing.T) {
	testCases := []struct {
		name              string
		query             string
		expectedStatus    int
		expectedOk        bool
		expectedMsgStatus string
		expectedError     bool
	}{
		{
			name:              "should return 200 with failed messages when status = failed",
			query:             "/messages?status=failed",
			expectedStatus:    http.StatusOK,
			expectedOk:        true,
			expectedMsgStatus: string(core.MessageStatusFailed),
		},
		{
			name:              "should return 200 with scheduled messages when status = scheduled",
			query:             "/messages?status=scheduled",
			expectedStatus:    http.StatusOK,
			expectedOk:        true,
			expectedMsgStatus: string(core.MessageStatusScheduled),
		},
		{
			name:              "should return 200 sent messages when status = sent",
			query:             "/messages?status=sent",
			expectedStatus:    http.StatusOK,
			expectedOk:        true,
			expectedMsgStatus: string(core.MessageStatusSent),
		},
		{
			name:           "should return 400 when status is invalid",
			query:          "/messages?status=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedOk:     false,
			expectedError:  true,
		},
		{
			name:           "should return 200 with all messages when no status filter",
			query:          "/messages",
			expectedStatus: http.StatusOK,
			expectedOk:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := newTestAPI()

			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			req.SetBasicAuth("admin", "admin")

			w := httptest.NewRecorder()

			api.serveGetMessages(w, req)

			body := parseBody(w)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedOk, body["ok"].(bool))

			if tc.expectedError {
				assert.NotNil(t, body["err"])
				return
			}

			data, ok := body["data"].([]interface{})
			assert.True(t, ok)
			assert.NotEmpty(t, data)

			if tc.expectedMsgStatus != "" {
				for _, item := range data {
					msg := item.(map[string]interface{})
					assert.Equal(t, tc.expectedMsgStatus, msg["status"])
				}
			}
		})
	}
}

func TestRetryMessage(t *testing.T) {
	tests := []struct {
		name           string
		messageID      string
		expectedStatus int
		expectedOK     bool
		expectedErr    string
		expectedMsg    string
	}{
		{
			name:           "message not found",
			messageID:      "200",
			expectedStatus: http.StatusNotFound,
			expectedOK:     false,
			expectedErr:    "ERR_MESSAGES_NOT_FOUND",
			expectedMsg:    "message id not found",
		},
		{
			name:           "message found",
			messageID:      "1",
			expectedStatus: http.StatusOK,
			expectedOK:     true,
			expectedErr:    "",
			expectedMsg:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newTestAPI()

			reqBody := RetryMessageRequest{
				ScheduledSendingAt: 1234567890,
			}

			jsonBody, _ := json.Marshal(reqBody)

			req := httptest.NewRequest(
				http.MethodPost,
				"/messages/"+tt.messageID+"/retry",
				bytes.NewBuffer(jsonBody),
			)

			req.SetBasicAuth("admin", "admin")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.messageID)

			req = req.WithContext(
				context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
			)

			w := httptest.NewRecorder()

			api.serveRetryMessage(w, req)

			body := parseBody(w)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedOK, body["ok"])

			if tt.expectedErr != "" {
				assert.Equal(t, tt.expectedErr, body["err"])
			}
			if tt.expectedMsg != "" {
				assert.Equal(t, tt.expectedMsg, body["msg"])
			}
		})
	}
}
