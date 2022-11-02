package hie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannelIter(t *testing.T) {
	ch := make(chan int, 10)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	it := Chan(ch)
	require.True(t, it.HasNext())
	require.Equal(t, 1, it.Next())
	require.True(t, it.HasNext())
	require.Equal(t, 2, it.Next())
	require.True(t, it.HasNext())
	require.Equal(t, 3, it.Next())
	require.False(t, it.HasNext())
}
