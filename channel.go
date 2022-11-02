package hie

func Chan[T any](under <-chan T) Iter[T] {
	return &chanIter[T]{
		ch: under,
	}
}

type chanIter[T any] struct {
	ch        <-chan T
	lastMatch T
}

func (c *chanIter[T]) HasNext() bool {
	val, closed := <-c.ch
	if !closed {
		return false
	}
	c.lastMatch = val
	return true
}

func (c *chanIter[T]) Next() T {
	lm := c.lastMatch
	var zero T
	c.lastMatch = zero
	return lm
}
