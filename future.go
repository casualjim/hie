package hie

import (
	"context"
	"sync"
)

type result[T any] struct {
	Value T
	Err   error
	Ctx   context.Context
	_     struct{} // avoid unkeyed usage
}

// Future represents a value that will become available in the future.
//
// This is loosely based on this paper: http://www.home.hs-karlsruhe.de/~suma0002/publications/events-to-futures.pdf
type Future[T any] interface {
	AndThen(func(context.Context, T) (T, context.Context, error)) Future[T]
	OrElse(func(context.Context, error) (T, context.Context, error)) Future[T]
	Get() (T, context.Context, error)
	Cancel()
}

// Func wraps a nullary function with a context handling function
func Func[T any](fn func() (T, error)) func(context.Context) (T, context.Context, error) {
	return func(ctx context.Context) (T, context.Context, error) {
		v, e := fn()
		return v, ctx, e
	}
}

// ThenFunc wraps a continuation function with a context handling function
func ThenFunc[T any](fn func(T) (T, error)) func(context.Context, T) (T, context.Context, error) {
	return func(ctx context.Context, vo T) (T, context.Context, error) {
		v, e := fn(vo)
		return v, ctx, e
	}
}

// ElseFunc wraps an error handler continuation with a context handling function
func ElseFunc[T any](fn func(error) (T, error)) func(context.Context, error) (T, context.Context, error) {
	return func(ctx context.Context, err error) (T, context.Context, error) {
		v, e := fn(err)
		return v, ctx, e
	}
}

// Do creates a future that executes the function in a go routine
// The context that is passed into the function will provide a cancellation signal
func Do[T any](fn func(context.Context) (T, context.Context, error)) Future[T] {
	return DoWithContext(context.Background(), fn)
}

// DoWithContext creates a future that executes the function in a go routine.
// The context is passed into the function so that it can be used for handling cancellation
// The function is expected to either pass the original context along or provide a new context based off the one passed in.
func DoWithContext[T any](ctx context.Context, fn func(context.Context) (T, context.Context, error)) Future[T] {
	inner, cancel := context.WithCancel(ctx)
	c := make(chan result[T], 1)
	go func() {
		defer close(c)
		v, ctx, e := fn(inner)
		c <- result[T]{Value: v, Err: e, Ctx: ctx}
	}()

	return &future[T]{
		C:      c,
		cancel: cancel,
		ctx:    inner,
		once:   new(sync.Once),
	}
}

type future[T any] struct {
	C      chan result[T]
	cancel context.CancelFunc
	ctx    context.Context
	once   *sync.Once
	val    *result[T]
}

func (f *future[T]) Get() (T, context.Context, error) {
	f.once.Do(func() {
		select {
		case <-f.ctx.Done():
			f.val = &result[T]{Err: f.ctx.Err(), Ctx: f.ctx}
		case val := <-f.C:
			f.val = &val
		}
	})
	if f.val == nil {
		var zeroT T
		return zeroT, f.ctx, nil
	}
	return f.val.Value, f.val.Ctx, f.val.Err
}

func (f *future[T]) Cancel() {
	f.cancel()
}

func (f *future[T]) AndThen(fn func(context.Context, T) (T, context.Context, error)) Future[T] {
	c := make(chan result[T], 1)
	go func() {
		defer close(c)

		v, ctx, e := f.Get()
		if e != nil { // on error we fail here
			c <- result[T]{Value: v, Ctx: ctx, Err: e}
			return
		}

		select {
		case <-ctx.Done():
			c <- result[T]{Value: v, Ctx: ctx, Err: ctx.Err()}
			select {
			case <-f.ctx.Done():
			default:
				f.cancel() // ensure closed if out of scope
			}
		case <-f.ctx.Done():
			c <- result[T]{Value: v, Ctx: ctx, Err: f.ctx.Err()}
		default:
			vv, ctx2, e2 := fn(ctx, v)
			c <- result[T]{Value: vv, Ctx: ctx2, Err: e2}
		}

	}()
	return &future[T]{
		C:      c,
		cancel: f.cancel,
		ctx:    f.ctx,
		once:   new(sync.Once),
	}
}

func (f *future[T]) OrElse(fn func(context.Context, error) (T, context.Context, error)) Future[T] {
	c := make(chan result[T], 1)
	go func() {
		defer close(c)

		v, ctx, e := f.Get()
		if e == nil { // on error we fail here
			c <- result[T]{Value: v, Ctx: ctx}
			return
		}

		select {
		case <-ctx.Done():
			c <- result[T]{Value: v, Ctx: ctx, Err: ctx.Err()}
			select {
			case <-f.ctx.Done():
			default:
				f.cancel() // ensure closed if out of scope
			}
		case <-f.ctx.Done():
			c <- result[T]{Value: v, Ctx: ctx, Err: f.ctx.Err()}
		default:
			vv, ctx2, e2 := fn(ctx, e)
			c <- result[T]{Value: vv, Ctx: ctx2, Err: e2}
		}

	}()
	return &future[T]{
		C:      c,
		cancel: f.cancel,
		ctx:    f.ctx,
		once:   new(sync.Once),
	}
}
