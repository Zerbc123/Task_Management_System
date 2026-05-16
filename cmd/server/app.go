package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"task-management/internal/middleware"
	"task-management/internal/task/handler"
	"task-management/internal/task/repository"
	"task-management/internal/task/services"
)

func main() {
	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())

	taskRepo := repository.NewMemoryTaskRepository()
	taskService := services.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	taskHandler.RegisterRoutes(r)

	log.Println("Server running at http://localhost:8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}