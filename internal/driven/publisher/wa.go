package wa

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ghazlabs/wa-scheduler/internal/core"
	"github.com/go-resty/resty/v2"
	"gopkg.in/validator.v2"
)

type WaPublisher struct {
	WaPublisherConfig
}

func NewWaPublisher(cfg WaPublisherConfig) (*WaPublisher, error) {
	// validate config
	err := validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &WaPublisher{
		WaPublisherConfig: cfg,
	}, nil
}

type WaPublisherConfig struct {
	HttpClient   *resty.Client `validate:"nonnil"`
	Username     string        `validate:"nonzero"`
	Password     string        `validate:"nonzero"`
	WaApiBaseUrl string        `validate:"nonzero"`
}

func (n *WaPublisher) Publish(ctx context.Context, msg core.Message) error {
	for _, recID := range msg.RecipientNumbers {
		err := n.sendMessage(ctx, recID, msg.Content)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *WaPublisher) sendMessage(ctx context.Context, recID string, content string) error {
	// send notification to whatsapp
	var rsp RespSendMessage
	resp, err := n.HttpClient.R().
		SetContext(ctx).
		SetBasicAuth(n.Username, n.Password).
		SetError(&rsp).
		SetBody(map[string]interface{}{
			"phone":   recID,
			"message": content,
		}).
		Post(fmt.Sprintf("%v/send/message", n.WaApiBaseUrl))
	if err != nil {
		return fmt.Errorf("unable to make http request: %w", err)
	}
	if resp.IsError() {
		slog.Error("failed to send wa message", slog.String("response", rsp.String()))
		if rsp.IsSessionExpired() {
			return core.ErrSessionExpired
		}

		return fmt.Errorf("failed to send message: %s", resp.String())
	}

	return nil
}
