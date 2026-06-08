package websocket

import "github.com/google/uuid"

type TaskUpdatedEvent struct {
	Type   string    `json:"type"`
	TaskID uuid.UUID `json:"task_id"`
	Title  string    `json:"title"`
	Status string    `json:"status"`
}

type CommentCreatedEvent struct {
	Type      string    `json:"type"`
	CommentID uuid.UUID `json:"comment_id"`
	TaskID    uuid.UUID `json:"task_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
}
