package load

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/sashajdn/orderbook/benchmarks/executor"
)

type Stage struct {
	Name                string
	RelativeStartTime   time.Duration
	Duration            time.Duration
	ThroughputPerMinute int
	NumberOfExecutors   int
	Executor            executor.Executor
	LoadCurve           LoadCurve
}

func (s *Stage) Run(ctx context.Context) error {
	slog.Info(
		"Running stage: ",
		"stage", s.Name,
		"throughput_per_minute", fmt.Sprintf("%d", s.ThroughputPerMinute),
		"duration", s.Duration.String(),
		"number_of_executors", fmt.Sprintf("%d", s.NumberOfExecutors),
	)

	var lcg LoadCurveGenerator
	switch s.LoadCurve {
	case LoadCurveLinear:
		lcg = NewLinearLoadCurveGenerator(s.ThroughputPerMinute, int(s.Duration/time.Minute), time.Minute)
	default:
		return fmt.Errorf("Load curve not supported: %s", s.LoadCurve)
	}

	ch := make(chan int64, 2<<16)
	go func() {
		// Populate channel
		defer close(ch)
		for i := 0; i < s.ThroughputPerMinute*int(s.Duration/time.Minute); i++ {
			select {
			case ch <- int64(i):
			case <-ctx.Done():
				return
			}
		}
	}()

	deadline := time.Now().UTC().Add(s.Duration)
	workerCtx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	var wg sync.WaitGroup
	for executionWorker := 0; executionWorker < s.NumberOfExecutors; executionWorker++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range ch {
				select {
				case <-ctx.Done():
					return
				default:
				}

				lcg.Wait()
				if err := s.Executor.RunIteration(workerCtx); err != nil {
					slog.Error("execute work", "error", err, "idx", fmt.Sprintf("%d", executionWorker))
				}
			}
		}()
	}

	done := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		slog.Info("stage done", "stage", s.Name)
	case <-time.After(s.Duration): // TODO: calc timer above
		slog.Warn(`stage not finished execution in alloted timeframe`, "stage", s.Name)
	case <-ctx.Done():
		slog.Warn("Context cancelled exiting...")
	}

	return nil
}

type Stages []Stage

func (s Stages) Len() int {
	return len(s)
}

func (s Stages) Less(i, j int) bool {
	return s[i].RelativeStartTime < (s[j].RelativeStartTime)
}

func (s Stages) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
