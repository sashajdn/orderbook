package slog

import (
	"log/slog"
	"os"
)

type Level slog.Level

const (
	Debug Level = Level(slog.LevelDebug)
	Info  Level = Level(slog.LevelInfo)
)

func Init(level Level) {
	var l = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(level),
	}))

	slog.SetDefault(l)
}
