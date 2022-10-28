package hie

type Iterator[T any] func(T) bool

func ForEach[T any](iter Iter[T], fn Iterator[T]) {
	i := iter
	for i.HasNext() {
		if !fn(i.Next()) {
			break
		}
	}
}

type Predicate[T any] func(T) bool

func Filter[T any](iter Iter[T], predicate Predicate[T]) Iter[T] {
	return &filterIter[T]{
		under:     iter,
		predicate: predicate,
	}
}

type filterIter[T any] struct {
	under     Iter[T]
	predicate Predicate[T]
	lastMatch Option[T]
}

func (f *filterIter[T]) HasNext() bool {
	for f.under.HasNext() {
		val := f.under.Next()
		if f.predicate(val) {
			f.lastMatch = Some(val)
			break
		}
	}
	return f.lastMatch.IsSome()
}

func (f *filterIter[T]) Next() T {
	if !f.lastMatch.IsSome() {
		panic("iterating beyond end")
	}
	res := f.lastMatch
	f.lastMatch = None[T]()
	return res.Value()
}

type Mapper[T, R any] func(T) R

func Map[T, R any](iter Iter[T], fn Mapper[T, R]) Iter[R] {
	return &mapperIter[T, R]{
		under:    iter,
		mapperFn: fn,
	}
}

type mapperIter[T, R any] struct {
	under    Iter[T]
	mapperFn Mapper[T, R]
}

func (m mapperIter[T, R]) HasNext() bool {
	return m.under.HasNext()
}

func (m *mapperIter[T, R]) Next() R {
	return m.mapperFn(m.under.Next())
}

type FlatMapper[T, R any] func(T) Iter[R]

func FlatMap[T, R any](iter Iter[T], fn FlatMapper[T, R]) Iter[R] {

	return &flatMapperIter[T, R]{
		mapperFn: fn,
		under:    iter,
	}
}

type flatMapperIter[T, R any] struct {
	under    Iter[T]
	mapperFn FlatMapper[T, R]
	current  Iter[R]
}

func (m *flatMapperIter[T, R]) HasNext() bool {
	return (m.current != nil && m.current.HasNext()) || m.under.HasNext()
}

func (m *flatMapperIter[T, R]) Next() R {
	if (m.current == nil || !m.current.HasNext()) && m.under.HasNext() {
		m.current = m.mapperFn(m.under.Next())
	}
	return m.current.Next()
}

func Slice[T any](value ...T) Iter[T] {
	return &sliceIter[T]{
		under: value,
		idx:   -1,
	}
}

type sliceIter[T any] struct {
	under []T
	idx   int
}

func (s sliceIter[T]) HasNext() bool {
	return s.idx < len(s.under)-1
}

func (s *sliceIter[T]) Next() T {
	if !s.HasNext() {
		panic("iterating beyond end")
	}

	s.idx++
	return s.under[s.idx]
}

func (s *sliceIter[T]) Collect() []T {
	return Collect(Slice(s.under...))
}

func Collect[T any](iter Iter[T]) []T {
	return Fold(iter, nil, func(t1 []T, t2 T) ([]T, bool) {
		return append(t1, t2), true
	})
}

type Iter[T any] interface {
	HasNext() bool
	Next() T
}

var _ Iter[any] = &cons[any]{}

func Concat[T any](left Iter[T], right Iter[T], others ...Iter[T]) Iter[T] {
	conc, ok := left.(*concat[T])
	if !ok {
		conc = &concat[T]{}
		conc.Append(left)
	}
	conc.Append(right)
	for _, other := range others {
		conc.Append(other)
	}
	return conc
}

type concat[T any] struct {
	head *cons[T]
	tail *cons[T]
}

func (c *concat[T]) Append(i Iter[T]) {
	newItem := &cons[T]{
		under: i,
	}

	if c.head == nil {
		c.head = newItem
		c.tail = newItem
	} else {
		c.tail.next = newItem
		c.tail = newItem
	}
}

func (c *concat[T]) HasNext() bool {
	return c.head != nil && c.head.HasNext()
}

func (c *concat[T]) Next() T {
	if c == nil || c.head == nil {
		panic("iterating an empty concat iterator")
	}
	return c.head.Next()
}

type cons[T any] struct {
	under Iter[T]
	next  *cons[T]
}

func (c *cons[T]) HasNext() bool {
	return c.under.HasNext() || (c.next != nil && c.next.HasNext())
}

func (c *cons[T]) Next() T {
	if c.under.HasNext() {
		return c.under.Next()
	}
	if c.next != nil && c.next.HasNext() {
		return c.next.Next()
	}
	panic("iterating beyond end")
}

func (c *cons[T]) Append(i Iter[T]) *cons[T] {
	next := &cons[T]{
		under: i,
		next:  c,
	}

	return next
}

type AccumulatorLeft[A, T any] func(A, T) (A, bool)

func Fold[A, T any](iter Iter[T], initialValue A, folder AccumulatorLeft[A, T]) A {
	it := iter
	acc := initialValue
	var shouldContinue bool
	for it.HasNext() {
		acc, shouldContinue = folder(acc, it.Next())
		if !shouldContinue {
			break
		}
	}
	return acc
}
