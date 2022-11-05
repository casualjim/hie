package iter

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/casualjim/hie"
	"github.com/stretchr/testify/require"
)

func TestClonableIter(t *testing.T) {
	i := makeConstantIter(3)

	ii, ok := Clone(i)
	require.True(t, ok)
	i.(*constantIter).val = 4

	require.Equal(t, 3, ii.(*constantIter).val)

	ic2, ok := Clone(ii)
	require.True(t, ok)
	ii.(*constantIter).val = 5
	require.Equal(t, 3, ic2.(*constantIter).val)
	require.Equal(t, 5, ii.(*constantIter).val)
}

func TestClonableNestedIter(t *testing.T) {
	t.Skip("skipping until clone can deal with nested clonable iter")
	total := &totalCount{}
	var s1 hie.Iter[int] = &countingCloneIter{w: hie.Slice(1, 2).AsIter(), total: total}
	var s2 hie.Iter[int] = &countingCloneIter{w: hie.Slice(3, 4).AsIter(), total: total}
	var s3 hie.Iter[int] = &countingCloneIter{w: hie.Slice(5, 6).AsIter(), total: total}
	var s hie.Iter[hie.Iter[int]] = &countingCloneIterIter{w: hie.Slice(s1, s2, s3).AsIter(), total: total}

	tt := reflect.TypeOf(s)
	fmt.Println(tt)
	_, ok := Clone(s)
	require.True(t, ok)
	require.Equal(t, 4, total.Total())
}

type constantIter struct {
	val int
}

func (c *constantIter) HasNext() bool { return true }
func (c *constantIter) Next() int     { return c.val }
func (c *constantIter) Clone() *constantIter {
	return &constantIter{val: c.val}
}
func makeConstantIter(val int) hie.Iter[int] {
	return &constantIter{val: val}
}
