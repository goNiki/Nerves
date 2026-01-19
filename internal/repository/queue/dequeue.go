package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	"github.com/redis/go-redis/v9"
)

func (r *Repository) DequeueWebhook(ctx context.Context) (*models.WebhookTask, error) {
	const op = "repository.queue.dequeuewebhook"

	result, err := r.redis.BRPop(ctx, dequeueTimeout, webhookQueueKey).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrDequeueWebhook, err)
	}

	if len(result) < 2 {
		return nil, nil
	}

	var task models.WebhookTask

	if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrDequeueWebhook, err)
	}

	return &task, nil

}
