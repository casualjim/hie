package iter

import (
	"reflect"

	"github.com/casualjim/hie"
)

func IsClonable[T any](i hie.Iter[T]) bool {
	if _, ok := i.(hie.ClonableIter[T]); ok {
		return true
	}
	_, hasClone := reflect.ValueOf(i).Type().MethodByName("Clone")
	return hasClone
}

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
	return Map(in, func(i T) T { return i.Clone() })
}
