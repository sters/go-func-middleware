package middleware

import (
	"context"
	"fmt"
	"testing"
)

var calledKey = struct{}{}

type called struct {
	num int
}

func getCalledCount(ctx context.Context) (int, error) {
	c, ok := ctx.Value(calledKey).(*called)
	if !ok {
		return 0, fmt.Errorf("failed to get context value as *called")
	}
	return c.num, nil
}

func incrementCalled(ctx context.Context) error {
	c, ok := ctx.Value(calledKey).(*called)
	if !ok {
		return fmt.Errorf("failed to get context value as *called")
	}
	c.num++
	return nil
}

func initCtx(c *called) context.Context {
	return context.WithValue(
		context.Background(),
		calledKey,
		c,
	)
}

func TestInvoke_empty(t *testing.T) {
	c := &called{}
	err := Invoke(
		initCtx(c),
		nil,
		//nolint
		func(ctx context.Context) error {
			return incrementCalled(ctx)
		},
	)

	if err != nil {
		t.Errorf("got error: %+v", err)
	}

	if c.num != 1 {
		t.Errorf("root function does not called")
	}
}

func TestInvoke(t *testing.T) {
	c := &called{}
	err := Invoke(
		initCtx(c),
		[]Middleware{
			func(ctx context.Context, next Function) error {
				if got, err := getCalledCount(ctx); err != nil {
					return err
				} else if want := 0; got != want {
					return fmt.Errorf("want = %v, got = %v", want, got)
				}

				if err := incrementCalled(ctx); err != nil {
					return err
				}

				err := next(ctx)
				if err != nil {
					return err
				}

				if got, err := getCalledCount(ctx); err != nil {
					return err
				} else if want := 6; got != want {
					return fmt.Errorf("want = %v, got = %v", want, got)
				}

				if err := incrementCalled(ctx); err != nil {
					return err
				}

				return nil
			},
			func(ctx context.Context, next Function) error {
				if got, err := getCalledCount(ctx); err != nil {
					return err
				} else if want := 1; got != want {
					return fmt.Errorf("want = %v, got = %v", want, got)
				}

				if err := incrementCalled(ctx); err != nil {
					return err
				}

				err := next(ctx)
				if err != nil {
					return err
				}

				if got, err := getCalledCount(ctx); err != nil {
					return err
				} else if want := 5; got != want {
					return fmt.Errorf("want = %v, got = %v", want, got)
				}

				if err := incrementCalled(ctx); err != nil {
					return err
				}

				return nil
			},
			func(ctx context.Context, next Function) error {
				if got, err := getCalledCount(ctx); err != nil {
					return err
				} else if want := 2; got != want {
					return fmt.Errorf("want = %v, got = %v", want, got)
				}

				if err := incrementCalled(ctx); err != nil {
					return err
				}

				err := next(ctx)
				if err != nil {
					return err
				}

				if got, err := getCalledCount(ctx); err != nil {
					return err
				} else if want := 4; got != want {
					return fmt.Errorf("want = %v, got = %v", want, got)
				}

				if err := incrementCalled(ctx); err != nil {
					return err
				}

				return nil
			},
		},
		func(ctx context.Context) error {
			if got, err := getCalledCount(ctx); err != nil {
				return err
			} else if want := 3; got != want {
				return fmt.Errorf("want = %v, got = %v", want, got)
			}

			return incrementCalled(ctx)
		},
	)

	if err != nil {
		t.Errorf("got error: %+v", err)
	}

	if want := 7; c.num != want {
		t.Errorf("want = %+v, got = %+v", want, c.num)
	}
}
