package log

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/whaoa/public-api-bridge/pkg/errors"
)

func normalizeLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func isDebugMode(level slog.Level) bool {
	return level <= slog.LevelDebug
}

func newAttrProcessor(debug bool) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, attr slog.Attr) slog.Attr {
		switch attr.Key {
		case slog.TimeKey:
			attr.Value = slog.StringValue(attr.Value.Time().UTC().Format(time.RFC3339))
		case slog.LevelKey:
			attr.Value = slog.StringValue(strings.ToLower(attr.Value.String()))
		case slog.SourceKey:
			if !debug {
				break
			}
			if src, ok := attr.Value.Any().(*slog.Source); ok {
				attr.Key = "caller"
				source := fmt.Sprintf("%v:%v", filepath.Base(src.File), src.Line)
				attr.Value = slog.StringValue(source)
			} else {
				attr = slog.Attr{}
			}
		case attrKeyError:
			if err, ok := attr.Value.Any().(error); ok {
				var value string
				if debug {
					value = fmt.Sprintf("%+v", err)
				} else if stacks := errors.Stacks(err); len(stacks) > 0 {
					value = fmt.Sprintf("%s [%n at %v]", err, stacks[0], stacks[0])
				} else {
					value = err.Error()
				}
				attr.Value = slog.StringValue(value)
			}
		}
		return attr
	}
}
