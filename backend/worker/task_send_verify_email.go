package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	UserID string `json:"user_id"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("Task enqueued successfully")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadSendVerifyEmail

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Error().Err(err).Bytes("payload", task.Payload()).Msg("Failed to unmarshal task payload")
		return asynq.SkipRetry
	}
	userID, err := strconv.Atoi(payload.UserID)
	if err != nil {
		log.Error().Err(err).Str("user_id", payload.UserID).Msg("Invalid user ID format")
		return asynq.SkipRetry
	}
	_, err = processor.store.GetUserByID(ctx, int32(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			return asynq.SkipRetry
		}
		return fmt.Errorf("failed to get user by ID: %w", err)
	}

	// TODO: Implement the logic to send the verification email.
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", "debdipmukherjee52@gmail.com").Msg("Processing task to send verification email")
	return nil
}
