package iter

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/casualjim/hie"
	"github.com/stretchr/testify/require"
)

type totalCount struct {
	v int
}

func (t *totalCount) Inc() {
	t.v++
}

func (t *totalCount) Total() int {
	return t.v
}

type countingCloneIter struct {
	w     hie.Iter[int]
	total *totalCount
}

func (n *countingCloneIter) Next() int {
	return n.w.Next()
}

func (n *countingCloneIter) HasNext() bool {
	return n.w.HasNext()
}

func (n *countingCloneIter) Clone() hie.Iter[int] {
	n.total.Inc()
	r, c := Clone(n.w)
	if !c {
		panic("clone called on an unclonable iterator")
	}
	return &countingCloneIter{
		w:     r,
		total: n.total,
	}
}

type countingCloseIter struct {
	w     hie.Iter[int]
	total *totalCount
}

func (n *countingCloseIter) Next() int {
	return n.w.Next()
}

func (n *countingCloseIter) HasNext() bool {
	return n.w.HasNext()
}

func (n *countingCloseIter) Close() error {
	n.total.Inc()
	return Close(n.w)
}

type countingCloneCloseIter struct {
	w      hie.Iter[int]
	clones *totalCount
	closes *totalCount
}

func (n *countingCloneCloseIter) Next() int {
	return n.w.Next()
}

func (n *countingCloneCloseIter) HasNext() bool {
	return n.w.HasNext()
}

func (n *countingCloneCloseIter) Clone() hie.Iter[int] {
	n.clones.Inc()
	return n
}

func (n *countingCloneCloseIter) Close() error {
	n.closes.Inc()
	return Close(n.w)
}

type countingCloneIterIter struct {
	w     hie.Iter[hie.Iter[int]]
	total *totalCount
}

func (n *countingCloneIterIter) Next() hie.Iter[int] {
	return n.w.Next()
}

func (n *countingCloneIterIter) HasNext() bool {
	return n.w.HasNext()
}

func (n *countingCloneIterIter) Clone() hie.Iter[hie.Iter[int]] {
	n.total.Inc()
	r, c := Clone(n.w)
	if !c {
		panic("clone called on an unclonable iterator")
	}
	return &countingCloneIterIter{
		w:     r,
		total: n.total,
	}
}

type countingCloseIterIter struct {
	w     hie.Iter[hie.Iter[int]]
	total *totalCount
}

func (n *countingCloseIterIter) Next() hie.Iter[int] {
	return n.w.Next()
}

func (n *countingCloseIterIter) HasNext() bool {
	return n.w.HasNext()
}

func (n *countingCloseIterIter) Close() error {
	n.total.Inc()
	return Close(n.w)
}

type countingCloneCloseIterIter struct {
	w      hie.Iter[hie.Iter[int]]
	clones *totalCount
	closes *totalCount
}

func (n *countingCloneCloseIterIter) Next() hie.Iter[int] {
	return n.w.Next()
}

func (n *countingCloneCloseIterIter) HasNext() bool {
	return n.w.HasNext()
}

func (n *countingCloneCloseIterIter) Clone() hie.Iter[hie.Iter[int]] {
	n.clones.Inc()
	return n
}

func (n *countingCloneCloseIterIter) Close() error {
	n.closes.Inc()
	return Close(n.w)
}

func TestConcatMixedUsePanics(t *testing.T) {
	total := &totalCount{}
	var s1 hie.Iter[int] = &countingCloseIter{w: hie.Slice(1, 2, 3).AsIter(), total: total}
	var s2 hie.Iter[int] = &countingCloneIter{w: hie.Slice(4, 5, 6).AsIter(), total: total}

	require.Panics(t, func() { Concat(s1, s2) })
}

