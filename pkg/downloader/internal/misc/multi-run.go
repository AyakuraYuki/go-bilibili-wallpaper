package misc

import (
	"context"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/sync/errgroup"
	"log"
)

// WorkFunc is simple work func
type WorkFunc func() error

// WorkFuncWithCtx is simple work func
type WorkFuncWithCtx func(ctx context.Context) error

// SafeGo runs fs in parallels, it returns quickly when first err occurred.
// caution: need to handle timeout yourself in WorkFuncWithCtx or it may cause leak.
func SafeGo(ctx context.Context, fs ...WorkFuncWithCtx) error {
	l := len(fs)
	if l == 0 {
		return nil
	}

	errChan := make(chan error, l)

	for i := range fs {
		if fs[i] == nil {
			errChan <- nil
			continue
		}
		go func(i int) (err error) {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("%v\n", err)
					if er, ok := e.(error); ok {
						err = er
					} else {
						err = errors.New(fmt.Sprint(e))
					}
				}
				errChan <- err
			}()
			return fs[i](ctx)
		}(i)
	}

	for i := 0; i < l; i++ {
		err := <-errChan
		if err != nil {
			return err
		}
	}

	return nil
}

// MultiRun run all WorkFunc, it will return first err or nil.
// cation: It waits all WorkFunc done, err won't interrupt or notice other WorkFunc when it happened.
func MultiRun(fs ...WorkFunc) error {
	if len(fs) == 0 {
		return nil
	}

	eg := &errgroup.Group{}
	for i := range fs {
		if fs[i] == nil {
			continue
		}
		eg.Go(fs[i])
	}

	return eg.Wait()
}

// MultiRunWithCtx is mostly like MultiRun, but with ctx notify.
func MultiRunWithCtx(fs ...WorkFuncWithCtx) error {
	if len(fs) == 0 {
		return nil
	}

	eg, ctx := errgroup.WithContext(context.Background())
	for i := range fs {
		if fs[i] == nil {
			continue
		}
		j := i
		eg.Go(func() error {
			return fs[j](ctx)
		})
	}

	return eg.Wait()
}

// Hystrix is hystrix circuit breaker
func Hystrix(name string, run func() error, fallback func(error) error) error {
	callBackFunc := func(err error) error {
		var circuitError hystrix.CircuitError
		if err != nil && errors.As(err, &circuitError) {
			log.Printf("[error] name: %s, err: %v, circuit_error: %v", name, err, circuitError)
		}
		if fallback == nil {
			return err
		}
		return fallback(err)
	}
	return hystrix.Do(name, run, callBackFunc)
}
