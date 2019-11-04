package eightpuzzle

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

type puzzle struct {
	initial [][]int
	goal [][]int
	logger *logrus.Logger
}

func NewPuzzle(initial, goal [][]int, logger *logrus.Logger) puzzle {
	return puzzle{
		initial: initial,
		goal:    goal,
		logger:  logger,
	}
}

func (p puzzle) CheckValidity(ctx context.Context) (bool, error) {
	entry := p.logger.WithFields(logrus.Fields{
		"id": ctx.Value("id"),
	})

	initialValues := make(map[int]bool)
	goalValues := make(map[int]bool)

	for i, row := range p.initial {
		rowEntry := entry.WithFields(logrus.Fields{
			"row": i,
		})
		for j, column := range row {
			currEntry := rowEntry.WithFields(logrus.Fields{
				"column": j,
			})

			if _, exists := initialValues[column]; exists {
				currEntry.WithField("value", column).Error("repeated value in initial matrix")
				return false, errors.New("repeated value in initial matrix")
			} else {
				initialValues[column] = true
			}
		}
	}

	for i, row := range p.goal {
		rowEntry := entry.WithFields(logrus.Fields{
			"row": i,
		})
		for j, column := range row {
			currEntry := rowEntry.WithFields(logrus.Fields{
				"column": j,
			})

			if _, exists := goalValues[column]; exists {
				currEntry.WithField("value", column).Error("repeated value")
				return false, errors.New("repeated value in goal matrix")
			} else {
				goalValues[column] = true
			}

			if _, exists := initialValues[column]; !exists {
				currEntry.WithField("value", column).Error("value in goal matrix does not exist in initial matrix")
				return false, errors.New("value in goal matrix does not exist in initial matrix")
			}
		}
	}

	if len(initialValues) != len(goalValues) {
		entry.WithFields(logrus.Fields{
			"initialValues": len(initialValues),
			"goalValues": len(goalValues),
		}).Error("values in initial matrix do not match values in goal matrix")
		return false, errors.New("values in initial matrix do not match values in goal matrix")
	}

	return true, nil
}
