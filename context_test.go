package context_test

import (
	"sync"
	"testing"
	"time"

	"github.com/posener/context"
)

func TestSetTimeout(t *testing.T) {
	t.Parallel()

	ctx := context.Init()
	setCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	context.Set(setCtx)

	var wg sync.WaitGroup
	wg.Add(3)

	context.Go(func() {
		assertWithDeadline(t)
		wg.Done()
	})

	context.GoCtx(setCtx, func() {
		assertWithDeadline(t)
		wg.Done()
	})

	context.GoCtx(ctx, func() {
		assertNoDeadline(t)
		wg.Done()
	})

	wg.Wait()
}

func TestNoSetTimeout(t *testing.T) {
	t.Parallel()
	ctx := context.Init()
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	context.Go(func() {
		assertNoDeadline(t)
		wg.Done()
	})

	context.GoCtx(ctx, func() {
		assertWithDeadline(t)
		wg.Done()
	})

	wg.Wait()
}

func assertWithDeadline(t *testing.T) {
	t.Helper()
	ctx := context.Get()
	if _, ok := ctx.Deadline(); !ok {
		t.Error("deadline did not propagated")
	}
	select {
	case <-ctx.Done():
	case <-time.After(500 * time.Millisecond):
		t.Error("deadline did not propagated")
	}
}

func assertNoDeadline(t *testing.T) {
	t.Helper()
	ctx := context.Get()
	if _, ok := ctx.Deadline(); ok {
		t.Error("no deadline was defined")
	}
	select {
	case <-ctx.Done():
		t.Error("no deadline was defined")
	case <-time.After(100 * time.Millisecond):
	}
}
