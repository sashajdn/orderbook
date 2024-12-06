package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/sashajdn/orderbook/benchmarks/client"
	"github.com/sashajdn/orderbook/benchmarks/executor"
	"github.com/sashajdn/orderbook/benchmarks/load"
	"github.com/sashajdn/orderbook/lob"
)

func main() {
	// LOB setup.
	lob := lob.NewOrderbook(2 << 16)

	// Client setup.
	client := client.NewLOBClient(lob)

	// Executor setup.
	marketMaker := executor.NewMarketMaker(10, client)

	now := time.Now()
	stages := []load.Stage{
		{
			Name:                "stage_1",
			StartTime:           now.Add(1 * time.Minute),
			Duration:            1 * time.Minute,
			ThroughputPerMinute: 10_000,
			NumberOfExecutors:   10,
			Executor:            marketMaker,
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	generator := load.NewGenerator(stages)
	if err := generator.Run(ctx); err != nil {
		slog.Error("Failed to run generator", "error", err)
	}
}
