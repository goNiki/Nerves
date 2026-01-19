package webhook

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
)

func (s *Service) Enqueue(ctx context.Context, task *models.WebhookTask) error {
	if err := s.dataSource.EnqueueWebhook(ctx, task); err != nil {
		s.logger.Error("failed to enqueue webhook", "task_id", task.ID, "error", err)
		return err
	}

	s.logger.Info("webhook task enqueued", "task_id", task.ID, "user_id", task.UserID)
	return nil
}
