package hie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	t.Parallel()
	slice := Slice("a", "b", "c")

	require.True(t, Contains(slice, "a"))
	require.False(t, Contains(slice, "d"))
	require.False(t, Contains(slice, "e"))
}

func isEven(i int) bool         { return i%2 == 0 }
func isDivisibleBy4(i int) bool { return i%4 == 0 }
func isLt10(i int) bool         { return i < 10 }
func isLt3(i int) bool          { return i < 3 }
func isNaught(i int) bool       { return i == 0 }

func TestExists(t *testing.T) {
	t.Parallel()
	slice := Slice(1, 2, 3)

	require.True(t, Exists(slice, isEven))
	require.False(t, Exists(slice, isDivisibleBy4))
}

func TestAll(t *testing.T) {
	t.Parallel()

	slice := Slice(1, 2, 3)

	require.True(t, All(slice, isLt10))
	require.False(t, All(slice, isLt3))
}

func TestNonExist(t *testing.T) {
	t.Parallel()

	slice := Slice(1, 2, 3)

	require.True(t, NoneExist(slice, isNaught))
	require.False(t, NoneExist(slice, isLt3))
	require.False(t, NoneExist(slice, isLt10))
}

func TestIsSubset(t *testing.T) {
	t.Parallel()

	slice := Slice(1, 2, 3, 4, 5)
	require.True(t, IsSubset(slice, Slice(1, 2)))
	require.True(t, IsSubset(slice, Slice[int]()))
	require.False(t, IsSubset(slice, Slice(1, 6)))
	require.False(t, IsSubset(slice, Slice(0, 6)))
}

func TestIsDisjoint(t *testing.T) {
	t.Parallel()

	slice := Slice(1, 2, 3, 4, 5)

	require.False(t, IsDisjoint(slice, Slice(1, 2)))
	require.True(t, IsDisjoint(slice, Slice[int]()))
	require.False(t, IsDisjoint(slice, Slice(1, 6)))
	require.True(t, IsDisjoint(slice, Slice(0, 6)))
}

func TestIntersect(t *testing.T) {
	t.Parallel()

	slice := Slice(0, 1, 2, 3, 4, 5)
	zsSlice := Slice(0, 6)

	require.Equal(t, []int{0, 2}, Collect(Intersect(slice, Slice(0, 2))))
	require.Equal(t, []int{0}, Collect(Intersect(slice, zsSlice)))
	require.False(t, Intersect(slice, Slice(-1, 6)).AsIter().HasNext())
	require.Equal(t, []int{0}, Collect(Intersect(zsSlice, slice)))
	require.Equal(t, []int{0}, Collect(Intersect(Slice(0, 6, 0), slice)))

}

func TestDifference(t *testing.T) {
	t.Parallel()

	s1 := Slice(0, 1, 2, 3, 4, 5)
	s2 := Slice(0, 2, 6)
	s3 := Slice(1, 2, 3, 4, 5)
	s4 := Slice(0, 6)

	require.Equal(t, []int{1, 3, 4, 5}, Collect(Difference(s1, s2)))
	require.Equal(t, []int{6}, Collect(Difference(s2, s1)))
	require.Equal(t, []int{1, 2, 3, 4, 5}, Collect(Difference(s3, s4)))
	require.Equal(t, []int{0, 6}, Collect(Difference(s4, s3)))
	require.Empty(t, Collect(Difference(s1, s1)))
}

func TestSymmetricDifference(t *testing.T) {
	t.Parallel()

	s1 := Slice(0, 1, 2, 3, 4, 5)
	s2 := Slice(0, 2, 6)
	s3 := Slice(1, 2, 3, 4, 5)
	s4 := Slice(0, 6)

	require.Equal(t, []int{1, 3, 4, 5, 6}, Collect(SymmetricDifference(s1, s2)))
	require.Equal(t, []int{6, 1, 3, 4, 5}, Collect(SymmetricDifference(s2, s1)))
	require.Equal(t, []int{1, 2, 3, 4, 5, 0, 6}, Collect(SymmetricDifference(s3, s4)))
	require.Equal(t, []int{0, 6, 1, 2, 3, 4, 5}, Collect(SymmetricDifference(s4, s3)))
	require.Empty(t, Collect(SymmetricDifference(s1, s1)))
}

func TestUnion(t *testing.T) {
	t.Parallel()

	s1 := Slice(0, 1, 2, 3, 4, 5)

	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 10}, Collect(Union(s1, Slice(0, 2, 10))))
	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7}, Collect(Union(s1, Slice(6, 7))))
	require.Equal(t, []int{0, 1, 2, 3, 4, 5}, Collect(Union(s1, Slice[int]())))
	require.Equal(t, []int{0, 1, 2}, Collect(Union(Slice(0, 1, 2), Slice(0, 1, 2))))
	require.Empty(t, Collect(Union(Slice[int](), Slice[int]())))

	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 10, 11}, Collect(Union(s1, Slice(0, 2, 10), Slice(0, 1, 11))))
	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, Collect(Union(s1, Slice(6, 7), Slice(8, 9))))
	require.Equal(t, []int{0, 1, 2, 3, 4, 5}, Collect(Union(s1, Slice[int](), Slice[int]())))
	require.Equal(t, []int{0, 1, 2}, Collect(Union(Slice(0, 1, 2), Slice(0, 1, 2))))
	require.Empty(t, Collect(Union(Slice[int](), Slice[int](), Slice[int]())))
}
