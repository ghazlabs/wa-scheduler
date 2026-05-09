package core

import (
	"context"
)

type Storage interface {
	GetAllMessages(ctx context.Context, message GetAllMessagesInput) ([]Message, error)
	SaveMessage(ctx context.Context, message Message) error
	UpdateMessage(ctx context.Context, message Message) error
	GetMessage(ctx context.Context, id string) (*Message, error)
}

type Scheduler interface {
	ScheduleMessage(ctx context.Context, message Message) error
	RetryMessage(ctx context.Context, msg Message) error
}
