package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"task-management/internal/project/dto"
	"task-management/internal/project/repository"
	"task-management/internal/project/services"
	response "task-management/internal/shared/response"
)

type ProjectHandler struct {
	service services.ProjectService
}

func NewProjectHandler(service services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		service: service,
	}
}

func (h *ProjectHandler) RegisterRoutes(r gin.IRouter) {
	projects := r.Group("/projects")
	{
		projects.POST("", h.CreateProject)
		projects.GET("", h.GetProjects)
		projects.GET("/:id", h.GetProjectByID)
		projects.PUT("/:id", h.UpdateProject)
		projects.DELETE("/:id", h.DeleteProject)
	}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	project, err := h.service.CreateProject(req, userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, "project created successfully", project)
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	projects, err := h.service.GetProjects()
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, "projects retrieved successfully", projects)
}

func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project id")
		return
	}

	project, err := h.service.GetProjectByID(id)
	if err != nil {
		handleProjectError(c, err)
		return
	}

	response.OK(c, "project retrieved successfully", project)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user id")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project id")
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	project, err := h.service.UpdateProject(id, userID, req)
	if err != nil {
		handleProjectError(c, err)
		return
	}

	response.OK(c, "project updated successfully", project)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user id")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project id")
		return
	}

	if err := h.service.DeleteProject(id, userID); err != nil {
		handleProjectError(c, err)
		return
	}

	response.OK(c, "project deleted successfully", nil)
}

func handleProjectError(c *gin.Context, err error) {
	if errors.Is(err, repository.ErrProjectNotFound) {
		response.NotFound(c, err.Error())
		return
	}

	if errors.Is(err, services.ErrForbidden) {
		response.Forbidden(c, err.Error())
		return
	}

	response.InternalServerError(c, err.Error())
}

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := value.(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}

	return userID, true
}