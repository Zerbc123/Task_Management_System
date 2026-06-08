package notification

import "github.com/google/uuid"

type NotificationJob struct {
	TaskID     uuid.UUID  `json:"task_id"`
	UserID     *uuid.UUID `json:"user_id"`
	EventType  string     `json:"event_type"`
	Message    string     `json:"message"`
	RetryCount int        `json:"retry_count"`
	MaxRetry   int        `json:"max_retry"`
}
