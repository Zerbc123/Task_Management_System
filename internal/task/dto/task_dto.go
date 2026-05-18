package dto

import "task-management/internal/task/model"

type CreateTaskRequest struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      model.TaskStatus `json:"status"`
	Priority    int              `json:"priority"`

	AssignedTo string `json:"assigned_to"`
}

type UpdateTaskRequest struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      model.TaskStatus `json:"status"`
	Priority    int              `json:"priority"`

	AssignedTo string `json:"assigned_to"`
}