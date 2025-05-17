package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey, slog.LevelKey, slog.SourceKey:
				return slog.Attr{} // omit key
			}
			return a
		},
	})
	Log = slog.New(handler)
	slog.SetDefault(Log)
}
