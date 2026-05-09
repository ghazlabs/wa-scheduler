package core

import (
	"encoding/json"
)

type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusFailed    MessageStatus = "failed"
	MessageStatusScheduled MessageStatus = "scheduled"
)

type Message struct {
	ID                 string        `json:"id"`
	Content            string        `json:"content"`
	RecipientNumbers   []string      `json:"recipient_numbers"`
	ScheduledSendingAt int64         `json:"scheduled_sending_at"`
	RetriedCount       int           `json:"retried_count"`
	Status             MessageStatus `json:"status"`
	SentAt             *int64        `json:"sent_at"`
	Reason             *string       `json:"reason"`
	CreatedAt          int64         `json:"created_at"`
	UpdatedAt          int64         `json:"updated_at"`
}

func (m *Message) String() string {
	jsonData, _ := json.Marshal(m)
	return string(jsonData)
}

type ScheduleMessageInput struct {
	Content            string   `json:"content"`
	RecipientNumbers   []string `json:"recipient_numbers"`
	ScheduledSendingAt int64    `json:"scheduled_sending_at"`
}

type GetAllMessagesInput struct {
	Status MessageStatus `json:"status"`
}

type RetryMessageInput struct {
	ID                 string `json:"id"`
	ScheduledSendingAt int64  `json:"scheduled_sending_at"`
}
