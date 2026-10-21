package log

import (
	"log/slog"
	"os"
)

func New(level string) *slog.Logger {
	l := normalizeLevel(level)

	debug := isDebugMode(l)

	options := &slog.HandlerOptions{
		Level:       l,
		AddSource:   debug,
		ReplaceAttr: newAttrProcessor(debug),
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, options))
}
