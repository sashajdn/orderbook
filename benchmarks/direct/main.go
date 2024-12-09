package main

import (
	"context"
	"fmt"
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
	slog.Info("Running direct benchmarking...")
	// LOB setup.
	lob := lob.NewOrderbook(2 << 16)

	// Client setup.
	client := client.NewLOBClient(lob)

	// Executor setup.
	taker := executor.NewTaker(executor.TakerConfig{
		Users:  10,
		Client: client,
	})

	maker := executor.NewMaker(executor.MakerConfig{
		Client:      client,
		Users:       10,
		Spread:      5,
		Midprice:    1000,
		LaplaceBeta: 1.0,
	})

	slog.Info("Direct benchmark setup complete")
	slog.Info(`Direct benchmark executing stages...`)

	stages := []*load.Stage{
		{
			Name:                "book_warmup",
			RelativeStartTime:   0,
			Duration:            1 * time.Minute,
			ThroughputPerMinute: 1000,
			NumberOfExecutors:   10,
			Executor:            maker,
			LoadCurve:           load.LoadCurveLinear,
		},
		{
			Name:                "maker",
			RelativeStartTime:   1 * time.Minute,
			Duration:            1 * time.Minute,
			ThroughputPerMinute: 100,
			NumberOfExecutors:   10,
			Executor:            maker,
			LoadCurve:           load.LoadCurveLinear,
		},
		{
			Name:                "taker",
			RelativeStartTime:   1 * time.Minute,
			Duration:            1 * time.Minute,
			ThroughputPerMinute: 100,
			NumberOfExecutors:   10,
			Executor:            taker,
			LoadCurve:           load.LoadCurveLinear,
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	generator := load.NewGenerator(stages)
	if err := generator.Run(ctx); err != nil {
		slog.Error("Failed to run generator", "error", err)
	}
}

var _ executor.Executor = &LiquidityPrinter{}

type LiquidityPrinter struct {
	lob *lob.Orderbook
}

func (l *LiquidityPrinter) RunIteration(_ context.Context) error {
	bv, av := l.lob.Volume()
	slog.Info("VOLUME", "bid", fmt.Sprintf("%.4f", bv), "ask", fmt.Sprintf("%.4f", av))
	return nil
}

func (l *LiquidityPrinter) Name() string { return "liquidity printer" }
