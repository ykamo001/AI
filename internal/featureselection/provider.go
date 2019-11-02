package featureselection

import (
	"context"

	"github.com/sirupsen/logrus"
	featureselectionservice "github.com/ykamo001/ai/rpc/featureselection"
)

type Provider interface {
	UploadData(ctx context.Context, request *featureselectionservice.UploadDataRequest) (*featureselectionservice.Empty, error)
}

type provider struct {
	logger *logrus.Logger
}

func NewProvider(logger *logrus.Logger) *provider {
	return &provider{
		logger: logger,
	}
}

func (p provider) UploadData(ctx context.Context, request *featureselectionservice.UploadDataRequest) (*featureselectionservice.Empty, error) {
	p.logger.Info("hello world")
	return &featureselectionservice.Empty{}, nil
}
