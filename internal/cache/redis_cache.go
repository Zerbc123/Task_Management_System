package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"task-management/internal/task/model"
)

const (
	taskListKey = "tasks:all"
	taskTTL     = 5 * time.Minute
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{
		client: client,
	}
}

func taskDetailKey(id uuid.UUID) string {
	return "tasks:" + id.String()
}

func (c *redisCache) GetTaskList(ctx context.Context) ([]*model.Task, bool) {
	data, err := c.client.Get(ctx, taskListKey).Result()
	if err != nil {
		return nil, false
	}

	var tasks []*model.Task
	if err := json.Unmarshal([]byte(data), &tasks); err != nil {
		log.Println("failed to unmarshal task list cache:", err)
		return nil, false
	}

	return tasks, true
}

func (c *redisCache) SetTaskList(ctx context.Context, tasks []*model.Task) {
	data, err := json.Marshal(tasks)
	if err != nil {
		log.Println("failed to marshal task list cache:", err)
		return
	}

	if err := c.client.Set(ctx, taskListKey, data, taskTTL).Err(); err != nil {
		log.Println("failed to set task list cache:", err)
	}
}

func (c *redisCache) DeleteTaskList(ctx context.Context) {
	if err := c.client.Del(ctx, taskListKey).Err(); err != nil {
		log.Println("failed to delete task list cache:", err)
	}
}

func (c *redisCache) GetTaskByID(ctx context.Context, id uuid.UUID) (*model.Task, bool) {
	data, err := c.client.Get(ctx, taskDetailKey(id)).Result()
	if err != nil {
		return nil, false
	}

	var task model.Task
	if err := json.Unmarshal([]byte(data), &task); err != nil {
		log.Println("failed to unmarshal task detail cache:", err)
		return nil, false
	}

	return &task, true
}

func (c *redisCache) SetTaskByID(ctx context.Context, task *model.Task) {
	data, err := json.Marshal(task)
	if err != nil {
		log.Println("failed to marshal task detail cache:", err)
		return
	}

	if err := c.client.Set(ctx, taskDetailKey(task.ID), data, taskTTL).Err(); err != nil {
		log.Println("failed to set task detail cache:", err)
	}
}

func (c *redisCache) DeleteTaskByID(ctx context.Context, id uuid.UUID) {
	if err := c.client.Del(ctx, taskDetailKey(id)).Err(); err != nil {
		log.Println("failed to delete task detail cache:", err)
	}
}