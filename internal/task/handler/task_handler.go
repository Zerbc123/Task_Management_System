package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"task-management/internal/task/dto"
	"task-management/internal/task/repository"
	"task-management/internal/task/services"

	response "task-management/internal/response"
)

type TaskHandler struct {
	service services.TaskService
}

func NewTaskHandler(service services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) RegisterRoutes(r *gin.Engine) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.CreateTask(req)
	if err != nil {
		response.InternalServerError(c, "Failed to create task")
		return
	}

	response.Created(c, "Task created successfully", task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.service.GetTasks()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve tasks")
		return
	}

	response.OK(c, "Tasks retrieved successfully", tasks)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, "Task retrieved successfully", task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	task, err := h.service.UpdateTask(id, req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, "Task updated successfully", task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	err = h.service.DeleteTask(id)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, "Task deleted successfully", nil)
}

func handleError(c *gin.Context, err error) {
	if errors.Is(err, repository.ErrTaskNotFound) {
		response.NotFound(c, "Task not found")
		return
	}

	response.InternalServerError(c, "Failed to perform operation")
}