package hie_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/casualjim/hie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tk string

func TestFunc(t *testing.T) {
	c := context.WithValue(context.Background(), tk("key"), "doesn't matter")
	f := hie.Func(func() (int, error) {
		return 5, nil
	})

	v, ctx, e := f(c)
	if assert.NoError(t, e) {
		assert.EqualValues(t, 5, v)
		assert.Equal(t, c, ctx)
	}
}

func TestGet_Success(t *testing.T) {
	f := hie.Do(hie.Func(func() (int, error) {
		time.Sleep(20 * time.Millisecond)
		return 8, nil
	}))

	v, ctx, e := f.Get()
	if assert.NoError(t, e) {
		assert.NotNil(t, ctx)
		assert.EqualValues(t, 8, v)
	}
}

func TestGet_Success_Repeat(t *testing.T) {
	f := hie.Do(hie.Func(func() (int, error) {
		time.Sleep(20 * time.Millisecond)
		return 8, nil
	}))

	for i := 0; i < 5; i++ {
		v, ctx, e := f.Get()
		if assert.NoError(t, e) {
			assert.NotNil(t, ctx)
			assert.EqualValues(t, 8, v)
		}
	}
}

func TestGet_Error(t *testing.T) {
	exp := errors.New("expected")
	f := hie.Do(hie.Func(func() (any, error) {
		time.Sleep(20 * time.Millisecond)
		return nil, exp
	}))

	for i := 0; i < 5; i++ {
		v, ctx, e := f.Get()
		if assert.Error(t, e) {
			assert.EqualError(t, e, exp.Error())
			assert.Nil(t, v)
			assert.NotNil(t, ctx)
		}
	}
}

func TestGet_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	cf := make(chan struct{}, 1)
	f := hie.DoWithContext(ctx, func(ctx context.Context) (int, context.Context, error) {
		select {
		case <-ctx.Done():
			cf <- struct{}{}
			return 0, ctx, ctx.Err()
		case <-time.After(3 * time.Second):
			return 10, ctx, nil
		}
	})

	result, ctx2, err := f.Get()
	assert.Equal(t, 0, result)
	if assert.Error(t, err) {
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Equal(t, context.DeadlineExceeded, ctx2.Err())
		assert.Equal(t, context.DeadlineExceeded, ctx.Err())
		select {
		case <-cf:
		case <-time.After(20 * time.Millisecond):
			assert.Fail(t, "expected cancellation message")
		}
	}
}

func TestGet_Cancel(t *testing.T) {
	f := hie.Do(hie.Func(func() (int, error) {
		time.Sleep(3 * time.Second)
		return 10, nil
	}))

	go func() {
		time.Sleep(500 * time.Millisecond)
		f.Cancel()
	}()

	result, ctx, err := f.Get()
	assert.Equal(t, 0, result)
	if assert.Error(t, err) {
		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, context.Canceled, ctx.Err())
	}
}

func TestGet_Cancel_Repeat(t *testing.T) {
	f := hie.Do(hie.Func(func() (int, error) {
		time.Sleep(3 * time.Second)
		return 10, nil
	}))

	go func() {
		time.Sleep(500 * time.Millisecond)
		f.Cancel()
	}()

	result, ctx, err := f.Get()
	assert.Equal(t, 0, result)
	if assert.Error(t, err) {
		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, context.Canceled, ctx.Err())
	}
	for i := 0; i < 5; i++ {
		f.Cancel()
	}
}

func TestGet_CancelInternal(t *testing.T) {
	ctx := context.Background()

	f := hie.DoWithContext(ctx, hie.Func(func() (int, error) {
		time.Sleep(3 * time.Second)
		return 10, nil
	}))

	go func() {
		time.Sleep(500 * time.Millisecond)
		f.Cancel()
	}()

	result, ctx2, err := f.Get()
	assert.Equal(t, 0, result)
	if assert.Error(t, err) {
		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, context.Canceled, ctx2.Err())
		assert.NoError(t, ctx.Err())
	}
}

func TestGet_CancelInternalOnlyScope(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	f := hie.DoWithContext(ctx, hie.Func(func() (int, error) {
		time.Sleep(3 * time.Second)
		return 10, nil
	}))

	go func() {
		time.Sleep(500 * time.Millisecond)
		f.Cancel()
	}()

	result, ctx2, err := f.Get()
	assert.Equal(t, 0, result)
	if assert.Error(t, err) {
		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, context.Canceled, ctx2.Err())
		assert.NoError(t, ctx.Err())
		cancel()
		assert.Error(t, ctx.Err())
	}
}

