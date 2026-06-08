package services

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"task-management/internal/notification"
	"task-management/internal/task/dto"
	"task-management/internal/task/model"
	memoryRepo "task-management/internal/task/repository/memory"
)

type mockCache struct {
	taskList map[string][]*model.Task
	tasks    map[uuid.UUID]*model.Task
}

func newMockCache() *mockCache {
	return &mockCache{
		taskList: make(map[string][]*model.Task),
		tasks:    make(map[uuid.UUID]*model.Task),
	}
}

func (m *mockCache) GetTaskList(ctx context.Context) ([]*model.Task, bool) {
	tasks, ok := m.taskList["tasks:all"]
	return tasks, ok
}

func (m *mockCache) SetTaskList(ctx context.Context, tasks []*model.Task) {
	m.taskList["tasks:all"] = tasks
}

func (m *mockCache) DeleteTaskList(ctx context.Context) {
	delete(m.taskList, "tasks:all")
}

func (m *mockCache) GetTaskByID(ctx context.Context, id uuid.UUID) (*model.Task, bool) {
	task, ok := m.tasks[id]
	return task, ok
}

func (m *mockCache) SetTaskByID(ctx context.Context, task *model.Task) {
	m.tasks[task.ID] = task
}

func (m *mockCache) DeleteTaskByID(ctx context.Context, id uuid.UUID) {
	delete(m.tasks, id)
}

type mockQueue struct {
	jobs []notification.NotificationJob
}

func newMockQueue() *mockQueue {
	return &mockQueue{
		jobs: make([]notification.NotificationJob, 0),
	}
}

func (m *mockQueue) Push(ctx context.Context, job notification.NotificationJob) error {
	m.jobs = append(m.jobs, job)
	return nil
}

func setupTaskService() (context.Context, TaskService, *mockQueue) {
	ctx := context.Background()

	repo := memoryRepo.NewMemoryTaskRepository()
	cache := newMockCache()
	queue := newMockQueue()

	service := NewTaskService(repo, cache, queue, nil)

	return ctx, service, queue
}

func TestCreateTask(t *testing.T) {
	ctx, service, queue := setupTaskService()

	projectID := uuid.New()
	assigneeID := uuid.New()

	req := dto.CreateTaskRequest{
		ProjectID:   projectID,
		Title:       "Learn Go",
		Description: "Study unit test",
		Status:      model.StatusTodo,
		AssigneeID:  &assigneeID,
	}

	task, err := service.Create(ctx, req)
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

	if len(queue.jobs) != 1 {
		t.Fatalf("expected 1 notification job, got %d", len(queue.jobs))
	}

	if queue.jobs[0].TaskID != task.ID {
		t.Fatalf("expected notification task ID %s, got %s", task.ID, queue.jobs[0].TaskID)
	}
}

func TestGetAll(t *testing.T) {
	ctx, service, _ := setupTaskService()

	_, err := service.Create(ctx, dto.CreateTaskRequest{
		ProjectID:   uuid.New(),
		Title:       "Task 1",
		Description: "Description 1",
		Status:      model.StatusTodo,
	})
	if err != nil {
		t.Fatal(err)
	}

	tasks, err := service.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
}

func TestGetByID(t *testing.T) {
	ctx, service, _ := setupTaskService()

	createdTask, err := service.Create(ctx, dto.CreateTaskRequest{
		ProjectID:   uuid.New(),
		Title:       "Task detail",
		Description: "Get task by ID",
		Status:      model.StatusTodo,
	})
	if err != nil {
		t.Fatal(err)
	}

	task, err := service.GetByID(ctx, createdTask.ID)
	if err != nil {
		t.Fatal(err)
	}

	if task.ID != createdTask.ID {
		t.Fatalf("expected task ID %s, got %s", createdTask.ID, task.ID)
	}
}

func TestUpdateTask(t *testing.T) {
	ctx, service, queue := setupTaskService()

	oldProjectID := uuid.New()
	newProjectID := uuid.New()
	assigneeID := uuid.New()

	task, err := service.Create(ctx, dto.CreateTaskRequest{
		ProjectID:   oldProjectID,
		Title:       "Old title",
		Description: "Old description",
		Status:      model.StatusTodo,
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedTask, err := service.Update(
		ctx,
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

	if len(queue.jobs) != 1 {
		t.Fatalf("expected 1 notification job after update, got %d", len(queue.jobs))
	}
}

func TestDeleteTask(t *testing.T) {
	ctx, service, _ := setupTaskService()

	task, err := service.Create(ctx, dto.CreateTaskRequest{
		ProjectID:   uuid.New(),
		Title:       "Task to delete",
		Description: "Delete this task",
		Status:      model.StatusTodo,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = service.Delete(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.GetByID(ctx, task.ID)
	if err == nil {
		t.Fatalf("expected error when getting deleted task")
	}
}