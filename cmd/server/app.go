package main

import (
	"log"

	"github.com/gin-gonic/gin"

	authHandler "task-management/internal/auth/handler"
	authServices "task-management/internal/auth/services"

	"task-management/internal/database"
	"task-management/internal/middleware"

	projectHandler "task-management/internal/project/handler"
	projectRepository "task-management/internal/project/repository"
	projectServices "task-management/internal/project/services"

	taskHandler "task-management/internal/task/handler"
	taskRepository "task-management/internal/task/repository"
	taskServices "task-management/internal/task/services"

	userRepository "task-management/internal/user/repository"
)

func main() {
	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())

	db := database.ConnectDB()

	userRepo := userRepository.NewGormUserRepository(db)

	authService := authServices.NewAuthService(userRepo)
	authHandler := authHandler.NewAuthHandler(authService)
	authHandler.RegisterRoutes(r)

	taskRepo := taskRepository.NewGormTaskRepository(db)
	taskService := taskServices.NewTaskService(taskRepo)
	taskHandler := taskHandler.NewTaskHandler(taskService)

	projectRepo := projectRepository.NewGormProjectRepository(db)
	projectService := projectServices.NewProjectService(projectRepo)
	projectHandler := projectHandler.NewProjectHandler(projectService)

	protected := r.Group("/")
	protected.Use(middleware.JWTAuth())

	projectHandler.RegisterRoutes(protected)
	taskHandler.RegisterRoutes(protected)

	log.Println("server running at http://localhost:8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}