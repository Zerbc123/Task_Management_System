package services

import (
	"context"
	"task-management/internal/comment/dto"
	"task-management/internal/comment/model"
	"task-management/internal/comment/repository"
	"task-management/internal/websocket"
	"time"

	"github.com/google/uuid"
)

type CommentService interface {
	Create(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, req dto.CreateCommentRequest) (*model.Comment, error)
	GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]*model.Comment, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateCommentRequest) (*model.Comment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type commentService struct {
	repo repository.CommentRepository
	hub  *websocket.Hub
}

func NewCommentService(
	repo repository.CommentRepository,
	hub *websocket.Hub,
) CommentService {
	return &commentService{
		repo: repo,
		hub:  hub,
	}
}

func (s *commentService) Create(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, req dto.CreateCommentRequest) (*model.Comment, error) {
	now := time.Now()

	comment := &model.Comment{
		ID:        uuid.New(),
		TaskID:    taskID,
		UserID:    userID,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, err
	}

	if s.hub != nil {
		s.hub.BroadcastJSON(
			websocket.CommentCreatedEvent{
				Type:      "comment.created",
				CommentID: comment.ID,
				TaskID:    comment.TaskID,
				UserID:    comment.UserID,
				Content:   comment.Content,
			},
		)
	}

	return comment, nil
}

func (s *commentService) GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]*model.Comment, error) {
	return s.repo.GetByTaskID(ctx, taskID)
}

func (s *commentService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCommentRequest) (*model.Comment, error) {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	comment.Content = req.Content
	comment.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
