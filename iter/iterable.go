package iter

import (
	"github.com/casualjim/hie"
	"github.com/casualjim/hie/opt"
	"go.uber.org/multierr"
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
	closed bool
}

func (c *closableFilterMapperIter[T, R]) HasNext() bool {
	return !c.closed && c.filterMapperIter.HasNext()
}

func (c *closableFilterMapperIter[T, R]) Next() R {
	if c.closed {
		panic("next called on a closed iterator")
	}
	return c.filterMapperIter.Next()
}

func (c *closableFilterMapperIter[T, R]) Close() error {
	c.closed = true
	return Close(c.under)
}

type clonableClosableFilterMapperIter[T, R any] struct {
	clonableFilterMapper[T, R]
	closed bool
}

func (c *clonableClosableFilterMapperIter[T, R]) HasNext() bool {
	return !c.closed && c.clonableFilterMapper.HasNext()
}

func (c *clonableClosableFilterMapperIter[T, R]) Next() R {
	if c.closed {
		panic("next called on a closed iterator")
	}
	return c.clonableFilterMapper.Next()
}
func (c *clonableClosableFilterMapperIter[T, R]) Close() error {
	c.closed = true
	return Close(c.under)
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
			under:     cu,
			mapperFn:  c.mapperFn,
			lastMatch: opt.None[R](),
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
	closed bool
}

func (c *closableMapperIter[T, R]) HasNext() bool {
	return !c.closed && c.mapperIter.HasNext()
}

func (c *closableMapperIter[T, R]) Next() R {
	if c.closed {
		panic("next called on a closed iterator")
	}
	return c.mapperIter.Next()
}

func (c *closableMapperIter[T, R]) Close() error {
	c.closed = true
	return Close(c.under)
}

type clonableClosableMapperIter[T, R any] struct {
	clonableMapper[T, R]
	closed bool
}

func (c *clonableClosableMapperIter[T, R]) HasNext() bool {
	return !c.closed && c.clonableMapper.HasNext()
}

func (c *clonableClosableMapperIter[T, R]) Next() R {
	if c.closed {
		panic("next called on a closed iterator")
	}
	return c.clonableMapper.Next()
}

func (c *clonableClosableMapperIter[T, R]) Close() error {
	c.closed = true
	return Close(c.under)
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
	closed bool
}

func (c *closableFlatMapperIter[T, R]) HasNext() bool {
	return !c.closed && c.flatMapperIter.HasNext()
}

func (c *closableFlatMapperIter[T, R]) Next() R {
	if !c.closed {
		panic("next called on a closed iterator")
	}
	return c.flatMapperIter.Next()
}

func (c *closableFlatMapperIter[T, R]) Close() error {
	c.closed = true
	return Close(c.under)
}

type clonableClosableFlatMapperIter[T, R any] struct {
	clonableFlatMapper[T, R]
	closed bool
}

func (c *clonableClosableFlatMapperIter[T, R]) HasNext() bool {
	return !c.closed && c.clonableFlatMapper.HasNext()
}

func (c *clonableClosableFlatMapperIter[T, R]) Next() R {
	if !c.closed {
		panic("next called on a closed iterator")
	}
	return c.clonableFlatMapper.Next()
}

func (c *clonableClosableFlatMapperIter[T, R]) Close() error {
	c.closed = true
	return Close(c.under)
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
			lastMatch: opt.None[T](),
		},
	}
}

type closableFilterIter[T any] struct {
	filterIter[T]
	closed bool
}

func (c *closableFilterIter[T]) HasNext() bool {
	return !c.closed && c.filterIter.HasNext()
}

func (c *closableFilterIter[T]) Next() T {
	if !c.closed {
		panic("next called on a closed iterator")
	}
	return c.filterIter.Next()
}

func (c *closableFilterIter[T]) Close() error {
	c.closed = true
	return Close(c.under)
}

type clonableClosableFilterIter[T any] struct {
	clonableFilterIter[T]
	closed bool
}

func (c *clonableClosableFilterIter[T]) HasNext() bool {
	return !c.closed && c.clonableFilterIter.HasNext()
}

