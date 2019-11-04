// +build integration

package eightpuzzle_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ykamo001/ai/internal/eightpuzzle"
	"os"
)

var provider eightpuzzle.Provider

func init() {
	logger := setupLogger()

	provider = eightpuzzle.NewProvider(logger)
}

func setupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)

	return logger
}

func makeContext() context.Context {
	ctx := context.Background()

	id, err := uuid.NewRandom()
	if err == nil {
		ctx = context.WithValue(ctx, "id", id.ID())
	}

	return ctx
}