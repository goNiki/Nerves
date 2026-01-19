package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

func (r *Repository) EnqueueWebhook(ctx context.Context, task *models.WebhookTask) error {
	const op = "repository.queue.enqueuewebhook"
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorapp.ErrEnqueueWebhook, err)
	}

	err = r.redis.LPush(ctx, webhookQueueKey, data).Err()
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorapp.ErrEnqueueWebhook, err)
	}

	return nil
}
