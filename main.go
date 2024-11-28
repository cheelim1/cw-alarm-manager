package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

const (
	maxRetries     = 3
	initialBackoff = time.Second // 1 second
)

// CloudWatchAPI defines the interface that will be implemented by both the real client and the mock client.
type CloudWatchAPI interface {
	SetAlarmState(ctx context.Context, input *cloudwatch.SetAlarmStateInput, opts ...func(*cloudwatch.Options)) (*cloudwatch.SetAlarmStateOutput, error)
}

func main() {
	// Setup structured logging with log/slog
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	alarmName := os.Getenv("CLOUDWATCH_ALARM_NAME")
	alarmState := os.Getenv("ALARM_STATE")
	alarmReason := os.Getenv("ALARM_REASON")

	if alarmName == "" || alarmState == "" || alarmReason == "" {
		logger.Error("Missing required environment variables",
			slog.String("alarmName", alarmName),
			slog.String("alarmState", alarmState),
			slog.String("alarmReason", alarmReason))
		os.Exit(1)
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Error("Failed to load AWS config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	client := cloudwatch.NewFromConfig(cfg)

	// Set CloudWatch alarm state with retry logic
	err = setAlarmStateWithRetry(client, alarmName, alarmState, alarmReason, logger)
	if err != nil {
		logger.Error("Failed to set CloudWatch alarm state after retries", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Successfully set CloudWatch alarm state", slog.String("AlarmName", alarmName), slog.String("AlarmState", alarmState))
}

func setAlarmStateWithRetry(client CloudWatchAPI, alarmName, alarmState, alarmReason string, logger *slog.Logger) error {
	var attempt int
	var backoff = initialBackoff

	for attempt = 1; attempt <= maxRetries; attempt++ {
		err := setAlarmState(client, alarmName, alarmState, alarmReason)
		if err == nil {
			return nil
		}

		logger.Error("Failed to set CloudWatch alarm state", slog.Int("attempt", attempt), slog.String("error", err.Error()))

		// Exponential backoff before retrying
		time.Sleep(backoff)
		backoff *= 2
	}

	return errors.New("exceeded maximum retries")
}

func setAlarmState(client CloudWatchAPI, alarmName, alarmState, alarmReason string) error {
	input := &cloudwatch.SetAlarmStateInput{
		AlarmName:   aws.String(alarmName),
		StateValue:  types.StateValue(alarmState),
		StateReason: aws.String(alarmReason),
	}

	_, err := client.SetAlarmState(context.TODO(), input)
	return err
}
