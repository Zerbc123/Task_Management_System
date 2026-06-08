package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	commentModel "task-management/internal/comment/model"
	projectModel "task-management/internal/project/model"
	taskModel "task-management/internal/task/model"
	userModel "task-management/internal/user/model"
)

func ConnectDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	err = db.AutoMigrate(
		&userModel.User{},
		&projectModel.Project{},
		&taskModel.Task{},
		&commentModel.Comment{},
	)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	log.Println("database connected")
	log.Println("database migrated")

	return db
}
