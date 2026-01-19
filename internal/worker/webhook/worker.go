package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"

	"github.com/goNiki/Nerves/internal/datasource/webhook"
	"github.com/goNiki/Nerves/internal/models"
)

type Worker struct {
	dataSource *webhook.DataSource
	httpClient *http.Client
	webhookURL string
	maxRetries int
	logger     *slog.Logger
}

func New(
	dataSource *webhook.DataSource,
	webhookURL string,
	timeout time.Duration,
	maxRetries int,
	logger *slog.Logger,
) *Worker {
	// добавляем http:// если не указан протокол
	if webhookURL != "" && webhookURL[:4] != "http" {
		webhookURL = "http://" + webhookURL
	}

	return &Worker{
		dataSource: dataSource,
		httpClient: &http.Client{Timeout: timeout},
		webhookURL: webhookURL,
		maxRetries: maxRetries,
		logger:     logger,
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.logger.Info("webhook worker started", "url", w.webhookURL)

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("webhook worker stopped")
			return
		default:
			task, err := w.dataSource.DequeueWebhook(ctx)
			if err != nil {
				w.logger.Error("failed to dequeue webhook", "error", err)
				time.Sleep(time.Second)
				continue
			}

			if task == nil {
				// очередь пуста - ждём перед следующей проверкой
				time.Sleep(100 * time.Millisecond)
				continue
			}

			w.processTask(ctx, task)
		}
	}
}

func (w *Worker) processTask(ctx context.Context, task *models.WebhookTask) {
	err := w.sendWebhook(ctx, task)
	if err != nil {
		task.Attempt++

		if task.Attempt < w.maxRetries {
			delay := time.Duration(math.Pow(2, float64(task.Attempt))) * time.Second

			w.logger.Warn("webhook delivery failed, scheduling retry",
				"task_id", task.ID,
				"attempt", task.Attempt,
				"delay", delay,
				"error", err,
			)

			time.Sleep(delay)

			if err := w.dataSource.EnqueueWebhook(ctx, task); err != nil {
				w.logger.Error("failed to re-enqueue webhook", "task_id", task.ID, "error", err)
			}
		} else {
			w.logger.Error("webhook delivery failed permanently",
				"task_id", task.ID,
				"attempts", task.Attempt,
				"user_id", task.UserID,
				"error", err,
			)
			// задача уже удалена из очереди при dequeue
		}
	} else {
		w.logger.Info("webhook delivered successfully",
			"task_id", task.ID,
			"user_id", task.UserID,
			"incidents_count", len(task.Incidents),
		)
		// задача уже удалена из очереди при dequeue
	}
}

func (w *Worker) sendWebhook(ctx context.Context, task *models.WebhookTask) error {
	payload := map[string]interface{}{
		"task_id":    task.ID,
		"user_id":    task.UserID,
		"latitude":   task.Latitude,
		"longitude":  task.Longitude,
		"incidents":  task.Incidents,
		"created_at": task.CreatedAt,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", w.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-ID", task.ID.String())

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
