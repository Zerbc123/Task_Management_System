package services

import (
	"context"
	"task-management/internal/comment/dto"
	"task-management/internal/comment/model"
	"task-management/internal/comment/repository"
	"testing"

	"github.com/google/uuid"
)

type mockCommentRepository struct {
	comments map[uuid.UUID]*model.Comment
}

func newMockCommentRepository() *mockCommentRepository {
	return &mockCommentRepository{comments: make(map[uuid.UUID]*model.Comment)}
}

func (m *mockCommentRepository) Create(cxt context.Context, comment *model.Comment) error {
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	comment, ok := m.comments[id]
	if !ok {
		return nil, repository.ErrCommentNotFound
	}

	return comment, nil
}

func (m *mockCommentRepository) GetByTaskID(ctx context.Context, TaskID uuid.UUID) ([]*model.Comment, error) {
	comments := make([]*model.Comment, 0)

	for _, comment := range m.comments {
		if comment.TaskID == TaskID {
			comments = append(comments, comment)
		}
	}

	return comments, nil
}

func (m *mockCommentRepository) Update(ctx context.Context, comment *model.Comment) error {
	_, ok := m.comments[comment.ID]
	if !ok {
		return repository.ErrCommentNotFound
	}

	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, ok := m.comments[id]
	if !ok {
		return repository.ErrCommentNotFound
	}

	delete(m.comments, id)
	return nil
}

func SetupCommentService() (context.Context, CommentService) {
	ctx := context.Background()

	repo := newMockCommentRepository()
	service := NewCommentService(repo, nil)

	return ctx, service
}

func TestCreateComment(t *testing.T) {
	ctx, service := SetupCommentService()

	taskID := uuid.New()
	userID := uuid.New()

	comment, err := service.Create(ctx, taskID, userID, dto.CreateCommentRequest{
		Content: "Redis queue completed",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if comment.ID == uuid.Nil {
		t.Fatal("expected comment ID to be generated")
	}
	
	if comment.TaskID != taskID {
		t.Fatalf("expected taskID %s, got %s", taskID, comment.TaskID)
	}

	if comment.UserID != userID {
		t.Fatalf("expected userID %s, got %s", userID, comment.UserID)
	}

	if comment.Content != "Redis queue completed" {
		t.Fatalf("expected Redis queue completed, got %s", comment.Content)
	}
}

func TestGetCommentByTaskID(t *testing.T) {
	ctx, service := SetupCommentService()

	taskID := uuid.New()
	userID := uuid.New()

	_, err := service.Create(ctx, taskID, userID, dto.CreateCommentRequest{
		Content: "comment 1",
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = service.Create(ctx, taskID, userID, dto.CreateCommentRequest{
		Content: "Comment 2",
	})

	if err != nil {
		t.Fatal(err)
	}

	comments, err := service.GetByTaskID(ctx, taskID)
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}
}

func TestUpdateComment(t *testing.T) {
	ctx, service := SetupCommentService()

	taskID := uuid.New()
	userID := uuid.New()

	comment, err := service.Create(ctx, taskID, userID, dto.CreateCommentRequest{
		Content: "Old comment",
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedComment, err := service.Update(ctx, comment.ID, dto.UpdateCommentRequest{
		Content: "Updated comment",
	})
	if err != nil {
		t.Fatal(err)
	}

	if updatedComment.Content != "Updated comment" {
		t.Fatalf("expected Updated comment, got %s", updatedComment.Content)
	}
}

func TestDeleteComment(t *testing.T) {
	ctx, service := SetupCommentService()

	taskID := uuid.New()
	userID := uuid.New()

	comment, err := service.Create(ctx, taskID, userID, dto.CreateCommentRequest{
		Content: "Comment to delete",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = service.Delete(ctx, comment.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.Update(ctx, comment.ID, dto.UpdateCommentRequest{
		Content: "Should fail",
	})
	if err == nil {
		t.Fatal("expected error after deleting comment")
	}
}