func TestGet_CancelReachesFunc(t *testing.T) {
	cf := make(chan struct{}, 1)

	f := hie.Do(func(ctx context.Context) (int, context.Context, error) {
		select {
		case <-ctx.Done():
			cf <- struct{}{}
			close(cf)
			return 0, ctx, ctx.Err()
		case <-time.After(3 * time.Second):
			return 5, ctx, nil
		}
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		f.Cancel()
	}()

	val, ctx, err := f.Get()
	assert.Equal(t, 0, val)
	if assert.Error(t, err) {
		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, context.Canceled, ctx.Err())
		select {
		case <-cf:
		case <-time.After(10 * time.Millisecond): // give some time to propagate cancel
			assert.Fail(t, "expected a message from function in hie")
		}
	}
}

func TestThen(t *testing.T) {
	f := hie.Do(hie.Func(func() (int, error) {
		return 10, nil
	})).AndThen(hie.ThenFunc(func(i int) (int, error) {
		return 2 * i, nil
	})).AndThen(hie.ThenFunc(func(i int) (int, error) {
		return 2 + i, nil
	}))

	result, ctx, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 22, result)
	assert.NotNil(t, ctx)

	exp := errors.New("expected")
	g := hie.Do(hie.Func(func() (int, error) {
		return 0, exp
	})).AndThen(hie.ThenFunc(func(i int) (int, error) {
		return 2 * i, nil
	}))

	result, ctx, err = g.Get()
	require.Error(t, err)
	assert.Equal(t, 0, result)
	assert.NotNil(t, ctx)
	assert.EqualError(t, exp, err.Error())

	exp2 := errors.New("expected later")
	h := hie.Do(hie.Func(func() (int, error) {
		return 10, nil
	})).AndThen(hie.ThenFunc(func(i int) (int, error) {
		return 0, exp2
	}))
	result, ctx, err = h.Get()
	require.Error(t, err)
	assert.Equal(t, 0, result)
	assert.NotNil(t, ctx)
	assert.EqualError(t, exp2, err.Error())
}

func TestThen_Cancel(t *testing.T) {
	var count int64
	f := hie.Do(func(ctx context.Context) (int, context.Context, error) {
		atomic.AddInt64(&count, 1)
		select {
		case <-ctx.Done():
			return 0, ctx, ctx.Err()
		case <-time.After(1 * time.Second):
			return 10, ctx, nil
		}
	}).AndThen(func(ctx context.Context, i int) (int, context.Context, error) {
		atomic.AddInt64(&count, 1)
		select {
		case <-ctx.Done():
			return 0, ctx, ctx.Err()
		case <-time.After(1 * time.Second):
			return 20, ctx, nil
		}
	})

	go func() {
		time.Sleep(200 * time.Millisecond)
		f.Cancel()
	}()

	result, _, err := f.Get()
	require.Error(t, err)
	assert.Equal(t, 0, result)
	assert.EqualValues(t, 1, atomic.LoadInt64(&count))

}

func TestOrElse(t *testing.T) {
	f := hie.Do(hie.Func(func() (int, error) {
		return 0, errors.New("transient")
	})).OrElse(hie.ElseFunc(func(i error) (int, error) {
		return 20, nil
	}))

	result, ctx, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 20, result)
	assert.NotNil(t, ctx)

	exp := errors.New("expected")
	g := hie.Do(hie.Func(func() (int, error) {
		return 0, exp
	})).OrElse(hie.ElseFunc(func(i error) (int, error) {
		return 21, nil
	})).OrElse(hie.ElseFunc(func(i error) (int, error) {
		return 22, nil // should not execute because previous is success
	}))

	result, ctx, err = g.Get()
	require.NoError(t, err)
	assert.Equal(t, 21, result)
	assert.NotNil(t, ctx)

}
func TestOrElse_Cancel(t *testing.T) {
	var count int64
	f := hie.Do(func(ctx context.Context) (int, context.Context, error) {
		atomic.AddInt64(&count, 1)
		select {
		case <-ctx.Done():
			return 0, ctx, ctx.Err()
		case <-time.After(1 * time.Second):
			return 10, ctx, nil
		}
	}).OrElse(func(ctx context.Context, err error) (int, context.Context, error) {
		atomic.AddInt64(&count, 1)
		select {
		case <-ctx.Done():
			return 0, ctx, ctx.Err()
		case <-time.After(1 * time.Second):
			return 20, ctx, nil
		}
	})

	go func() {
		time.Sleep(200 * time.Millisecond)
		f.Cancel()
	}()

	result, _, err := f.Get()
	require.Error(t, err)
	assert.Equal(t, 0, result)
	assert.EqualValues(t, 1, atomic.LoadInt64(&count))

}
