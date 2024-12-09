package load

import (
	"context"
	"fmt"
	"time"
)

type LoadCurve int8

const (
	LoadCurveLinear LoadCurve = iota + 1
)

var _ fmt.Stringer = new(LoadCurve)

func (l LoadCurve) String() string {
	switch l {
	case LoadCurveLinear:
		return "linear_load_curve"
	default:
		return "unknown"
	}
}

type LoadCurveGenerator interface {
	Wait()
	Close()
}

func NewLinearLoadCurveGenerator(rate int, length int, unit time.Duration) *LinearLoadCurveGenerator {
	ch := make(chan struct{}, rate)

	tokens := rate * length
	l := &LinearLoadCurveGenerator{
		tokens: ch,
	}

	go func() {
		t := time.NewTimer(time.Duration(rate) / unit)
		defer t.Stop()

		for {
			if tokens == 0 {
				return
			}

			select {
			case <-t.C:
			case <-l.ctx.Done():
				return
			}

			select {
			case ch <- struct{}{}:
				tokens--
			case <-l.ctx.Done():
				return
			}
		}

	}()

	l.ctx, l.cancel = context.WithCancel(context.Background())

	return l
}

var _ LoadCurveGenerator = &LinearLoadCurveGenerator{}

type LinearLoadCurveGenerator struct {
	tokens chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
}

func (l *LinearLoadCurveGenerator) Wait() {
	<-l.tokens
}

func (l *LinearLoadCurveGenerator) Close() {
	l.cancel()
}
