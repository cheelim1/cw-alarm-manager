package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/stretchr/testify/assert"
)

// MockCloudWatchClient implements the CloudWatchAPI interface to mock the SetAlarmState method.
type MockCloudWatchClient struct {
	err error
}

func (m *MockCloudWatchClient) SetAlarmState(ctx context.Context, input *cloudwatch.SetAlarmStateInput, opts ...func(*cloudwatch.Options)) (*cloudwatch.SetAlarmStateOutput, error) {
	return nil, m.err
}

func TestSetAlarmStateWithRetry_Success(t *testing.T) {
	mockClient := &MockCloudWatchClient{err: nil}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := setAlarmStateWithRetry(mockClient, "TestAlarm", "OK", "Test reason", logger)
	assert.NoError(t, err)
}

func TestSetAlarmStateWithRetry_Failure(t *testing.T) {
	mockClient := &MockCloudWatchClient{err: errors.New("API error")}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := setAlarmStateWithRetry(mockClient, "TestAlarm", "OK", "Test reason", logger)
	assert.Error(t, err)
	assert.Equal(t, "exceeded maximum retries", err.Error())
}