func (c *clonableClosableFilterIter[T]) Next() T {
	if !c.closed {
		panic("next called on a closed iterator")
	}
	return c.clonableFilterIter.Next()
}

func (c *clonableClosableFilterIter[T]) Close() error {
	c.closed = true
	return Close(c.under)
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
	switch conc := left.(type) {
	case *concat[T]:
		conc.Append(right)
		for _, other := range others {
			conc.Append(other)
		}
		return conc
	case *closableConcatIter[T]:
		conc.Append(right)
		for _, other := range others {
			conc.Append(other)
		}
		return conc
	case *clonableConcatIter[T]:
		conc.Append(right)
		for _, other := range others {
			conc.Append(other)
		}
		return conc
	case *clonableClosableConcatIter[T]:
		conc.Append(right)
		for _, other := range others {
			conc.Append(other)
		}
		return conc
	default:
		isClonable, isClosable := IsClonable(left), IsClosable(left)
		switch {
		case isClonable && isClosable:
			cc := &clonableClosableConcatIter[T]{}
			cc.Append(left)
			cc.Append(right)
			for _, other := range others {
				cc.Append(other)
			}
			return cc
		case isClonable:
			cc := &clonableConcatIter[T]{}
			cc.Append(left)
			cc.Append(right)
			for _, other := range others {
				cc.Append(other)
			}
			return cc
		case isClosable:
			cc := &closableConcatIter[T]{}
			cc.Append(left)
			cc.Append(right)
			for _, other := range others {
				cc.Append(other)
			}
			return cc
		default:
			cc := &concat[T]{}
			cc.Append(left)
			cc.Append(right)
			for _, other := range others {
				cc.Append(other)
			}
			return cc
		}
	}

}

type clonableConcatIter[T any] struct {
	head *clonableConsIter[T]
	tail *clonableConsIter[T]
}

