# go-func-middleware

[![go](https://github.com/sters/go-func-middleware/workflows/Go/badge.svg)](https://github.com/sters/go-func-middleware/actions?query=workflow%3AGo)
[![codecov](https://codecov.io/gh/sters/go-func-middleware/branch/main/graph/badge.svg)](https://codecov.io/gh/sters/go-func-middleware)
[![go-report](https://goreportcard.com/badge/github.com/sters/go-func-middleware)](https://goreportcard.com/report/github.com/sters/go-func-middleware)

More focus on your application core logic.

## Usage

```go
type yourApplication struct {
  middlewares []middleware.Middleware
}

func (app *yourApplication) doSomethingHandler(ctx context.Context) error {
  app.doSomething1()

  err := middleware.Invoke(ctx, app.middlewares, func(ctx context.Context) error {
    return app.doSomething2(ctx)
  })
  if err != nil {
    return err
  }

  return nil
}
```
