package handler

import (
	"net/http"
	"task-management/internal/audit/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuditHandler struct {
	service services.AuditService
}

func NewAuditHandler(service services.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) RegisterRoutes(r gin.IRouter) {
	r.GET("/tasks/:id/audit-logs", h.GetTaskAuditLogs)
}

func (h *AuditHandler) GetTaskAuditLogs(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	logs, err := h.service.GetByEntityID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}
