package driver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ghazlabs/wa-scheduler/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"gopkg.in/validator.v2"
)

type API struct {
	APIConfig
}

type APIConfig struct {
	Service            core.Service `validate:"nonnil"`
	ClientUsername     string       `validate:"nonzero"`
	ClientPassword     string       `validate:"nonzero"`
	WebClientPublicDir string       `validate:"nonzero"`
	DefaultNumbers     []string
}

func NewAPI(cfg APIConfig) (*API, error) {
	err := validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid API config: %w", err)
	}
	return &API{APIConfig: cfg}, nil
}

func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", a.serveWebFrontend)
	r.Group(func(r chi.Router) {
		r.Use(BasicAuth(a.ClientUsername, a.ClientPassword))
		r.Use(render.SetContentType(render.ContentTypeJSON))

		r.Get("/check", a.serveCheckSystem)
		r.Get("/messages", a.serveGetMessages)
		r.Post("/messages", a.serveScheduleMessage)
		r.Post("/messages/{id}/retry", a.serveRetryMessage)
	})

	return r
}

func (a *API) serveWebFrontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(a.WebClientPublicDir, "index.html"))
}

func (a *API) serveCheckSystem(w http.ResponseWriter, r *http.Request) {
	resp := NewSuccessResp(RespCheck{
		DefaultNumbers: a.DefaultNumbers,
	})
	render.Render(w, r, resp)
}

func (a *API) serveGetMessages(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	input := core.GetAllMessagesInput{}
	if status != "" {
		switch status {
		case string(core.MessageStatusScheduled):
			input.Status = core.MessageStatusScheduled
		case string(core.MessageStatusFailed):
			input.Status = core.MessageStatusFailed
		case string(core.MessageStatusSent):
			input.Status = core.MessageStatusSent
		default:
			render.Render(w, r, NewErrorResp(NewBadRequestError("invalid status")))
			return
		}
	}

	messages, err := a.Service.GetAllMessages(r.Context(), input)
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}

	resp := NewSuccessResp(messages)
	render.Render(w, r, resp)
}

func (a *API) serveScheduleMessage(w http.ResponseWriter, r *http.Request) {
	var req SendMessageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}

	err = a.Service.SendMessage(r.Context(), core.ScheduleMessageInput{
		Content:            req.Content,
		RecipientNumbers:   req.RecipientNumbers,
		ScheduledSendingAt: req.ScheduledSendingAt,
	})
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}

	resp := NewSuccessResp(nil)
	render.Render(w, r, resp)
}

func (a *API) serveRetryMessage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Render(w, r, NewErrorResp(NewBadRequestError("id is required")))
		return
	}

	var req RetryMessageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}
	if req.ScheduledSendingAt == 0 {
		req.ScheduledSendingAt = time.Now().Unix()
	}

	err = a.Service.RetryMessage(r.Context(), core.RetryMessageInput{
		ID:                 id,
		ScheduledSendingAt: req.ScheduledSendingAt,
	})
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}

	resp := NewSuccessResp(nil)
	render.Render(w, r, resp)
}
