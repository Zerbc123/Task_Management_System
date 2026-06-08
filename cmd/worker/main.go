package main

import (
	"context"
	"log"

	"task-management/internal/database"
	"task-management/internal/notification"
)

func main() {
	ctx := context.Background()

	redisClient := database.ConnectRedis()

	notificationQueue := notification.NewQueue(redisClient)

	worker := notification.NewWorker(notificationQueue)

	log.Println("worker process running")

	worker.Start(ctx)
}