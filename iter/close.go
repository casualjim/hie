package iter

import (
	"io"

	"github.com/casualjim/hie"
	"go.uber.org/multierr"
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

func CloseAll[T io.Closer](in hie.Iter[T]) error {
	return Fold(in, nil, func(err error, i T) (error, bool) {
		return multierr.Append(err, i.Close()), true
	})
}
