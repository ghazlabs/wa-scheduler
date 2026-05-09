package core

import (
	"context"
	"testing"
	"time"
)

func TestSendMessage(t *testing.T) {
	// 1. Setup your service (mocking dependencies if necessary)
	svc := &service{
		// Add your storage/scheduler mocks here
	}

	// 2. Define the Test Table
	tests := []struct {
		name          string
		input         ScheduleMessageInput
		expectedError string
	}{
		{
			name: "Empty Recipients",
			input: ScheduleMessageInput{
				RecipientNumbers:   []string{}, // Empty
				Content:            "Hello",
				ScheduledSendingAt: time.Now().Add(time.Hour).Unix(),
			},
			expectedError: "recipient numbers cannot be empty",
		},
		{
			name: "Empty Content",
			input: ScheduleMessageInput{
				RecipientNumbers:   []string{"12345"},
				Content:            "", // Empty
				ScheduledSendingAt: time.Now().Add(time.Hour).Unix(),
			},
			expectedError: "message content cannot be empty",
		},
		{
			name: "Past Timestamp",
			input: ScheduleMessageInput{
				RecipientNumbers:   []string{"12345"},
				Content:            "Hello",
				ScheduledSendingAt: time.Now().Add(-time.Hour).Unix(), // 1 hour ago
			},
			expectedError: "scheduled time must be in the future",
		},
		{
			name: "Zero Timestamp",
			input: ScheduleMessageInput{
				RecipientNumbers:   []string{"12345"},
				Content:            "Hello",
				ScheduledSendingAt: 0, // The "Zero" problem
			},
			expectedError: "scheduled time must be in the future",
		},
		{
			name: "Valid",
			input: ScheduleMessageInput{
				RecipientNumbers:   []string{"12345"},
				Content:            "Hello",
				ScheduledSendingAt: time.Now().Add(time.Hour).Unix(), // The "Zero" problem
			},
			expectedError: "",
		},
	}

	// 3. Run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.SendMessage(context.Background(), tc.input)

			if tc.expectedError == "" {
				if err != nil {
					t.Fatalf("expected no error, but got: %v", err)
				}
			}

			// Check if we got an error when we expected one
			if err == nil {
				t.Fatal("expected an error but got none")
			}

			// Check if the error message is what we expected
			if err.Error() != tc.expectedError {
				t.Errorf("expected error %q, got %q", tc.expectedError, err.Error())
			}
		})
	}
}
