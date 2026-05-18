package services

import (
	"testing"
	"task-management/internal/task/repository"
	"task-management/internal/task/dto"
	"github.com/google/uuid"
	"task-management/internal/task/model"
)

func TestCreateTask(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	service := NewTaskService(repo)

	req := dto.CreateTaskRequest{
		Title:       "Learn Go",
		Description: "Study unit test",
		Status:      model.StatusTodo,
		Priority:    3,
		AssignedTo:  "Duy",
	}

	task, err := service.CreateTask(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if task.ID == uuid.Nil {
		t.Fatalf("Expected task ID to be generated")
	}

	if task.Title != req.Title {
		t.Errorf("Expected title %s, got %s", req.Title, task.Title)
	}

	if task.Status != model.StatusTodo{
		t.Fatalf("Expected status TODO, got %s", task.Status)
	}
}

func TestGetTasks(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	service := NewTaskService(repo)

	_, err := service.CreateTask(dto.CreateTaskRequest{
		Title:       "Task 1",
		Description: "Description 1",
		Status:      model.StatusTodo,
		Priority:    1,
	})

	if err != nil {
		t.Fatal(err)
	}

	tasks, err := service.GetTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}
}

func TestUpdateTask(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	service := NewTaskService(repo)

	task, err := service.CreateTask(dto.CreateTaskRequest{
		Title:       "Old title",
		Description: "Old description",
		Status:      model.StatusTodo,
		Priority:    1,
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedTask, err := service.UpdateTask(task.ID, dto.UpdateTaskRequest{
		Title:       "New title",
		Description: "New description",
		Status:      model.StatusDone,
		Priority:    3,
	})
	if err != nil {
		t.Fatal(err)
	}

	if updatedTask.Title != "New title" {
		t.Fatalf("Expected title 'New title', got '%s'", updatedTask.Title)
	}

	if updatedTask.Status != model.StatusDone {
		t.Fatalf("Expected status DONE, got %s", updatedTask.Status)
	}
}

func TestDeleteTask(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	service := NewTaskService(repo)

	task, err := service.CreateTask(dto.CreateTaskRequest{
		Title:       "Task to delete",
		Description: "Delete this task",
		Status:      model.StatusTodo,
		Priority:    2,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = service.DeleteTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.GetTaskByID(task.ID)
	if err == nil {
		t.Fatalf("Expected error when getting deleted task")
	}
}