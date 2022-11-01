package iter

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/casualjim/hie"
	"github.com/stretchr/testify/require"
)

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

func TestSlice_FlatMap(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(hie.Slice(1, 2).AsIter(), hie.Slice(3, 4).AsIter(), hie.Slice(5, 6).AsIter())

	iter := FlatMap(slice.AsIter(), func(i hie.Iter[int]) hie.Iter[string] {
		return Map(i, func(i int) string { return fmt.Sprintf("%d", i) })
	})

	require.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, Collect(iter))
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

func TestFilterMap(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)

	result := FilterMap(slice.AsIter(), func(i int) (string, bool) { return strconv.Itoa(i), i%2 == 0 })

	require.Equal(t, []string{"2", "4", "6", "8"}, Collect(result))
}

func TestTakeN(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3, 4, 5, 6, 7, 8)
	result := TakeN(slice.AsIter(), 4)

	empty := hie.EmptyIter[int]()
	er := TakeN(empty, 4)

	empty2 := hie.Slice[int]()
	er2 := TakeN(empty2.AsIter(), 3)

	require.Equal(t, []int{1, 2, 3, 4}, Collect(result))
	require.Equal(t, []int(nil), Collect(er))
	require.Equal(t, []int(nil), Collect(er2))
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, Collect(TakeN(slice.AsIter(), 10)))
}
