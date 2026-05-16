package dto

import "task-management/internal/task/model"

type CreateTaskRequest struct {
	Title       string           `json:"title" binding:"required"`
	Description string           `json:"description"`
	Status      model.TaskStatus `json:"status"`
	Priority    int              `json:"priority"`
}

type UpdateTaskRequest struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      model.TaskStatus `json:"status"`
	Priority    int              `json:"priority"`
}