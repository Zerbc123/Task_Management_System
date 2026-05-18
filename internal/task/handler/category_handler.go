package handler

import(
	"errors"
	
	response "task-management/internal/response"
	"task-management/internal/task/dto"
	"task-management/internal/task/repository"
	"task-management/internal/task/services"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) RegisterRoutes(r *gin.Engine) {
	categories := r.Group("/categories")
	{
		categories.POST("/", h.CreateCategory)
		categories.GET("/", h.GetCategories)
		categories.GET("/:id", h.GetCategoryByID)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload")
		return
	}

	category, err := h.service.CreateCategory(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "Category created successfully", category)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, "Categories retrieved successfully", categories)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		handleCategoryError(c, err)
		return
	}
	response.OK(c, "Category retrieved successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := h.service.UpdateCategory(id, req)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	response.OK(c, "Category updated successfully", category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}
	if err := h.service.DeleteCategory(id); err != nil {
		handleCategoryError(c, err)
		return
	}
	response.OK(c, "Category deleted successfully", nil)
}

func handleCategoryError(c *gin.Context, err error) {
	if errors.Is(err, repository.ErrCategoryNotFound) {
		response.NotFound(c, err.Error())
		return
	}
	response.InternalServerError(c, err.Error())
}