package iter

import (
	"testing"

	"github.com/casualjim/hie"
	"github.com/stretchr/testify/require"
)

type ioCloseableIter struct {
	count int
}

func (c *ioCloseableIter) HasNext() bool {
	return true
}

func (c *ioCloseableIter) Next() int {
	return 1
}

func (c *ioCloseableIter) Close() error {
	c.count++
	return nil
}

func newIOCloserIter() hie.Iter[int] {
	return &ioCloseableIter{}
}

type closableIter struct {
	count int
}

func (c *closableIter) HasNext() bool {
	return true
}

func (c *closableIter) Next() int {
	return 1
}

func (c *closableIter) Close() error {
	c.count++
	return nil
}

func newCloserIter() hie.Iter[int] {
	return &closableIter{}
}

func TestIterIOCloser(t *testing.T) {
	r := newIOCloserIter()
	require.NoError(t, Close(r))
	require.Equal(t, 1, r.(*ioCloseableIter).count)
}

func TestCloseIter(t *testing.T) {
	r := newCloserIter()
	require.NoError(t, Close(r))
	require.Equal(t, 1, r.(*closableIter).count)
}
