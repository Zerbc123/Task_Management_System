package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"task-management/internal/shared/response"
	"task-management/internal/task/dto"
	"task-management/internal/task/repository"
	"task-management/internal/task/services"
)

type TaskHandler struct {
	service services.TaskService
}

func NewTaskHandler(service services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) RegisterRoutes(r gin.IRouter) {
	tasks := r.Group("/tasks")
	{
		tasks.POST("", h.CreateTask)
		tasks.GET("", h.GetTasks)
		tasks.GET("/:id", h.GetTaskByID)
		tasks.PUT("/:id", h.UpdateTask)
		tasks.DELETE("/:id", h.DeleteTask)
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	task, err := h.service.CreateTask(req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, "task created successfully", task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.service.GetTasks()
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, "tasks retrieved successfully", tasks)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		handleTaskError(c, err)
		return
	}

	response.OK(c, "task retrieved successfully", task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	task, err := h.service.UpdateTask(id, req)
	if err != nil {
		handleTaskError(c, err)
		return
	}

	response.OK(c, "task updated successfully", task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	if err := h.service.DeleteTask(id); err != nil {
		handleTaskError(c, err)
		return
	}

	response.OK(c, "task deleted successfully", nil)
}

func handleTaskError(c *gin.Context, err error) {
	if errors.Is(err, repository.ErrTaskNotFound) {
		response.NotFound(c, err.Error())
		return
	}

	response.InternalServerError(c, err.Error())
}