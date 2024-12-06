package load

import (
	"context"
	"log/slog"
	"sort"
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
	for _, stage := range g.stages {
		if err := stage.Run(ctx); err != nil {
			slog.Error("Failed to run stage", "stage", stage.Name, "error", err)
			return err
		}
	}

	return nil
}
