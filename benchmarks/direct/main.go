package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/sashajdn/orderbook/benchmarks/load"
	"github.com/sashajdn/orderbook/lob"
)

func main() {
	lob := lob.NewOrderbook(2 << 16)
	client := load.NewLOBClient(lob)

	now := time.Now()
	stages := []load.Stage{
		{
			Name:                "stage_1",
			StartTime:           now.Add(1 * time.Minute),
			Duration:            1 * time.Minute,
			ThroughputPerMinute: 10_000,
			NumberOfExecutors:   10,
			Executor: func(ctx context.Context, client load.Client) error {
				_, err := client.AddOrder(ctx, load.AddOrderRequest{})
				if err != nil {
					return fmt.Errorf("add order: %w", err)
				}

				return nil
			},
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	generator := load.NewGenerator(client, stages)
	if err := generator.Run(ctx); err != nil {
		slog.Error("Failed to run generator", "error", err)
	}
}
