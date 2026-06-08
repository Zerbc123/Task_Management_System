package notification

import (
	"context"
	"errors"
	"log"
	"time"
)

type Worker struct {
	queue *Queue
}

func NewWorker(queue *Queue) *Worker {
	return &Worker{
		queue: queue,
	}
}

func (w *Worker) Start(ctx context.Context) {
	log.Println("notification worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("notification worker stopped")
			return

		default:
			job, err := w.queue.Pop(ctx)
			if err != nil {
				log.Println("failed to pop notification job:", err)
				continue
			}

			w.handleJob(ctx, *job)
		}
	}
}

func (w *Worker) handleJob(ctx context.Context, job NotificationJob) {
	err := w.sendNotification(job)
	if err == nil {
		log.Printf(
			"notification sent successfully: task=%s user=%v event=%s",
			job.TaskID,
			job.UserID,
			job.EventType,
		)
		return
	}

	log.Printf(
		"notification failed: task=%s retry=%d/%d error=%v",
		job.TaskID,
		job.RetryCount,
		job.MaxRetry,
		err,
	)

	if job.RetryCount < job.MaxRetry {
		job.RetryCount++

		time.Sleep(2 * time.Second)

		log.Printf(
			"retry notification: task=%s retry=%d/%d",
			job.TaskID,
			job.RetryCount,
			job.MaxRetry,
		)

		if err := w.queue.Push(ctx, job); err != nil {
			log.Println("failed to retry notification job:", err)
		}

		return
	}

	log.Printf(
		"notification permanently failed: task=%s user=%v event=%s",
		job.TaskID,
		job.UserID,
		job.EventType,
	)
}

func (w *Worker) sendNotification(job NotificationJob) error {
	if job.UserID == nil {
		return errors.New("assignee is empty")
	}

	log.Printf(
		"send notification: user=%s task=%s event=%s message=%s",
		job.UserID.String(),
		job.TaskID,
		job.EventType,
		job.Message,
	)

	return nil
}
