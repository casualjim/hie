package iter

import (
	"io"

	"github.com/casualjim/hie"
)

func IsClosable[T any](i hie.Iter[T]) bool {
	switch i.(type) {
	case io.Closer:
		return true
	case interface{ Close() }:
		return true
	}

	return false
}

func Close[T any](i hie.Iter[T]) error {
	switch it := i.(type) {
	case io.Closer:
		return it.Close()
	case interface{ Close() }:
		it.Close()
		return nil
	}

	return nil
}
