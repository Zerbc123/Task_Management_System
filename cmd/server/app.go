package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"task-management/internal/cache"
	"task-management/internal/database"
	"task-management/internal/health"
	"task-management/internal/middleware"
	"task-management/internal/notification"

	authHandler "task-management/internal/auth/handler"
	authServices "task-management/internal/auth/services"

	commentHandler "task-management/internal/comment/handler"
	commentRepository "task-management/internal/comment/repository"
	commentServices "task-management/internal/comment/services"

	projectHandler "task-management/internal/project/handler"
	projectRepository "task-management/internal/project/repository"
	projectServices "task-management/internal/project/services"

	taskHandler "task-management/internal/task/handler"
	taskRepository "task-management/internal/task/repository"
	taskServices "task-management/internal/task/services"

	userRepository "task-management/internal/user/repository"

	websocketHandler "task-management/internal/websocket"
)

func main() {
	ctx := context.Background()

	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())

	wsHub := websocketHandler.NewHub()
	wsHandler := websocketHandler.NewHandler(wsHub)
	wsHandler.RegisterRoutes(r)

	db := database.ConnectDB()

	redisClient := database.ConnectRedis()
	redisCache := cache.NewRedisCache(redisClient)

	healthHandler := health.NewHandler(db, redisClient)
	healthHandler.RegisterRoutes(r)

	notificationQueue := notification.NewQueue(redisClient)

	// Nếu đã chạy worker bằng cmd/worker hoặc Docker service task-worker,
	// thì KHÔNG nên start worker trong API nữa.
	notificationWorker := notification.NewWorker(notificationQueue)
	go notificationWorker.Start(ctx)

	userRepo := userRepository.NewGormUserRepository(db)
	authService := authServices.NewAuthService(userRepo)
	authHandler := authHandler.NewAuthHandler(authService)
	authHandler.RegisterRoutes(r)

	taskRepo := taskRepository.NewGormTaskRepository(db)
	taskService := taskServices.NewTaskService(
		taskRepo,
		redisCache,
		notificationQueue,
		wsHub,
	)
	taskHandler := taskHandler.NewTaskHandler(taskService)

	projectRepo := projectRepository.NewGormProjectRepository(db)
	projectService := projectServices.NewProjectService(projectRepo)
	projectHandler := projectHandler.NewProjectHandler(projectService)

	commentRepo := commentRepository.NewGormCommentRepository(db)
	commentService := commentServices.NewCommentService(commentRepo, wsHub)
	commentHandler := commentHandler.NewCommentHandler(commentService)

	protected := r.Group("/")
	protected.Use(middleware.JWTAuth())

	projectHandler.RegisterRoutes(protected)
	taskHandler.RegisterRoutes(protected)
	commentHandler.RegisterRoutes(protected)

	log.Println("server running at http://localhost:8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}