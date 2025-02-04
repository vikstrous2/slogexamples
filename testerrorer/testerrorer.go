package testerrorer

import (
	"testing"

	"golang.org/x/exp/slog"
)

type testErrorer struct {
	TB   testing.TB
	Next func(groups []string, a slog.Attr) slog.Attr
}

func (t *testErrorer) replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Value.Kind() == slog.KindAny {
		level, ok := a.Value.Any().(slog.Level)
		if ok && level >= slog.LevelError {
			t.TB.Errorf("An error was logged.")
		}
	}
	if t.Next == nil {
		return a
	}
	return t.Next(groups, a)
}

// NewTestErrorer creates a function that can be used as slog's ReplaceAttr function in the standard library implementations of TextHandler or JSONHandler
func NewTestErrorer(tb testing.TB, next func(groups []string, a slog.Attr) slog.Attr) func(groups []string, a slog.Attr) slog.Attr {
	return (&testErrorer{tb, next}).replaceAttr
}