func TestConcatClonableClosable(t *testing.T) {
	clones := &totalCount{}
	closes := &totalCount{}
	var s1 hie.Iter[int] = &countingCloneCloseIter{w: hie.Slice(1, 2, 3).AsIter(), clones: clones, closes: closes}
	var s2 hie.Iter[int] = &countingCloneCloseIter{w: hie.Slice(4, 5, 6).AsIter(), clones: clones, closes: closes}
	var s3 hie.Iter[int] = &countingCloneCloseIter{w: hie.Slice(7, 8, 9).AsIter(), clones: clones, closes: closes}
	var s4 hie.Iter[int] = &countingCloneCloseIter{w: hie.Slice(10, 11, 12).AsIter(), clones: clones, closes: closes}

	list := Concat(s1, s2, s3, s4)

	nlist, cloned := Clone(list)
	require.True(t, cloned)
	require.Equal(t, 4, clones.Total())

	require.True(t, list.HasNext())
	require.NoError(t, Close(list))
	require.Equal(t, 4, closes.Total())

	require.False(t, list.HasNext())
	require.Panics(t, func() { list.Next() })

	require.True(t, nlist.HasNext())
	require.Equal(t, 1, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 2, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 3, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 4, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 5, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 6, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 7, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 8, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 9, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 10, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 11, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 12, nlist.Next())
	require.Panics(t, func() { nlist.Next() })
}

func TestConcatClosable(t *testing.T) {
	total := &totalCount{}
	var s1 hie.Iter[int] = &countingCloseIter{w: hie.Slice(1, 2, 3).AsIter(), total: total}
	var s2 hie.Iter[int] = &countingCloseIter{w: hie.Slice(4, 5, 6).AsIter(), total: total}
	var s3 hie.Iter[int] = &countingCloseIter{w: hie.Slice(7, 8, 9).AsIter(), total: total}
	var s4 hie.Iter[int] = &countingCloseIter{w: hie.Slice(10, 11, 12).AsIter(), total: total}

	list := Concat(s1, s2, s3, s4)

	require.True(t, list.HasNext())
	require.NoError(t, Close(list))
	require.Equal(t, 4, total.Total())

	require.False(t, list.HasNext())
	require.Panics(t, func() { list.Next() })
}

func TestConcatClonable(t *testing.T) {
	total := &totalCount{}
	var s1 hie.Iter[int] = &countingCloneIter{w: hie.Slice(1, 2, 3).AsIter(), total: total}
	var s2 hie.Iter[int] = &countingCloneIter{w: hie.Slice(4, 5, 6).AsIter(), total: total}
	var s3 hie.Iter[int] = &countingCloneIter{w: hie.Slice(7, 8, 9).AsIter(), total: total}
	var s4 hie.Iter[int] = &countingCloneIter{w: hie.Slice(10, 11, 12).AsIter(), total: total}

	list := Concat(s1, s2, s3, s4)

	nlist, cloned := Clone(list)
	require.True(t, cloned)
	require.Equal(t, 4, total.Total())

	require.True(t, nlist.HasNext())
	require.Equal(t, 1, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 2, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 3, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 4, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 5, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 6, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 7, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 8, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 9, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 10, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 11, nlist.Next())
	require.True(t, nlist.HasNext())
	require.Equal(t, 12, nlist.Next())
	require.Panics(t, func() { nlist.Next() })
}

func TestConcat(t *testing.T) {
	t.Parallel()

	s1 := hie.Slice(1, 2, 3)
	s2 := hie.Slice(4, 5, 6)
	s3 := hie.Slice(7, 8, 9)

	list := Concat(s1.AsIter(), s2.AsIter(), s3.AsIter())

	require.True(t, list.HasNext())
	require.Equal(t, 1, list.Next())
	require.True(t, list.HasNext())
	require.Equal(t, 2, list.Next())
	require.True(t, list.HasNext())
	require.Equal(t, 3, list.Next())
	require.True(t, list.HasNext())
	require.Equal(t, 4, list.Next())
	require.True(t, list.HasNext())
	require.Equal(t, 5, list.Next())
	require.True(t, list.HasNext())
	require.Equal(t, 6, list.Next())
	require.True(t, list.HasNext())
	require.Equal(t, 7, list.Next())
	require.True(t, list.HasNext())
	require.Equal(t, 8, list.Next())
	require.True(t, list.HasNext())
	require.Equal(t, 9, list.Next())
	require.Panics(t, func() { list.Next() })
}

