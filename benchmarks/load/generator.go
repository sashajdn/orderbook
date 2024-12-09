package load

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"time"
)

func NewGenerator(stages Stages) *Generator {
	sort.Sort(stages)

	return &Generator{
		stages: stages,
	}
}

type Generator struct {
	stages Stages
}

func (g *Generator) Run(ctx context.Context) error {
	s := &scheduler{}
	s.run(ctx, g.stages)

	return nil
}

type scheduler struct{}

func (s *scheduler) run(ctx context.Context, stages []*Stage) error {
	var wg sync.WaitGroup
	wg.Add(len(stages))

	then := time.Now()

	for _, stage := range stages {
		stage := stage
		go func() {
			defer wg.Done()

			time.Sleep(stage.RelativeStartTime - time.Since(then))

			slog.Info("====== Starting stage", "stage", stage.Name)

			if err := stage.Run(ctx); err != nil {
				slog.Error("Failed to run stage", "stage", stage.Name, "error", err)
				return
			}

			slog.Info("====== Stage finished", "stage", stage.Name)
		}()
	}

	wg.Wait()

	return nil
}
