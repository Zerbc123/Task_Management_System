package notification

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

const notificationQueueKey = "queue:notifications"

type JobQueue interface {
	Push(ctx context.Context, job NotificationJob) error
}

type Queue struct {
	redis *redis.Client
}

func NewQueue(redisClient *redis.Client) *Queue {
	return &Queue{
		redis: redisClient,
	}
}

func (q *Queue) Push(ctx context.Context, job NotificationJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.redis.LPush(ctx, notificationQueueKey, data).Err()
}

func (q *Queue) Pop(ctx context.Context) (*NotificationJob, error) {
	result, err := q.redis.BRPop(ctx, 0, notificationQueueKey).Result()
	if err != nil {
		return nil, err
	}

	var job NotificationJob
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, err
	}

	return &job, nil
}
