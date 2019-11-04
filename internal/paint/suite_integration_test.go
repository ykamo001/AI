// +build integration

package paint_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ykamo001/backend/internal/paint"
	"os"
)

var provider paint.Provider

func init() {
	logger := setupLogger()

	provider = paint.NewProvider(logger)
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
