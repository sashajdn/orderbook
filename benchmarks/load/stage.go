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
	StartTime           time.Time
	Duration            time.Duration
	ThroughputPerMinute int
	NumberOfExecutors   int
	Executor            executor.Executor
}

func (s *Stage) Run(ctx context.Context) error {
	ch := make(chan int64, 2<<16)

	go func() {
		// Populate channel
		for i := 0; i < s.ThroughputPerMinute*int(s.Duration); i++ {
			ch <- int64(i)
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
	}

	return nil
}

type Stages []Stage

func (s Stages) Len() int {
	return len(s)
}

func (s Stages) Less(i, j int) bool {
	return s[i].StartTime.Before(s[j].StartTime)
}

func (s Stages) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
