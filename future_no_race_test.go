//go:build !race

package hie_test

import (
	"sync"
	"testing"

	"github.com/casualjim/hie"
)

func TestCancelConc(t *testing.T) {
	loop := func() {
		const N = 8000
		start := make(chan int)
		var done sync.WaitGroup
		done.Add(N)
		f := hie.Do(hie.Func(func() (int, error) {
			select {} //block
			return 1, nil
		}))
		for i := 0; i < N; i++ {
			go func() {
				defer done.Done()
				<-start
				f.Cancel()
			}()
		}
		close(start)
		done.Wait()
	}

	for i := 0; i < 500; i++ {
		loop()
	}

}
