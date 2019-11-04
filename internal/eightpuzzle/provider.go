package eightpuzzle

import (
	"context"
	"errors"
	"github.com/twitchtv/twirp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	eightpuzzleservice "github.com/ykamo001/ai/rpc/eightpuzzle"
)

type Provider interface {
	FindPath(ctx context.Context, request *eightpuzzleservice.FindPathRequest) (*eightpuzzleservice.FindPathResponse, error)
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

	initial, err := p.createArrayFromString(ctx, request.Initial)
	if err != nil {
		entry.WithField("parameter", "initial").Error("createArrayFromString")
		return &eightpuzzleservice.FindPathResponse{}, twirp.NewError(twirp.InvalidArgument, "invalid argument")
	}

	goal, err := p.createArrayFromString(ctx, request.Goal)
	if err != nil {
		entry.WithField("parameter", "goal").Error("createArrayFromString")
		return &eightpuzzleservice.FindPathResponse{}, twirp.NewError(twirp.InvalidArgument, "invalid argument")
	}

	puzzle := NewPuzzle(initial, goal, p.logger)

	if valid, err := puzzle.CheckValidity(ctx); !valid || err != nil {
		entry.Error("CheckValidity")
		return &eightpuzzleservice.FindPathResponse{}, twirp.NewError(twirp.InvalidArgument, "invalid argument")
	}

	return &eightpuzzleservice.FindPathResponse{
		Path: "hello",
	}, nil
}


func (p provider) createArrayFromString(ctx context.Context, strArr string) ([][]int, error) {
	entry := p.logger.WithFields(logrus.Fields{
		"id": ctx.Value("id"),
		"strArray": strArr,
	})

	var matrix [][]int
	strInit := strings.Split(strings.TrimSpace(strArr), ",")
	n := 0
	for i, strRow := range strInit {
		entry = entry.WithField("row", i)
		var row []int
		nums := strings.Split(strings.TrimSpace(strRow), " ")
		if n == 0 {
			n = len(nums)
			entry = entry.WithField("n", n)
		}
		if n != len(nums) {
			entry.WithField("currentRowLength", len(nums)).Error("row length does not match rest of matrix")
			return [][]int{}, errors.New("array dimension mismatch")
		}
		for j, strNum := range nums {
			entry = entry.WithField("column", j)
			num, err := strconv.Atoi(strNum)
			if err != nil {
				entry.WithField("value", strNum).WithError(err).Error("invalid value in array")
				return [][]int{}, err
			}
			row = append(row, num)
		}
		matrix = append(matrix, row)
	}

	return matrix, nil
}