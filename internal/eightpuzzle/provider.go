package eightpuzzle

import (
	"context"

	"github.com/sirupsen/logrus"
	eightpuzzleservice "github.com/ykamo001/ai/rpc/eightpuzzle"
)

type Provider interface {
	UploadDFindPathata(ctx context.Context, request *eightpuzzleservice.FindPathRequest) (*eightpuzzleservice.FindPathResponse, error)
}

type provider struct {
	logger *logrus.Logger
}

func NewProvider(logger *logrus.Logger) *provider {
	return &provider{
		logger: logger,
	}
}

func (p provider) FindPath(ctx context.Context, request *eightpuzzleservice.FindPathRequest) (*eightpuzzleservice.FindPathResponse, error) {
	entry := p.logger.WithFields(logrus.Fields{
		"id": ctx.Value("id"),
	})

	entry.Info("test")
	return &eightpuzzleservice.FindPathResponse{
		Path: "hello",
	}, nil
}
