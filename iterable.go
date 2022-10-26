package hie

type AsIter[T any] interface {
	AsIter() Iter[T]
}

type Iterator[T any] func(T) bool

func ForEach[T any](iter AsIter[T], fn Iterator[T]) {
	i := iter.AsIter()
	for i.HasNext() {
		if !fn(i.Next()) {
			break
		}
	}
}

type Predicate[T any] func(T) bool

func Filter[T any](iter AsIter[T], predicate Predicate[T]) AsIter[T] {
	return &filterIter[T]{
		under:     iter.AsIter(),
		predicate: predicate,
	}
}

type filterIter[T any] struct {
	under     Iter[T]
	predicate Predicate[T]
	lastMatch T
	matched   bool
}

func (f *filterIter[T]) HasNext() bool {
	for f.under.HasNext() {
		val := f.under.Next()
		if f.predicate(val) {
			f.lastMatch = val
			f.matched = true
			break
		}
	}
	return f.matched
}

func (f *filterIter[T]) Next() T {
	if !f.matched {
		panic("iterating beyond end")
	}
	res := f.lastMatch
	f.matched = false
	return res
}

func (f *filterIter[T]) AsIter() Iter[T] {
	return f
}

type Mapper[T, R any] func(T) R

func Map[T, R any](iter AsIter[T], fn Mapper[T, R]) AsIter[R] {
	return &mapperIter[T, R]{
		under:    iter.AsIter(),
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

func (m *mapperIter[T, R]) AsIter() Iter[R] {
	return m
}

type FlatMapper[T, R any] func(T) AsIter[R]

func FlatMap[T, R any](iter AsIter[T], fn FlatMapper[T, R]) AsIter[R] {

	return &flatMapperIter[T, R]{
		mapperFn: fn,
		under:    iter.AsIter(),
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
		m.current = m.mapperFn(m.under.Next()).AsIter()
	}
	return m.current.Next()
}

func (m *flatMapperIter[T, R]) AsIter() Iter[R] {
	return m
}

func Slice[T any](value ...T) AsIter[T] {
	return SliceAsIter[T](value)
}

type SliceAsIter[T any] []T

func (s SliceAsIter[T]) AsIter() Iter[T] {
	return &sliceIter[T]{
		under: []T(s),
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

func Collect[T any](iter AsIter[T]) []T {
	return Fold(iter, nil, func(t1 []T, t2 T) ([]T, bool) {
		return append(t1, t2), true
	})
}

type Iter[T any] interface {
	HasNext() bool
	Next() T
}

var _ Iter[any] = &cons[any]{}
var _ AsIter[any] = &cons[any]{}

func Concat[T any](left AsIter[T], right AsIter[T], others ...AsIter[T]) AsIter[T] {
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

func (c *concat[T]) Append(i AsIter[T]) {
	newItem := &cons[T]{
		under: i.AsIter(),
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

func (c *concat[T]) AsIter() Iter[T] {
	return c
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

func (c *cons[T]) Append(i AsIter[T]) *cons[T] {
	next := &cons[T]{
		under: i.AsIter(),
		next:  c,
	}

	return next
}

func (c *cons[T]) AsIter() Iter[T] {
	return c
}

type AccumulatorLeft[A, T any] func(A, T) (A, bool)

func Fold[A, T any](iter AsIter[T], initialValue A, folder AccumulatorLeft[A, T]) A {
	it := iter.AsIter()
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
