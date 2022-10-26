package hie

import (
	"sync"
)

func Synchronize(locker ...sync.Locker) *Synchronizer {
	if len(locker) > 1 {
		panic("only 1 locker can be specified")
	}
	holder := &Synchronizer{}
	if len(locker) == 0 {
		holder.mu = new(sync.Mutex)
	}
	return holder
}

type Synchronizer struct {
	mu sync.Locker
}

func (l *Synchronizer) Do(thunk func() error) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return thunk()
}
