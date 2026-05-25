package core

import (
	"context"
	"errors"
	"testing"
)

type fakeStorage struct {
	message *Message
}

func (s *fakeStorage) GetAllMessages(ctx context.Context, input GetAllMessagesInput) ([]Message, error) {
	return nil, nil
}

func (s *fakeStorage) SaveMessage(ctx context.Context, message Message) error {
	return nil
}

func (s *fakeStorage) UpdateMessage(ctx context.Context, message Message) error {
	return nil
}

func (s *fakeStorage) GetMessage(ctx context.Context, id string) (*Message, error) {
	return s.message, nil
}

type fakeScheduler struct {
	retryCalled bool
	retryMsg    Message
}

func (s *fakeScheduler) ScheduleMessage(ctx context.Context, message Message) error {
	return nil
}

func (s *fakeScheduler) RetryMessage(ctx context.Context, message Message) error {
	s.retryCalled = true
	s.retryMsg = message
	return nil
}

func TestServiceRetryMessage(t *testing.T) {
	tests := []struct {
		name            string
		message         *Message
		wantErr         error
		wantRetryCalled bool
		wantScheduledAt int64
	}{
		{
			name: "failed message can be retried",
			message: &Message{
				ID:     "message-1",
				Status: MessageStatusFailed,
			},
			wantRetryCalled: true,
			wantScheduledAt: 12345,
		},
		{
			name: "sent message cannot be retried",
			message: &Message{
				ID:     "message-1",
				Status: MessageStatusSent,
			},
			wantErr: ErrRetryNonFailed,
		},
		{
			name: "scheduled message cannot be retried",
			message: &Message{
				ID:     "message-1",
				Status: MessageStatusScheduled,
			},
			wantErr: ErrRetryNonFailed,
		},
		{
			name:    "missing message returns not found",
			message: nil,
			wantErr: ErrMessageNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &fakeStorage{message: tt.message}
			scheduler := &fakeScheduler{}
			svc, err := NewService(ServiceConfig{
				Storage:   storage,
				Scheduler: scheduler,
			})
			if err != nil {
				t.Fatalf("NewService() error = %v", err)
			}

			err = svc.RetryMessage(context.Background(), RetryMessageInput{
				ID:                 "message-1",
				ScheduledSendingAt: 12345,
			})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("RetryMessage() error = %v, want %v", err, tt.wantErr)
			}
			if scheduler.retryCalled != tt.wantRetryCalled {
				t.Fatalf("RetryMessage() scheduler called = %v, want %v", scheduler.retryCalled, tt.wantRetryCalled)
			}
			if tt.wantRetryCalled && scheduler.retryMsg.ScheduledSendingAt != tt.wantScheduledAt {
				t.Fatalf("RetryMessage() scheduled at = %v, want %v", scheduler.retryMsg.ScheduledSendingAt, tt.wantScheduledAt)
			}
		})
	}
}
