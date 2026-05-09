package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gopkg.in/validator.v2"
)

type Service interface {
	InitializeService(ctx context.Context)
	GetAllMessages(ctx context.Context, input GetAllMessagesInput) ([]Message, error)
	SendMessage(ctx context.Context, inputMsg ScheduleMessageInput) error
	RetryMessage(ctx context.Context, input RetryMessageInput) error
}

type ServiceConfig struct {
	Storage   Storage   `validate:"nonnil"`
	Scheduler Scheduler `validate:"nonnil"`
}

type service struct {
	ServiceConfig
}

func NewService(config ServiceConfig) (Service, error) {
	err := validator.Validate(config)
	if err != nil {
		return nil, err
	}

	return &service{
		ServiceConfig: config,
	}, nil
}

func (s *service) InitializeService(ctx context.Context) {
	scheduledMessages, err := s.Storage.GetAllMessages(ctx, GetAllMessagesInput{
		Status: MessageStatusScheduled,
	})
	if err != nil {
		fmt.Printf("failed to retrieve scheduled messages: %v\n", err)
		return
	}
	for _, msg := range scheduledMessages {
		err := s.Scheduler.ScheduleMessage(ctx, msg)
		if err != nil {
			fmt.Printf("failed to schedule message: %v\n", err)
			continue
		}
	}
	fmt.Println("Service initialized and scheduled messages processed.")
}

func (s *service) GetAllMessages(ctx context.Context, input GetAllMessagesInput) ([]Message, error) {
	messages, err := s.Storage.GetAllMessages(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}

	return messages, nil
}

func (s *service) SendMessage(ctx context.Context, input ScheduleMessageInput) error {
	msg := Message{
		ID:                 uuid.New().String(),
		Content:            input.Content,
		RecipientNumbers:   input.RecipientNumbers,
		ScheduledSendingAt: input.ScheduledSendingAt,
		Status:             MessageStatusScheduled,
	}
	err := s.Storage.SaveMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to schedule message: %w", err)
	}

	err = s.Scheduler.ScheduleMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to schedule message: %w", err)
	}

	return nil
}

func (s *service) RetryMessage(ctx context.Context, input RetryMessageInput) error {
	msg, err := s.Storage.GetMessage(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve message: %w", err)
	}
	if msg == nil {
		return ErrMessageNotFound
	}
	msg.ScheduledSendingAt = input.ScheduledSendingAt

	err = s.Scheduler.RetryMessage(ctx, *msg)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}
