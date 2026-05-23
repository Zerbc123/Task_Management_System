package services

import (
	"testing"

	"github.com/google/uuid"

	"task-management/internal/task/dto"
	"task-management/internal/task/model"
	memoryRepo "task-management/internal/task/repository/memory"
)

func TestCreateTask(t *testing.T) {
	repo := memoryRepo.NewMemoryTaskRepository()
	service := NewTaskService(repo)

	projectID := uuid.New()
	assigneeID := uuid.New()

	req := dto.CreateTaskRequest{
		ProjectID:   projectID,
		Title:       "Learn Go",
		Description: "Study unit test",
		Status:      model.StatusTodo,
		AssigneeID:  &assigneeID,
	}

	task, err := service.CreateTask(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID == uuid.Nil {
		t.Fatalf("expected task ID to be generated")
	}

	if task.ProjectID != projectID {
		t.Fatalf("expected project ID %s, got %s", projectID, task.ProjectID)
	}

	if task.Title != req.Title {
		t.Errorf("expected title %s, got %s", req.Title, task.Title)
	}

	if task.Status != model.StatusTodo {
		t.Fatalf("expected status TODO, got %s", task.Status)
	}

	if task.AssigneeID == nil {
		t.Fatalf("expected assignee ID, got nil")
	}

	if *task.AssigneeID != assigneeID {
		t.Fatalf("expected assignee ID %s, got %s", assigneeID, *task.AssigneeID)
	}
}

func TestGetTasks(t *testing.T) {
	repo := memoryRepo.NewMemoryTaskRepository()
	service := NewTaskService(repo)

	_, err := service.CreateTask(dto.CreateTaskRequest{
		ProjectID:   uuid.New(),
		Title:       "Task 1",
		Description: "Description 1",
		Status:      model.StatusTodo,
	})
	if err != nil {
		t.Fatal(err)
	}

	tasks, err := service.GetTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
}

func TestUpdateTask(t *testing.T) {
	repo := memoryRepo.NewMemoryTaskRepository()
	service := NewTaskService(repo)

	oldProjectID := uuid.New()
	newProjectID := uuid.New()
	assigneeID := uuid.New()

	task, err := service.CreateTask(dto.CreateTaskRequest{
		ProjectID:   oldProjectID,
		Title:       "Old title",
		Description: "Old description",
		Status:      model.StatusTodo,
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedTask, err := service.UpdateTask(
		task.ID,
		dto.UpdateTaskRequest{
			ProjectID:   newProjectID,
			Title:       "New title",
			Description: "New description",
			Status:      model.StatusDone,
			AssigneeID:  &assigneeID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if updatedTask.ProjectID != newProjectID {
		t.Fatalf("expected project ID %s, got %s", newProjectID, updatedTask.ProjectID)
	}

	if updatedTask.Title != "New title" {
		t.Fatalf("expected title 'New title', got '%s'", updatedTask.Title)
	}

	if updatedTask.Status != model.StatusDone {
		t.Fatalf("expected status DONE, got %s", updatedTask.Status)
	}

	if updatedTask.AssigneeID == nil {
		t.Fatalf("expected assignee ID, got nil")
	}

	if *updatedTask.AssigneeID != assigneeID {
		t.Fatalf("expected assignee ID %s, got %s", assigneeID, *updatedTask.AssigneeID)
	}
}

func TestDeleteTask(t *testing.T) {
	repo := memoryRepo.NewMemoryTaskRepository()
	service := NewTaskService(repo)

	task, err := service.CreateTask(dto.CreateTaskRequest{
		ProjectID:   uuid.New(),
		Title:       "Task to delete",
		Description: "Delete this task",
		Status:      model.StatusTodo,
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
		t.Fatalf("expected error when getting deleted task")
	}
}