package main

import (
	"log"

	"github.com/gin-gonic/gin"

	categoryHandler "task-management/internal/task/handler"
	categoryRepository "task-management/internal/task/repository"
	categoryServices "task-management/internal/task/services"

	"task-management/internal/middleware"

	taskHandler "task-management/internal/task/handler"
	taskRepository "task-management/internal/task/repository"
	taskServices "task-management/internal/task/services"
)

func main() {
	// Create Gin server
	r := gin.New()

	// Middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())

	// =========================
	// TASK MODULE
	// =========================

	taskRepo := taskRepository.NewMemoryTaskRepository()

	taskService := taskServices.NewTaskService(taskRepo)

	taskHandler := taskHandler.NewTaskHandler(taskService)

	taskHandler.RegisterRoutes(r)

	// =========================
	// CATEGORY MODULE
	// =========================

	categoryRepo := categoryRepository.NewMemoryCategoryRepository()

	categoryService := categoryServices.NewCategoryService(categoryRepo)

	categoryHandler := categoryHandler.NewCategoryHandler(categoryService)

	categoryHandler.RegisterRoutes(r)

	// Health Check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Task Management API is running",
		})
	})

	log.Println("Server running at http://localhost:8080")

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}