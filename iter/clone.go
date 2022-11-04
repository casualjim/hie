package iter

import (
	"reflect"

	"github.com/casualjim/hie"
)

func Clone[T any](i hie.Iter[T]) (hie.Iter[T], bool) {
	if ic, ok := i.(hie.ClonableIter[T]); ok {
		return ic.Clone(), true
	}

	val := reflect.ValueOf(i)
	tpe := val.Type()
	mthd, hasClone := tpe.MethodByName("Clone")
	if hasClone {
		res := mthd.Func.Call([]reflect.Value{val})
		return res[0].Interface().(hie.Iter[T]), true
	}

	return i, false
}

func Cloned[T hie.Clonable[T]](in hie.Iter[T]) hie.Iter[T] {
	return &clonedIter[T]{under: in}
}

type clonedIter[T hie.Clonable[T]] struct {
	under hie.Iter[T]
}

func (c *clonedIter[T]) HasNext() bool {
	return c.under.HasNext()
}

func (c *clonedIter[T]) Next() T {
	return c.under.Next().Clone()
}
