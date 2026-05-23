package handler

import (
	"github.com/gin-gonic/gin"
	authServices "task-management/internal/auth/services"
	response "task-management/internal/shared/response"
	"task-management/internal/auth/dto"
	userRepository "task-management/internal/user/repository"
	"errors"
)

type AuthHandler struct {
	service authServices.AuthService
}

func NewAuthHandler(service authServices.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, userRepository.ErrEmailAlreadyExists) {
			response.BadRequest(c, err.Error())
			return
		}

		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, "register successfully", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	token, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, authServices.ErrInvalidCredentials) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, "login successfully", token)
}