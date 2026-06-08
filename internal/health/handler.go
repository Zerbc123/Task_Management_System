package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHandler(
	db *gorm.DB,
	redis *redis.Client,
) *Handler {
	return &Handler{
		db:    db,
		redis: redis,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/health", h.Health)
}

func (h *Handler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	dbStatus := "UP"
	redisStatus := "UP"

	if err := h.db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		dbStatus = "DOWN"
	}

	if err := h.redis.Ping(ctx).Err(); err != nil {
		redisStatus = "DOWN"
	}

	status := "UP"

	if dbStatus == "DOWN" || redisStatus == "DOWN" {
		status = "DOWN"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   status,
		"database": dbStatus,
		"redis":    redisStatus,
	})
}