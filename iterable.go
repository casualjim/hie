package hie

type Iterator[T any] func(T) bool

type Iter[T any] interface {
	HasNext() bool
	Next() T
}

type AsIter[T any] interface {
	AsIter() Iter[T]
}

func ForEach[T any](iter AsIter[T], fn Iterator[T]) {
	i := iter.AsIter()
	for i.HasNext() {
		if !fn(i.Next()) {
			break
		}
	}
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

func (m *mapperIter[T, R]) HasNext() bool {
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

type Predicate[T any] func(T) bool

func Filter[T any](iter AsIter[T], predicate Predicate[T]) AsIter[T] {
	return &filterIter[T]{
		under:     iter.AsIter(),
		predicate: predicate,
		lastMatch: None[T](),
	}
}

type filterIter[T any] struct {
	under     Iter[T]
	predicate Predicate[T]
	lastMatch Option[T]
}

func (f *filterIter[T]) HasNext() bool {
	if f.lastMatch.IsSome() { // we were already called
		return true
	}

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

func (f *filterIter[T]) AsIter() Iter[T] {
	return f
}

func Collect[T any](iter AsIter[T]) []T {
	return Fold(iter.AsIter(), nil, func(t1 []T, t2 T) ([]T, bool) {
		return append(t1, t2), true
	})
}

var _ Iter[any] = &cons[any]{}

func Concat[T any](left AsIter[T], right AsIter[T], others ...AsIter[T]) Iter[T] {
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