func (c *clonableConcatIter[T]) Append(i hie.Iter[T]) {
	isClonable, isClosable := IsClonable(i), IsClosable(i)
	if !isClonable && isClosable {
		panic("Adding an unclonable, closable iterator to a collection of clonable, unclosable iterators")
	} else if !isClonable {
		panic("Adding an unclonable iterator to a collection of clonable, unclosable iterators")
	} else if isClosable {
		panic("Adding a closable iterator to a collection of clonable, unclosable iterators")
	}

	newItem := &clonableConsIter[T]{
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

func (c *clonableConcatIter[T]) Clone() hie.Iter[T] {
	var head = c.head
	var cur = head
	var prev *clonableConsIter[T]
	var tail *clonableConsIter[T]

	for {
		if cur.next == nil {
			tail = cur.Clone().(*clonableConsIter[T])
			prev.next = tail
			break
		}

		if prev == nil {
			nxt := cur.next
			head = cur.Clone().(*clonableConsIter[T])
			prev = head
			cur = nxt
			continue
		}

		nxt := cur.next
		c := cur.Clone().(*clonableConsIter[T])
		cur = nxt
		prev.next = c
		prev = c

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
	head   *closableConsIter[T]
	tail   *closableConsIter[T]
	closed bool
}

func (c *closableConcatIter[T]) HasNext() bool {
	return !c.closed && c.head != nil && c.head.HasNext()
}

func (c *closableConcatIter[T]) Next() T {
	if c == nil {
		panic("iterating an empty clonable concat iterator")
	}
	if c.closed {
		panic("next called on a closed iterator")
	}
	if c.head == nil {
		panic("iterating an empty clonable concat iterator")
	}
	return c.head.Next()
}

func (c *closableConcatIter[T]) Close() error {
	c.closed = true
	head := c.head
	cur := head
	var err error
	for {
		if cur == nil {
			break
		}
		err = multierr.Append(err, cur.Close())
		cur = cur.next
	}
	return err
}

func (c *closableConcatIter[T]) Append(i hie.Iter[T]) {
	isClonable, isClosable := IsClonable(i), IsClosable(i)
	if isClonable && !isClosable {
		panic("Adding a clonable, unclosable iterator to a collection of closable iterators")
	} else if isClonable {
		panic("Adding a clonable iterator to a collection of unclonable, closable iterators")
	} else if !isClosable {
		panic("Adding an unclosable iterator to a collection of unclonable, closable iterators")
	}

	newItem := &closableConsIter[T]{
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

type clonableClosableConcatIter[T any] struct {
	head   *clonableClosableConsIter[T]
	tail   *clonableClosableConsIter[T]
	closed bool
}

func (c *clonableClosableConcatIter[T]) Clone() hie.Iter[T] {
	var head = c.head
	var cur = head
	var prev *clonableClosableConsIter[T]
	var tail *clonableClosableConsIter[T]

	for {
		if cur.next == nil {
			tail = cur.Clone().(*clonableClosableConsIter[T])
			prev.next = tail
			break
		}

		if prev == nil {
			nxt := cur.next
			head = cur.Clone().(*clonableClosableConsIter[T])
			prev = head
			cur = nxt
			continue
		}

		nxt := cur.next
		c := cur.Clone().(*clonableClosableConsIter[T])
		cur = nxt
		prev.next = c
		prev = c

	}

	return &clonableClosableConcatIter[T]{
		head:   head,
		tail:   tail,
		closed: c.closed,
	}
}

func (c *clonableClosableConcatIter[T]) Append(i hie.Iter[T]) {
	isClonable, isClosable := IsClonable(i), IsClosable(i)
	if !isClonable && !isClosable {
		panic("Adding an unclonable, unclosable iterator to a collection of clonable, closable iterators")
	} else if !isClonable {
		panic("Adding an unclonable iterator to a collection of clonable, closable iterators")
	} else if !isClosable {
		panic("Adding an unclosable iterator to a collection of clonable, closable iterators")
	}

	newItem := &clonableClosableConsIter[T]{
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

func (c *clonableClosableConcatIter[T]) HasNext() bool {
	return !c.closed && c.head != nil && c.head.HasNext()
}

func (c *clonableClosableConcatIter[T]) Next() T {
	if c == nil {
		panic("iterating an empty clonable concat iterator")
	}
	if c.closed {
		panic("next called on a closed iterator")
	}
	if c.head == nil {
		panic("iterating an empty clonable concat iterator")
	}
	return c.head.Next()
}

func (c *clonableClosableConcatIter[T]) Close() error {
	c.closed = true
	head := c.head
	cur := head
	var err error
	for {
		if cur == nil {
			break
		}
		err = multierr.Append(err, cur.Close())
		cur = cur.next
	}
	return err
}

type concat[T any] struct {
	head *cons[T]
	tail *cons[T]
}

func (c *concat[T]) Append(i hie.Iter[T]) {
	isClonable, isClosable := IsClonable(i), IsClosable(i)
	if isClonable && isClosable {
		panic("Adding a clonable, closable iterator to a collection of unclosable and unclonable iterators")
	} else if isClonable {
		panic("Adding a clonable iterator to a collection of unclonable iterators")
	} else if isClosable {
		panic("Adding a closable iterator to a collection of unclosable iterators")
	}

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
		next:  c.next, // clone for the next element is handled in the container
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
	return Close(c.under)
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
		next:  c.next, // clone for the next element is handled in the container
	}
}

func (c *clonableClosableConsIter[T]) Close() error {
	return Close(c.under)
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
	closed bool
}

func (c *closableTakeNIter[T]) HasNext() bool {
	return !c.closed && c.under.HasNext()
}

func (c *closableTakeNIter[T]) Next() T {
	if c.closed {
		panic("next called on a closed iterator")
	}
	return c.under.Next()
}

func (c *closableTakeNIter[T]) Close() error {
	c.closed = true
	return Close(hie.Iter[T](c.under))
}

type clonableClosableTakeNIter[T any] struct {
	clonableTakeNIter[T]
	closed bool
}

func (c *clonableClosableTakeNIter[T]) HasNext() bool {
	return !c.closed && c.under.HasNext()
}

func (c *clonableClosableTakeNIter[T]) Next() T {
	if c.closed {
		panic("next called on a closed iterator")
	}
	return c.under.Next()
}

func (c *clonableClosableTakeNIter[T]) Close() error {
	c.closed = true
	return Close(hie.Iter[T](c.under))
}
