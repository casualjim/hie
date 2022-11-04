package iter

import (
	"io"

	"github.com/casualjim/hie"
)

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
