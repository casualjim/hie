package iter

import (
	"testing"

	"github.com/casualjim/hie"
	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	t.Parallel()
	slice := hie.Slice("a", "b", "c")

	require.True(t, Contains(slice.AsIter(), "a"))
	require.False(t, Contains(slice.AsIter(), "d"))
	require.False(t, Contains(slice.AsIter(), "e"))
}

func isEven(i int) bool         { return i%2 == 0 }
func isDivisibleBy4(i int) bool { return i%4 == 0 }
func isLt10(i int) bool         { return i < 10 }
func isLt3(i int) bool          { return i < 3 }
func isNaught(i int) bool       { return i == 0 }

func TestExists(t *testing.T) {
	t.Parallel()
	slice := hie.Slice(1, 2, 3)

	require.True(t, Exists(slice.AsIter(), isEven))
	require.False(t, Exists(slice.AsIter(), isDivisibleBy4))
}

func TestAll(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3)

	require.True(t, All(slice.AsIter(), isLt10))
	require.False(t, All(slice.AsIter(), isLt3))
}

func TestNonExist(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3)

	require.True(t, NoneExist(slice.AsIter(), isNaught))
	require.False(t, NoneExist(slice.AsIter(), isLt3))
	require.False(t, NoneExist(slice.AsIter(), isLt10))
}

func TestIsSubset(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3, 4, 5)
	require.True(t, IsSubset(slice.AsIter(), hie.Slice(1, 2).AsIter()))
	require.True(t, IsSubset(slice.AsIter(), hie.Slice[int]().AsIter()))
	require.False(t, IsSubset(slice.AsIter(), hie.Slice(1, 6).AsIter()))
	require.False(t, IsSubset(slice.AsIter(), hie.Slice(0, 6).AsIter()))
}

func TestIsDisjoint(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(1, 2, 3, 4, 5)

	require.False(t, IsDisjoint(slice.AsIter(), hie.Slice(1, 2).AsIter()))
	require.True(t, IsDisjoint(slice.AsIter(), hie.Slice[int]().AsIter()))
	require.False(t, IsDisjoint(slice.AsIter(), hie.Slice(1, 6).AsIter()))
	require.True(t, IsDisjoint(slice.AsIter(), hie.Slice(0, 6).AsIter()))
}

func TestIntersect(t *testing.T) {
	t.Parallel()

	slice := hie.Slice(0, 1, 2, 3, 4, 5)
	zsSlice := hie.Slice(0, 6)

	require.Equal(t, []int{0, 2}, Collect(Intersect(slice.AsIter(), hie.Slice(0, 2).AsIter())))
	require.Equal(t, []int{0}, Collect(Intersect(slice.AsIter(), zsSlice.AsIter())))
	require.False(t, Intersect(slice.AsIter(), hie.Slice(-1, 6).AsIter()).HasNext())
	require.Equal(t, []int{0}, Collect(Intersect(zsSlice.AsIter(), slice.AsIter())))
	require.Equal(t, []int{0}, Collect(Intersect(hie.Slice(0, 6, 0).AsIter(), slice.AsIter())))

}

func TestDifference(t *testing.T) {
	t.Parallel()

	s1 := hie.Slice(0, 1, 2, 3, 4, 5)
	s2 := hie.Slice(0, 2, 6)
	s3 := hie.Slice(1, 2, 3, 4, 5)
	s4 := hie.Slice(0, 6)

	require.Equal(t, []int{1, 3, 4, 5}, Collect(Difference(s1.AsIter(), s2.AsIter())))
	require.Equal(t, []int{6}, Collect(Difference(s2.AsIter(), s1.AsIter())))
	require.Equal(t, []int{1, 2, 3, 4, 5}, Collect(Difference(s3.AsIter(), s4.AsIter())))
	require.Equal(t, []int{0, 6}, Collect(Difference(s4.AsIter(), s3.AsIter())))
	require.Empty(t, Collect(Difference(s1.AsIter(), s1.AsIter())))
}

func TestSymmetricDifference(t *testing.T) {
	t.Parallel()

	s1 := hie.Slice(0, 1, 2, 3, 4, 5)
	s2 := hie.Slice(0, 2, 6)
	s3 := hie.Slice(1, 2, 3, 4, 5)
	s4 := hie.Slice(0, 6)

	require.Equal(t, []int{1, 3, 4, 5, 6}, Collect(SymmetricDifference(s1.AsIter(), s2.AsIter())))
	require.Equal(t, []int{6, 1, 3, 4, 5}, Collect(SymmetricDifference(s2.AsIter(), s1.AsIter())))
	require.Equal(t, []int{1, 2, 3, 4, 5, 0, 6}, Collect(SymmetricDifference(s3.AsIter(), s4.AsIter())))
	require.Equal(t, []int{0, 6, 1, 2, 3, 4, 5}, Collect(SymmetricDifference(s4.AsIter(), s3.AsIter())))
	require.Empty(t, Collect(SymmetricDifference(s1.AsIter(), s1.AsIter())))
}

func TestUnion(t *testing.T) {
	t.Parallel()

	s1 := hie.Slice(0, 1, 2, 3, 4, 5)

	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 10}, Collect(Union(s1.AsIter(), hie.Slice(0, 2, 10).AsIter())))
	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7}, Collect(Union(s1.AsIter(), hie.Slice(6, 7).AsIter())))
	require.Equal(t, []int{0, 1, 2, 3, 4, 5}, Collect(Union(s1.AsIter(), hie.Slice[int]().AsIter())))
	require.Equal(t, []int{0, 1, 2}, Collect(Union(hie.Slice(0, 1, 2).AsIter(), hie.Slice(0, 1, 2).AsIter())))
	require.Empty(t, Collect(Union(hie.Slice[int]().AsIter(), hie.Slice[int]().AsIter())))

	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 10, 11}, Collect(Union(s1.AsIter(), hie.Slice(0, 2, 10).AsIter(), hie.Slice(0, 1, 11).AsIter())))
	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, Collect(Union(s1.AsIter(), hie.Slice(6, 7).AsIter(), hie.Slice(8, 9).AsIter())))
	require.Equal(t, []int{0, 1, 2, 3, 4, 5}, Collect(Union(s1.AsIter(), hie.Slice[int]().AsIter(), hie.Slice[int]().AsIter())))
	require.Equal(t, []int{0, 1, 2}, Collect(Union(hie.Slice(0, 1, 2).AsIter(), hie.Slice(0, 1, 2).AsIter())))
	require.Empty(t, Collect(Union(hie.Slice[int]().AsIter(), hie.Slice[int]().AsIter(), hie.Slice[int]().AsIter())))
}