func TestSlice_Iter(t *testing.T) {
	t.Parallel()

	asIter := hie.Slice(1, 2, 3)
	iter := asIter.AsIter()

	require.True(t, iter.HasNext())
	require.Equal(t, 1, iter.Next())

	require.True(t, iter.HasNext())
	require.Equal(t, 2, iter.Next())

	require.True(t, iter.HasNext())
	require.Equal(t, 3, iter.Next())

	require.False(t, iter.HasNext())
	require.Panics(t, func() { iter.Next() })
}

func TestSlice_ForEach(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3)

	var result int
	ForEach(slice.AsIter(), func(i int) bool {
		result += i
		return true
	})

	require.Equal(t, 6, result)
}

func TestSlice_Map(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3)

	iter := Map(slice.AsIter(), func(i int) string { return fmt.Sprintf("%d", i) })

	require.Equal(t, []string{"1", "2", "3"}, Collect(iter))
}

func TestMapClonable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(1, 2, 3)
	total := &totalCount{}
	var slice hie.Iter[int] = &countingCloneIter{w: sl.AsIter(), total: total}

	result := Map(slice, func(i int) string { return fmt.Sprintf("%d", i) })
	cres, cloned := Clone(result)
	require.True(t, cloned)
	require.Equal(t, 1, total.Total())

	require.Equal(t, []string{"1", "2", "3"}, Collect(result))
	require.Equal(t, []string{"1", "2", "3"}, Collect(cres))
}

func TestMapClosable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(1, 2, 3)
	total := &totalCount{}
	var slice hie.Iter[int] = &countingCloseIter{w: sl.AsIter(), total: total}

	result := Map(slice, func(i int) string { return fmt.Sprintf("%d", i) })
	require.NoError(t, Close(result))
	require.Equal(t, 1, total.Total())

	require.Equal(t, []string(nil), Collect(result))
}

func TestMapClonableClosable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(1, 2, 3)
	clones := &totalCount{}
	closes := &totalCount{}
	var slice hie.Iter[int] = &countingCloneCloseIter{w: sl.AsIter(), clones: clones, closes: closes}

	result := Map(slice, func(i int) string { return fmt.Sprintf("%d", i) })
	cres, cloned := Clone(result)
	require.True(t, cloned)
	require.Equal(t, 1, clones.Total())

	require.NoError(t, Close(result))
	require.Equal(t, 1, closes.Total())

	require.Equal(t, []string(nil), Collect(result))
	require.Equal(t, []string{"1", "2", "3"}, Collect(cres))
}

func TestSlice_FlatMap(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(hie.Slice(1, 2).AsIter(), hie.Slice(3, 4).AsIter(), hie.Slice(5, 6).AsIter())

	iter := FlatMap(slice.AsIter(), func(i hie.Iter[int]) hie.Iter[string] {
		return Map(i, func(i int) string { return fmt.Sprintf("%d", i) })
	})

	require.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, Collect(iter))
}

func TestSlice_FlatMapClonable(t *testing.T) {
	t.Skip("skipping until clone can deal with nested clonable iter")
	t.Parallel()

	sl := hie.Slice(hie.Slice(1, 2).AsIter(), hie.Slice(3, 4).AsIter(), hie.Slice(5, 6).AsIter())
	total := &totalCount{}
	var slice hie.Iter[hie.Iter[int]] = &countingCloneIterIter{w: sl.AsIter(), total: total}

	iter := FlatMap(slice, func(i hie.Iter[int]) hie.Iter[string] {
		return Map(i, func(i int) string { return fmt.Sprintf("%d", i) })
	})

	cres, cloned := Clone(iter)
	require.True(t, cloned)
	require.Equal(t, 1, total.Total())

	require.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, Collect(iter))
	require.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, Collect(cres))
}

