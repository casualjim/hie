package iter

import (
	"github.com/casualjim/hie"
	"github.com/casualjim/hie/opt"
)

func Empty[T any]() hie.Iter[T] {
	return emptyIter[T]{}
}

type emptyIter[T any] struct{}

func (emptyIter[T]) HasNext() bool { return false }
func (emptyIter[T]) Next() T       { panic("next called on an empty iter") }

func ForEach[T any](iter hie.Iter[T], fn hie.Iterator[T]) {
	i := iter
	for i.HasNext() {
		if !fn(i.Next()) {
			break
		}
	}
}

func Flatten[T any](iter hie.Iter[hie.Iter[T]]) hie.Iter[T] {
	return FlatMap(iter, hie.Identity[hie.Iter[T]])
}

type FilterMapper[T, R any] func(T) (R, bool)

func FilterMap[T, R any](iter hie.Iter[T], fn FilterMapper[T, R]) hie.Iter[R] {
	return &filterMapperIter[T, R]{
		under:     iter,
		mapperFn:  fn,
		lastMatch: opt.None[R](),
	}
}

type filterMapperIter[T, R any] struct {
	under     hie.Iter[T]
	mapperFn  FilterMapper[T, R]
	lastMatch opt.Option[R]
}

func (f *filterMapperIter[T, R]) HasNext() bool {
	if f.lastMatch.IsSome() {
		return true
	}

	for f.under.HasNext() {
		elem := f.under.Next()
		if nv, ok := f.mapperFn(elem); ok {
			f.lastMatch = opt.Some(nv)
			return true
		}
	}
	return false
}

func (f *filterMapperIter[T, R]) Next() R {
	res := f.lastMatch
	f.lastMatch = opt.None[R]()
	return res.Value()
}

type Mapper[T, R any] func(T) R

func Map[T, R any](iter hie.Iter[T], fn Mapper[T, R]) hie.Iter[R] {
	return &mapperIter[T, R]{
		under:    iter,
		mapperFn: fn,
	}
}

type mapperIter[T, R any] struct {
	under    hie.Iter[T]
	mapperFn Mapper[T, R]
}

func (m *mapperIter[T, R]) HasNext() bool {
	return m.under.HasNext()
}

func (m *mapperIter[T, R]) Next() R {
	return m.mapperFn(m.under.Next())
}

type FlatMapper[T, R any] func(T) hie.Iter[R]

func FlatMap[T, R any](iter hie.Iter[T], fn FlatMapper[T, R]) hie.Iter[R] {

	return &flatMapperIter[T, R]{
		mapperFn: fn,
		under:    iter,
	}
}

type flatMapperIter[T, R any] struct {
	under    hie.Iter[T]
	mapperFn FlatMapper[T, R]
	current  hie.Iter[R]
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

type Predicate[T any] func(T) bool

func Filter[T any](iter hie.Iter[T], predicate Predicate[T]) hie.Iter[T] {
	return &filterIter[T]{
		under:     iter,
		predicate: predicate,
		lastMatch: opt.None[T](),
	}
}

type filterIter[T any] struct {
	under     hie.Iter[T]
	predicate Predicate[T]
	lastMatch opt.Option[T]
}

func (f *filterIter[T]) HasNext() bool {
	if f.lastMatch.IsSome() { // we were already called
		return true
	}

	for f.under.HasNext() {
		val := f.under.Next()
		if f.predicate(val) {
			f.lastMatch = opt.Some(val)
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
	f.lastMatch = opt.None[T]()
	return res.Value()
}

func Collect[T any](iter hie.Iter[T]) []T {
	return Fold(iter, nil, func(t1 []T, t2 T) ([]T, bool) {
		return append(t1, t2), true
	})
}

var _ hie.Iter[any] = &cons[any]{}

func Concat[T any](left hie.Iter[T], right hie.Iter[T], others ...hie.Iter[T]) hie.Iter[T] {
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

func (c *concat[T]) Append(i hie.Iter[T]) {
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
	under hie.Iter[T]
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

func (c *cons[T]) Append(i hie.Iter[T]) *cons[T] {
	next := &cons[T]{
		under: i,
		next:  c,
	}

	return next
}

type AccumulatorLeft[A, T any] func(A, T) (A, bool)

func Fold[A, T any](iter hie.Iter[T], initialValue A, folder AccumulatorLeft[A, T]) A {
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

func TakeN[T any](iter hie.Iter[T], n int) hie.Iter[T] {
	return &takeNIter[T]{
		max:   n,
		under: iter,
	}
}

type takeNIter[T any] struct {
	max   int
	count int
	under hie.Iter[T]
}

func (n *takeNIter[T]) HasNext() bool {
	return n.count < n.max && n.under.HasNext()
}

func (n *takeNIter[T]) Next() T {
	elem := n.under.Next()
	n.count++
	return elem
}
