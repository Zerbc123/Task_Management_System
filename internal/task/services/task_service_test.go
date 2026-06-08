package services

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"task-management/internal/audit/model"
	"task-management/internal/notification"
	"task-management/internal/task/dto"
	taskModel "task-management/internal/task/model"
	memoryRepo "task-management/internal/task/repository/memory"
)

type mockCache struct {
	taskList map[string][]*taskModel.Task
	tasks    map[uuid.UUID]*taskModel.Task
}

func newMockCache() *mockCache {
	return &mockCache{
		taskList: make(map[string][]*taskModel.Task),
		tasks:    make(map[uuid.UUID]*taskModel.Task),
	}
}

func (m *mockCache) GetTaskList(ctx context.Context) ([]*taskModel.Task, bool) {
	tasks, ok := m.taskList["tasks:all"]
	return tasks, ok
}

func (m *mockCache) SetTaskList(ctx context.Context, tasks []*taskModel.Task) {
	m.taskList["tasks:all"] = tasks
}

func (m *mockCache) DeleteTaskList(ctx context.Context) {
	delete(m.taskList, "tasks:all")
}

func (m *mockCache) GetTaskByID(ctx context.Context, id uuid.UUID) (*taskModel.Task, bool) {
	task, ok := m.tasks[id]
	return task, ok
}

func (m *mockCache) SetTaskByID(ctx context.Context, task *taskModel.Task) {
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

type mockAuditService struct {
	logs []*model.AuditLog
}

func newMockAuditService() *mockAuditService {
	return &mockAuditService{
		logs: make([]*model.AuditLog, 0),
	}
}

func (m *mockAuditService) Log(
	ctx context.Context,
	userID uuid.UUID,
	entityID uuid.UUID,
	entityType string,
	eventType string,
	description string,
) error {
	m.logs = append(m.logs, &model.AuditLog{
		ID:          uuid.New(),
		UserID:      userID,
		EntityID:    entityID,
		EntityType:  entityType,
		EventType:   eventType,
		Description: description,
	})

	return nil
}

func (m *mockAuditService) GetByEntityID(
	ctx context.Context,
	entityID uuid.UUID,
) ([]*model.AuditLog, error) {
	result := make([]*model.AuditLog, 0)

	for _, log := range m.logs {
		if log.EntityID == entityID {
			result = append(result, log)
		}
	}

	return result, nil
}

func setupTaskService() (context.Context, TaskService, *mockQueue, *mockAuditService) {
	ctx := context.Background()

	repo := memoryRepo.NewMemoryTaskRepository()
	cache := newMockCache()
	queue := newMockQueue()
	audit := newMockAuditService()

	service := NewTaskService(repo, cache, queue, nil, audit)

	return ctx, service, queue, audit
}

func TestCreateTask(t *testing.T) {
	ctx, service, queue, audit := setupTaskService()

	userID := uuid.New()
	projectID := uuid.New()
	assigneeID := uuid.New()

	req := dto.CreateTaskRequest{
		ProjectID:   projectID,
		Title:       "Learn Go",
		Description: "Study unit test",
		Status:      taskModel.StatusTodo,
		AssigneeID:  &assigneeID,
	}

	task, err := service.Create(ctx, userID, req)
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

	if task.Status != taskModel.StatusTodo {
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

	if len(audit.logs) != 1 {
		t.Fatalf("expected 1 audit log, got %d", len(audit.logs))
	}

	if audit.logs[0].EventType != "task.created" {
		t.Fatalf("expected audit event task.created, got %s", audit.logs[0].EventType)
	}
}

func TestGetAll(t *testing.T) {
	ctx, service, _, _ := setupTaskService()

	userID := uuid.New()

	_, err := service.Create(ctx, userID, dto.CreateTaskRequest{
		ProjectID:   uuid.New(),
		Title:       "Task 1",
		Description: "Description 1",
		Status:      taskModel.StatusTodo,
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
	ctx, service, _, _ := setupTaskService()

	userID := uuid.New()

	createdTask, err := service.Create(ctx, userID, dto.CreateTaskRequest{
		ProjectID:   uuid.New(),
		Title:       "Task detail",
		Description: "Get task by ID",
		Status:      taskModel.StatusTodo,
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
	ctx, service, queue, audit := setupTaskService()

	userID := uuid.New()
	oldProjectID := uuid.New()
	newProjectID := uuid.New()
	assigneeID := uuid.New()

	task, err := service.Create(ctx, userID, dto.CreateTaskRequest{
		ProjectID:   oldProjectID,
		Title:       "Old title",
		Description: "Old description",
		Status:      taskModel.StatusTodo,
	})
	if err != nil {
		t.Fatal(err)
	}

	queue.jobs = nil
	audit.logs = nil

	updatedTask, err := service.Update(
		ctx,
		task.ID,
		userID,
		dto.UpdateTaskRequest{
			ProjectID:   newProjectID,
			Title:       "New title",
			Description: "New description",
			Status:      taskModel.StatusDone,
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

	if updatedTask.Status != taskModel.StatusDone {
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

	if len(audit.logs) != 2 {
		t.Fatalf("expected 2 audit logs after update, got %d", len(audit.logs))
	}
}

func TestDeleteTask(t *testing.T) {
	ctx, service, _, audit := setupTaskService()

	userID := uuid.New()

	task, err := service.Create(ctx, userID, dto.CreateTaskRequest{
		ProjectID:   uuid.New(),
		Title:       "Task to delete",
		Description: "Delete this task",
		Status:      taskModel.StatusTodo,
	})
	if err != nil {
		t.Fatal(err)
	}

	audit.logs = nil

	err = service.Delete(ctx, task.ID, userID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.GetByID(ctx, task.ID)
	if err == nil {
		t.Fatalf("expected error when getting deleted task")
	}

	if len(audit.logs) != 1 {
		t.Fatalf("expected 1 audit log after delete, got %d", len(audit.logs))
	}

	if audit.logs[0].EventType != "task.deleted" {
		t.Fatalf("expected task.deleted audit log, got %s", audit.logs[0].EventType)
	}
}