package log

import (
	"log/slog"
	"time"
)

const (
	attrKeyScope    = "scope"
	attrKeyDuration = "duration"
	attrKeyError    = "error"
)

func Scope(scope string) slog.Attr {
	return slog.String(attrKeyScope, scope)
}

func Dur(duration time.Duration) slog.Attr {
	return slog.String(attrKeyDuration, duration.String())
}

func Err(err error) slog.Attr {
	return slog.Any(attrKeyError, err)
}
