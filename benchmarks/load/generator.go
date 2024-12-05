package load

import (
	"context"
	"log/slog"
	"sort"
)

func NewGenerator(client Client, stages Stages) *Generator {
	return &Generator{
		client: client,
		stages: stages,
	}
}

type Generator struct {
	client Client
	stages Stages
}

func (g *Generator) Run(ctx context.Context) error {
	sort.Sort(g.stages)

	for _, stage := range g.stages {
		if err := stage.Run(ctx, g.client); err != nil {
			slog.Error("Failed to run stage", "stage", stage.Name, "error", err)
			return err
		}
	}

	return nil
}
