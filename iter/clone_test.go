package iter

import (
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
