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
	ii := filterMapperIter[T, R]{
		under:     iter,
		mapperFn:  fn,
		lastMatch: opt.None[R](),
	}

	if IsClonable(iter) {
		ic := clonableFilterMapper[T, R]{
			filterMapperIter: ii,
		}
		if IsClosable(iter) {
			return &clonableClosableFilterMapperIter[T, R]{
				clonableFilterMapper: ic,
			}
		}
		return &ic
	}

	if IsClosable(iter) {
		return &closableFilterMapperIter[T, R]{
			filterMapperIter: ii,
		}
	}
	return &ii
}

type closableFilterMapperIter[T, R any] struct {
	filterMapperIter[T, R]
}

func (c *closableFilterMapperIter[T, R]) Close() error {
	return Close(hie.Iter[R](&c.filterMapperIter))
}

type clonableClosableFilterMapperIter[T, R any] struct {
	clonableFilterMapper[T, R]
}

func (c *clonableClosableFilterMapperIter[T, R]) Close() error {
	return Close(hie.Iter[R](&c.filterMapperIter))
}

type clonableFilterMapper[T, R any] struct {
	filterMapperIter[T, R]
}

func (c *clonableFilterMapper[T, R]) Clone() hie.Iter[R] {
	cu, cloned := Clone(c.under)
	if !cloned {
		panic("Clone called on an unclonable iterator")
	}
	return &clonableFilterMapper[T, R]{
		filterMapperIter: filterMapperIter[T, R]{
			under:    cu,
			mapperFn: c.mapperFn,
		},
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
	mi := mapperIter[T, R]{
		under:    iter,
		mapperFn: fn,
	}

	if IsClonable(iter) {
		cm := clonableMapper[T, R]{
			mapperIter: mi,
		}
		if IsClosable(iter) {
			return &clonableClosableMapperIter[T, R]{
				clonableMapper: cm,
			}
		}
		return &cm
	}

	if IsClosable(iter) {
		return &closableMapperIter[T, R]{
			mapperIter: mi,
		}
	}
	return &mi
}

type clonableMapper[T, R any] struct {
	mapperIter[T, R]
}

func (c *clonableMapper[T, R]) Clone() hie.Iter[R] {
	cu, cloned := Clone(c.under)
	if !cloned {
		panic("Clone called on an unclonable iterator")
	}
	return &clonableMapper[T, R]{
		mapperIter: mapperIter[T, R]{
			under:    cu,
			mapperFn: c.mapperFn,
		},
	}
}

type closableMapperIter[T, R any] struct {
	mapperIter[T, R]
}

func (c *closableMapperIter[T, R]) Close() error {
	return Close(hie.Iter[R](&c.mapperIter))
}

type clonableClosableMapperIter[T, R any] struct {
	clonableMapper[T, R]
}

func (c *clonableClosableMapperIter[T, R]) Close() error {
	return Close(hie.Iter[R](&c.clonableMapper))
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

	fm := flatMapperIter[T, R]{
		mapperFn: fn,
		under:    iter,
	}

	if IsClonable(iter) {
		cm := clonableFlatMapper[T, R]{
			flatMapperIter: fm,
		}
		if IsClosable(iter) {
			return &clonableClosableFlatMapperIter[T, R]{
				clonableFlatMapper: cm,
			}
		}
		return &cm
	}

	if IsClosable(iter) {
		return &closableFlatMapperIter[T, R]{
			flatMapperIter: fm,
		}
	}
	return &fm
}

type clonableFlatMapper[T, R any] struct {
	flatMapperIter[T, R]
}

func (c *clonableFlatMapper[T, R]) Clone() hie.Iter[R] {
	cu, cloned := Clone(c.under)
	if !cloned {
		panic("Clone called on an unclonable iterator")
	}
	return &clonableFlatMapper[T, R]{
		flatMapperIter: flatMapperIter[T, R]{
			under:    cu,
			mapperFn: c.mapperFn,
		},
	}
}

type closableFlatMapperIter[T, R any] struct {
	flatMapperIter[T, R]
}

func (c *closableFlatMapperIter[T, R]) Close() error {
	return Close(hie.Iter[R](&c.flatMapperIter))
}

type clonableClosableFlatMapperIter[T, R any] struct {
	clonableFlatMapper[T, R]
}

func (c *clonableClosableFlatMapperIter[T, R]) Close() error {
	return Close(hie.Iter[R](&c.clonableFlatMapper))
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
	fi := filterIter[T]{
		under:     iter,
		predicate: predicate,
		lastMatch: opt.None[T](),
	}

	if IsClonable(iter) {
		cm := clonableFilterIter[T]{
			filterIter: fi,
		}
		if IsClosable(iter) {
			return &clonableClosableFilterIter[T]{
				clonableFilterIter: cm,
			}
		}
		return &cm
	}

	if IsClosable(iter) {
		return &closableFilterIter[T]{
			filterIter: fi,
		}
	}
	return &fi
}

type clonableFilterIter[T any] struct {
	filterIter[T]
}

func (c *clonableFilterIter[T]) Clone() hie.Iter[T] {
	cu, cloned := Clone(c.under)
	if !cloned {
		panic("Clone called on an unclonable iterator")
	}
	return &clonableFilterIter[T]{
		filterIter: filterIter[T]{
			under:     cu,
			predicate: c.predicate,
		},
	}
}

type closableFilterIter[T any] struct {
	filterIter[T]
}

func (c *closableFilterIter[T]) Close() error {
	return Close(hie.Iter[T](&c.filterIter))
}

type clonableClosableFilterIter[T any] struct {
	clonableFilterIter[T]
}

func (c *clonableClosableFilterIter[T]) Close() error {
	return Close(hie.Iter[T](&c.clonableFilterIter))
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

type clonableConcatIter[T any] struct {
	head *clonableConsIter[T]
	tail *clonableConsIter[T]
}

func (c *clonableConcatIter[T]) Clone() hie.Iter[T] {
	cui, cloned := Clone(hie.Iter[T](c.head))
	if !cloned {
		panic("Clone called on an unclonable iterator")
	}

	head := cui.(*clonableConsIter[T])
	var tail *clonableConsIter[T]
	var cur = head
	for {
		if cur == c.tail {
			tail = cur
			break
		}
		prev := cur
		cur = cur.next.Clone().(*clonableConsIter[T])
		prev.next = cur
	}
	return &clonableConcatIter[T]{
		head: head,
		tail: tail,
	}
}

func (c *clonableConcatIter[T]) HasNext() bool {
	return c.head != nil && c.head.HasNext()
}

func (c *clonableConcatIter[T]) Next() T {
	if c == nil || c.head == nil {
		panic("iterating an empty clonable concat iterator")
	}
	return c.head.Next()
}

type closableConcatIter[T any] struct {
	head *closableConcatIter[T]
	tail *closableConcatIter[T]
}

func (c *closableConcatIter[T]) HasNext() bool {
	return c.head != nil && c.head.HasNext()
}

func (c *closableConcatIter[T]) Next() T {
	if c == nil || c.head == nil {
		panic("iterating an empty clonable concat iterator")
	}
	return c.head.Next()
}

func (c *closableConcatIter[T]) Close() error {
	head := c.head
	cur := head
	for {
		if cur == tail {
			break
		}
	}
	return Close(hie.Iter[T](c.head))
}

type clonableClosableConcatIter[T any] struct {
	clonableConcatIter[T]
}

func (c *clonableClosableConcatIter[T]) Close() error {
	return Close(hie.Iter[T](&c.clonableConcatIter))
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

type clonableConsIter[T any] struct {
	under hie.Iter[T]
	next  *clonableConsIter[T]
}

func (c *clonableConsIter[T]) Clone() hie.Iter[T] {
	cu, cloned := Clone(c.under)
	if !cloned {
		panic("Clone called on an unclonable iterator")
	}

	return &clonableConsIter[T]{
		under: cu,
		next:  c.next.Clone().(*clonableConsIter[T]),
	}
}

func (c *clonableConsIter[T]) HasNext() bool {
	return c.under.HasNext() || (c.next != nil && c.next.HasNext())
}

func (c *clonableConsIter[T]) Next() T {
	if c.under.HasNext() {
		return c.under.Next()
	}
	if c.next != nil && c.next.HasNext() {
		return c.next.Next()
	}
	panic("iterating beyond end")
}

func (c *clonableConsIter[T]) Append(i hie.Iter[T]) *clonableConsIter[T] {
	next := &clonableConsIter[T]{
		under: i,
		next:  c,
	}

	return next
}

type closableConsIter[T any] struct {
	under hie.Iter[T]
	next  *closableConsIter[T]
}

func (c *closableConsIter[T]) Close() error {
	if c.next != nil {
		if e := c.next.Close(); e != nil {
			return e
		}
	}
	err := Close(c.under)
	return err
}

func (c *closableConsIter[T]) HasNext() bool {
	return c.under.HasNext() || (c.next != nil && c.next.HasNext())
}

func (c *closableConsIter[T]) Next() T {
	if c.under.HasNext() {
		return c.under.Next()
	}
	if c.next != nil && c.next.HasNext() {
		return c.next.Next()
	}
	panic("iterating beyond end")
}

func (c *closableConsIter[T]) Append(i hie.Iter[T]) *closableConsIter[T] {
	next := &closableConsIter[T]{
		under: i,
		next:  c,
	}

	return next
}

type clonableClosableConsIter[T any] struct {
	under hie.Iter[T]
	next  *clonableClosableConsIter[T]
}

func (c *clonableClosableConsIter[T]) Clone() hie.Iter[T] {
	cu, cloned := Clone(c.under)
	if !cloned {
		panic("Clone called on an unclonable iterator")
	}

	return &clonableClosableConsIter[T]{
		under: cu,
		next:  c.next.Clone().(*clonableClosableConsIter[T]),
	}
}

func (c *clonableClosableConsIter[T]) Close() error {
	err := Close(c.under)
	if e := c.next.Close(); e != nil {
		return e
	}
	return err
}

func (c *clonableClosableConsIter[T]) HasNext() bool {
	return c.under.HasNext() || (c.next != nil && c.next.HasNext())
}

func (c *clonableClosableConsIter[T]) Next() T {
	if c.under.HasNext() {
		return c.under.Next()
	}
	if c.next != nil && c.next.HasNext() {
		return c.next.Next()
	}
	panic("iterating beyond end")
}

func (c *clonableClosableConsIter[T]) Append(i hie.Iter[T]) *clonableClosableConsIter[T] {
	next := &clonableClosableConsIter[T]{
		under: i,
		next:  c,
	}

	return next
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
	tn := takeNIter[T]{
		max:   n,
		under: iter,
	}

	if IsClonable(iter) {
		cm := clonableTakeNIter[T]{
			takeNIter: tn,
		}
		if IsClosable(iter) {
			return &clonableClosableTakeNIter[T]{
				clonableTakeNIter: cm,
			}
		}
		return &cm
	}

	if IsClosable(iter) {
		return &closableTakeNIter[T]{
			takeNIter: tn,
		}
	}

	return &tn
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

type clonableTakeNIter[T any] struct {
	takeNIter[T]
}

func (c *clonableTakeNIter[T]) Clone() hie.Iter[T] {
	cu, cloned := Clone(c.under)
	if !cloned {
		panic("Clone called on an unclonable iterator")
	}
	return &clonableTakeNIter[T]{
		takeNIter: takeNIter[T]{
			under: cu,
			max:   c.max,
		},
	}
}

type closableTakeNIter[T any] struct {
	takeNIter[T]
}

func (c *closableTakeNIter[T]) Close() error {
	return Close(hie.Iter[T](&c.takeNIter))
}

type clonableClosableTakeNIter[T any] struct {
	clonableTakeNIter[T]
}

func (c *clonableClosableTakeNIter[T]) Close() error {
	return Close(hie.Iter[T](&c.clonableTakeNIter))
}
