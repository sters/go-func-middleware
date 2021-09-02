package middleware

import (
	"context"
	"sync"
	"time"
)

type Function func(ctx context.Context) error
type Middleware func(ctx context.Context, next Function) error

type ApplySide int

const (
	ApplySideBefore = iota
	ApplySideAfter
	ApplySideAfterNoErrorOnly
	ApplySideAfterErrorOnly
	ApplySideBoth
)

func Invoke(ctx context.Context, middlewares []Middleware, root Function) error {
	wrapper := func(m Middleware, current Function) Function {
		return func(ctx context.Context) error {
			return m(ctx, current)
		}
	}

	current := root
	for i := len(middlewares) - 1; i >= 0; i-- {
		current = wrapper(middlewares[i], current)
	}

	return current(ctx)
}

func WithTimeout(d time.Duration) Middleware {
	return func(ctx context.Context, next Function) error {
		ctx, cancel := context.WithTimeout(ctx, d)
		defer cancel()
		return next(ctx)
	}
}

func WithLock(mux *sync.Mutex) Middleware {
	return func(ctx context.Context, next Function) error {
		mux.Lock()
		defer mux.Unlock()
		return next(ctx)
	}
}

func WithRWLock(mux *sync.RWMutex) Middleware {
	return func(ctx context.Context, next Function) error {
		mux.RLock()
		defer mux.RUnlock()
		return next(ctx)
	}
}

func WithSleep(d time.Duration, a ApplySide) Middleware {
	return func(ctx context.Context, next Function) error {
		if a == ApplySideBefore || a == ApplySideBoth {
			time.Sleep(d)
		}

		result := next(ctx)

		if a == ApplySideAfter ||
			a == ApplySideBoth ||
			((result != nil) == (a == ApplySideAfterErrorOnly)) ||
			((result == nil) == (a == ApplySideAfterNoErrorOnly)) {
			time.Sleep(d)
		}

		return result
	}
}
