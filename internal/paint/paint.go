package paint

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

type picture struct {
	values [][]string
	logger *logrus.Logger
}

func NewPicture(values [][]string, logger *logrus.Logger) picture {
	return picture{
		values: values,
		logger: logger,
	}
}

// add back-pressure to prevent massive goroutines
func (p picture) FillIn(ctx context.Context, value string, x, y int64) error {
	entry := p.logger.WithFields(logrus.Fields{
		"id":    ctx.Value("id"),
		"x":     x,
		"y":     y,
		"value": value,
	})

	if x < 0 || x >= int64(len(p.values)) {
		entry.WithField("maxRows", len(p.values)).Error("x is out of range")
		return errors.New("x is out of range")
	}

	if y < 0 || y >= int64(len(p.values[x])) {
		entry.WithField("maxColumnForRow", len(p.values[x])).Error("y is out of range")
		return errors.New("y is out of range")
	}

	ctx = context.WithValue(ctx, "initValue", p.values[x][y])
	ctx = context.WithValue(ctx, "fillInValue", value)

	p.values[x][y] = value
	wg := sync.WaitGroup{}
	wg.Add(4)
	go p.fillIn(ctx, x+1, y, &wg)
	go p.fillIn(ctx, x-1, y, &wg)
	go p.fillIn(ctx, x, y-1, &wg)
	go p.fillIn(ctx, x, y+1, &wg)

	wg.Wait()
	return nil
}

func (p picture) fillIn(ctx context.Context, x, y int64, wg *sync.WaitGroup) error {
	defer wg.Done()

	entry := p.logger.WithFields(logrus.Fields{
		"id":          ctx.Value("id"),
		"initValue":   ctx.Value("initValue"),
		"fillInValue": ctx.Value("fillInValue"),
		"x":           x,
		"y":           y,
	})

	initValue, ok := ctx.Value("initValue").(string)
	if !ok {
		err := errors.New("could not convert context initValue to string")
		entry.Error(err)
		return err
	}

	fillInValue, ok := ctx.Value("fillInValue").(string)
	if !ok {
		err := errors.New("could not convert context fillInValue to string")
		entry.Error(err)
		return err
	}

	if x < 0 || x >= int64(len(p.values)) {
		return nil
	}

	if y < 0 || y >= int64(len(p.values[x])) {
		return nil
	}

	if p.values[x][y] == initValue {
		p.values[x][y] = fillInValue
	} else {
		return nil
	}

	currWg := sync.WaitGroup{}
	currWg.Add(4)
	go p.fillIn(ctx, x+1, y, &currWg)
	go p.fillIn(ctx, x-1, y, &currWg)
	go p.fillIn(ctx, x, y-1, &currWg)
	go p.fillIn(ctx, x, y+1, &currWg)

	currWg.Wait()
	return nil
}
