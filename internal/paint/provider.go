package paint

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/twitchtv/twirp"
	paintservice "github.com/ykamo001/backend/rpc/paint"
)

type Provider interface {
	FillIn(ctx context.Context, request *paintservice.FillInRequest) (*paintservice.FillInResponse, error)
}

type provider struct {
	logger *logrus.Logger
}

func NewProvider(logger *logrus.Logger) *provider {
	return &provider{
		logger: logger,
	}
}

func (p provider) FillIn(ctx context.Context, request *paintservice.FillInRequest) (*paintservice.FillInResponse, error) {
	entry := p.logger.WithFields(logrus.Fields{
		"id": ctx.Value("id"),
	})

	var matrix [][]string
	for _, row := range request.Matrix {
		var currRow []string
		for _, value := range row.Array {
			currRow = append(currRow, value)
		}
		matrix = append(matrix, currRow)
	}

	picture := NewPicture(matrix, p.logger)
	err := picture.FillIn(ctx, request.Value, request.X, request.Y)
	if err != nil {
		entry.Error("FillIn")
		return &paintservice.FillInResponse{}, twirp.NewError(twirp.Internal, "internal error")
	}

	var grpcMatrix []*paintservice.InternalArray
	for _, row := range picture.values {
		var currRow paintservice.InternalArray
		for _, value := range row {
			currRow.Array = append(currRow.Array, value)
		}
		grpcMatrix = append(grpcMatrix, &currRow)
	}
	return &paintservice.FillInResponse{
		Matrix: grpcMatrix,
	}, nil
}
