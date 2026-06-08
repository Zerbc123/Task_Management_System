package dto

import (
	"task-management/internal/task/model"

	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	ProjectID   uuid.UUID        `json:"project_id" binding:"required"`
	Title       string           `json:"title" binding:"required"`
	Description string           `json:"description"`
	Status      model.TaskStatus `json:"status"`
	AssigneeID  *uuid.UUID       `json:"assignee_id"`
}

type UpdateTaskRequest struct {
	ProjectID   uuid.UUID        `json:"project_id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      model.TaskStatus `json:"status"`
	AssigneeID  *uuid.UUID       `json:"assignee_id"`
}
