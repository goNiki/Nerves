package webhook

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
)

func (d *DataSource) DequeueWebhook(ctx context.Context) (*models.WebhookTask, error) {
	task, err := d.queueRepo.DequeueWebhook(ctx)
	if err != nil {
		d.logger.Error("failed to dequeue webhook", "error", err)
		return nil, err
	}

	if task != nil {
		d.logger.Debug("webhook task dequeued", "task_id", task.ID)
	}

	return task, nil
}
