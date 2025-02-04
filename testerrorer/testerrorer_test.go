package testerrorer_test

import (
	"os"
	"testing"

	"github.com/vikstrous/slogexamples/testerrorer"
	"golang.org/x/exp/slog"
)

type TestT struct {
	*testing.T
	DidError bool
}

func (t *TestT) Errorf(s string, args ...any) {
	t.DidError = true
}

func TestErrorerCallsNext(t *testing.T) {
	nextWasCalled := false
	next := func(groups []string, a slog.Attr) slog.Attr {
		nextWasCalled = true
		return a
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		ReplaceAttr: testerrorer.NewTestErrorer(nil, next),
	}))
	logger.Info("example")
	if !nextWasCalled {
		t.Fatal("next was not called")
	}
}

func TestErrorerErrors(t *testing.T) {
	wrappedT := &TestT{T: t}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		ReplaceAttr: testerrorer.NewTestErrorer(wrappedT, nil),
	}))
	allocsPerRun := testing.AllocsPerRun(1, func() {
		logger.Log(nil, slog.LevelInfo, "example")
	})
	if allocsPerRun > 1 {
		// 1 allocation comes from the use of ReplaceAttr, so we can't reach 0
		t.Fatalf("extra allocations introduced in info path %.0f", allocsPerRun)
	}
	if wrappedT.DidError {
		t.Fatal("error below error level")
	}
	allocsPerRun = testing.AllocsPerRun(1, func() {
		logger.Error("example")
	})
	if !wrappedT.DidError {
		t.Fatal("did not error at error level")
	}
	if allocsPerRun > 1 {
		// 1 allocation comes from the use of ReplaceAttr, so we can't reach 0
		t.Fatalf("extra allocations introduced in error path: %.0f", allocsPerRun)
	}
	wrappedT.DidError = false
	logger.Log(nil, slog.LevelError+1, "example")
	if !wrappedT.DidError {
		t.Fatal("did not error above error level")
	}
	// Fialing the test in this case is actually a limitation of this approach rather than desired behaviour.
	wrappedT.DidError = false
	logger.Info("example", "l", slog.LevelError)
	if !wrappedT.DidError {
		t.Fatal("error level attribute in list didn't trigger a test fialure")
	}
}

func ExampleNewTestErrorer() {
	var t *testing.T
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		ReplaceAttr: testerrorer.NewTestErrorer(t, nil),
	}))
	logger.Info("use the logger as normal")
}