func TestSlice_FlatMapClosable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(hie.Slice(1, 2).AsIter(), hie.Slice(3, 4).AsIter(), hie.Slice(5, 6).AsIter())
	total := &totalCount{}
	var slice hie.Iter[hie.Iter[int]] = &countingCloseIterIter{w: sl.AsIter(), total: total}

	iter := FlatMap(slice, func(i hie.Iter[int]) hie.Iter[string] {
		return Map(i, func(i int) string { return fmt.Sprintf("%d", i) })
	})

	require.NoError(t, Close(iter))
	require.Equal(t, 1, total.Total())

	require.Equal(t, []string(nil), Collect(iter))
}

func TestSlice_FlatMapClonableClosable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(hie.Slice(1, 2).AsIter(), hie.Slice(3, 4).AsIter(), hie.Slice(5, 6).AsIter())
	clones := &totalCount{}
	closes := &totalCount{}
	var slice hie.Iter[hie.Iter[int]] = &countingCloneCloseIterIter{w: sl.AsIter(), clones: clones, closes: closes}

	iter := FlatMap(slice, func(i hie.Iter[int]) hie.Iter[string] {
		return Map(i, func(i int) string { return fmt.Sprintf("%d", i) })
	})

	cres, cloned := Clone(iter)
	require.True(t, cloned)
	require.Equal(t, 1, clones.Total())

	require.NoError(t, Close(iter))
	require.Equal(t, 1, closes.Total())

	require.Equal(t, []string(nil), Collect(iter))
	require.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, Collect(cres))
}

func TestSlice_Flatten(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(hie.Slice(1, 2).AsIter(), hie.Slice(3, 4).AsIter(), hie.Slice(5, 6).AsIter())

	iter := Flatten(slice.AsIter())

	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, Collect(iter))
}

func TestFilter(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)

	result := Filter(slice.AsIter(), func(i int) bool { return i%2 == 0 })

	require.Equal(t, []int{2, 4, 6, 8}, Collect(result))
}

func TestFilterClonable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	total := &totalCount{}
	var slice hie.Iter[int] = &countingCloneIter{w: sl.AsIter(), total: total}

	result := Filter(slice, func(i int) bool { return i%2 == 0 })
	cres, cloned := Clone(result)
	require.True(t, cloned)
	require.Equal(t, 1, total.Total())

	require.Equal(t, []int{2, 4, 6, 8}, Collect(result))
	require.Equal(t, []int{2, 4, 6, 8}, Collect(cres))
}

func TestFilterClosable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	total := &totalCount{}
	var slice hie.Iter[int] = &countingCloseIter{w: sl.AsIter(), total: total}

	result := Filter(slice, func(i int) bool { return i%2 == 0 })
	require.NoError(t, Close(result))
	require.Equal(t, 1, total.Total())

	require.Equal(t, []int(nil), Collect(result))
}

func TestFilterClonableClosable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	clones := &totalCount{}
	closes := &totalCount{}
	var slice hie.Iter[int] = &countingCloneCloseIter{w: sl.AsIter(), clones: clones, closes: closes}

	result := Filter(slice, func(i int) bool { return i%2 == 0 })
	cres, cloned := Clone(result)
	require.True(t, cloned)
	require.Equal(t, 1, clones.Total())

	require.NoError(t, Close(result))
	require.Equal(t, 1, closes.Total())

	require.Equal(t, []int(nil), Collect(result))
	require.Equal(t, []int{2, 4, 6, 8}, Collect(cres))
}

func TestFilterMap(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)

	result := FilterMap(slice.AsIter(), func(i int) (string, bool) { return strconv.Itoa(i), i%2 == 0 })

	require.Equal(t, []string{"2", "4", "6", "8"}, Collect(result))
}

