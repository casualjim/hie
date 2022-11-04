package hie

import (
	"testing"
	"time"

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

func TestChannelIterNil(t *testing.T) {

	it := Chan(chan int(nil))
	require.False(t, it.HasNext())
}

func TestChannelIterStartsEmpty(t *testing.T) {
	ch := make(chan int, 10)
	it := Chan(ch)
	go func() {
		<-time.After(100 * time.Millisecond)
		ch <- 1
	}()
	require.True(t, it.HasNext())
}
