package webhook

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
)

func (d *DataSource) EnqueueWebhook(ctx context.Context, task *models.WebhookTask) error {
	if err := d.queueRepo.EnqueueWebhook(ctx, task); err != nil {
		d.logger.Error("failed to enqueue webhook", "task_id", task.ID, "error", err)
		return err
	}

	d.logger.Debug("webhook task enqueued", "task_id", task.ID, "user_id", task.UserID)
	return nil
}
