package hie

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcat(t *testing.T) {
	s1 := Slice(1, 2, 3)
	s2 := Slice(4, 5, 6)
	s3 := Slice(7, 8, 9)

	list := &concat[int]{}
	list.Append(s1)
	list.Append(s2)
	list.Append(s3)

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
	asIter := Slice(1, 2, 3)
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
	asIter := Slice(1, 2, 3)

	var result int
	ForEach(asIter, func(i int) bool {
		result += i
		return true
	})

	require.Equal(t, 6, result)
}

func TestSlice_Map(t *testing.T) {
	asIter := Slice(1, 2, 3)

	iter := Map(asIter, func(i int) string { return fmt.Sprintf("%d", i) })

	require.Equal(t, []string{"1", "2", "3"}, Collect(iter))
}

func TestSlice_FlatMap(t *testing.T) {
	asIter := Slice(Slice(1, 2), Slice(3, 4), Slice(5, 6))

	iter := FlatMap(asIter, func(i AsIter[int]) AsIter[string] {
		vals := Collect(Map(i, func(i int) string { return fmt.Sprintf("%d", i) }))

		return Slice(vals...)
	})

	require.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, Collect(iter))
}