func TestFilterMapClonable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	total := &totalCount{}
	var slice hie.Iter[int] = &countingCloneIter{w: sl.AsIter(), total: total}

	result := FilterMap(slice, func(i int) (string, bool) { return strconv.Itoa(i), i%2 == 0 })
	cres, cloned := Clone(result)
	require.True(t, cloned)
	require.Equal(t, 1, total.Total())

	require.Equal(t, []string{"2", "4", "6", "8"}, Collect(result))
	require.Equal(t, []string{"2", "4", "6", "8"}, Collect(cres))
}

func TestFilterMapClosable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	total := &totalCount{}
	var slice hie.Iter[int] = &countingCloseIter{w: sl.AsIter(), total: total}

	result := FilterMap(slice, func(i int) (string, bool) { return strconv.Itoa(i), i%2 == 0 })
	require.NoError(t, Close(result))
	require.Equal(t, 1, total.Total())

	require.Equal(t, []string(nil), Collect(result))
}

func TestFilterMapClonableClosable(t *testing.T) {
	t.Parallel()

	sl := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	clones := &totalCount{}
	closes := &totalCount{}
	var slice hie.Iter[int] = &countingCloneCloseIter{w: sl.AsIter(), clones: clones, closes: closes}

	result := FilterMap(slice, func(i int) (string, bool) { return strconv.Itoa(i), i%2 == 0 })
	cres, cloned := Clone(result)
	require.True(t, cloned)
	require.Equal(t, 1, clones.Total())

	require.NoError(t, Close(result))
	require.Equal(t, 1, closes.Total())

	require.Equal(t, []string(nil), Collect(result))
	require.Equal(t, []string{"2", "4", "6", "8"}, Collect(cres))
}

func TestTakeN(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	result := TakeN(slice.AsIter(), 4)

	empty := Empty[int]()
	er := TakeN(empty, 4)

	empty2 := hie.Slice[int]()
	er2 := TakeN(empty2.AsIter(), 3)

	require.Equal(t, []int{1, 2, 3, 4}, Collect(result))
	require.Equal(t, []int(nil), Collect(er))
	require.Equal(t, []int(nil), Collect(er2))
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, Collect(TakeN(slice.AsIter(), 10)))
}

func TestTakeNClonable(t *testing.T) {
	t.Parallel()

	total := &totalCount{}
	sl := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	var slice hie.Iter[int] = &countingCloneIter{w: sl.AsIter(), total: total}
	result := TakeN(slice, 4)
	cres, cloned := Clone(result)
	require.True(t, cloned)
	require.Equal(t, 1, total.Total())

	require.Equal(t, []int{1, 2, 3, 4}, Collect(result))
	require.Equal(t, []int{1, 2, 3, 4}, Collect(cres))
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, Collect(TakeN(sl.AsIter(), 10)))
}

func TestTakeNClosable(t *testing.T) {
	t.Parallel()

	total := &totalCount{}
	sl := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	var slice hie.Iter[int] = &countingCloseIter{w: sl.AsIter(), total: total}

	result := TakeN(slice, 4)
	require.NoError(t, Close(result))
	require.Equal(t, 1, total.Total())

	require.Equal(t, []int(nil), Collect(result))
}

func TestTakeNClonableClosable(t *testing.T) {
	t.Parallel()

	clones := &totalCount{}
	closes := &totalCount{}
	sl := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	var slice hie.Iter[int] = &countingCloneCloseIter{w: sl.AsIter(), clones: clones, closes: closes}
	result := TakeN(slice, 4)

	cres, cloned := Clone(result)
	require.True(t, cloned)
	require.Equal(t, 1, clones.Total())

	require.NoError(t, Close(result))
	require.Equal(t, 1, closes.Total())

	require.Equal(t, []int(nil), Collect(result))
	require.Equal(t, []int{1, 2, 3, 4}, Collect(cres))
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, Collect(TakeN(sl.AsIter(), 10)))
}
