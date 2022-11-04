package opt

import (
	"fmt"
	"reflect"

	"github.com/casualjim/hie"
)

type Defaulter[T any] func() T

type Option[T any] interface {
	IsNone() bool
	IsSome() bool
	Value() T
	ValueOrDefault() T
	ValueOr(T) T
	ValueOrElse(Defaulter[T]) T
	AsIter() hie.Iter[T]
	isOption()
}

func New[T any](val T) Option[T] {
	tpe := reflect.TypeOf(val)
	switch tpe.Kind() {
	case reflect.Pointer, reflect.UnsafePointer, reflect.Func, reflect.Chan, reflect.Slice, reflect.Interface, reflect.Map:
		vv := reflect.ValueOf(val)
		if vv.IsNil() {
			return none[T]{}
		}
		return &some[T]{
			value: val,
		}
	default:
		return &some[T]{
			value: val,
		}
	}
}

func Some[T any](val T) Option[T] {
	return &some[T]{value: val}
}

func None[T any]() Option[T] {
	return none[T]{}
}

type some[T any] struct {
	value T
}

func (some[T]) isOption()    {} //nolint:unused
func (some[T]) IsNone() bool { return false }
func (some[T]) IsSome() bool { return true }
func (s *some[T]) Value() T  { return s.value }
func (s *some[T]) ValueOrDefault() T {
	return s.value
}
func (s *some[T]) ValueOr(defaultValue T) T                { return s.value }
func (s *some[T]) ValueOrElse(defaultValue Defaulter[T]) T { return s.value }
func (s *some[T]) AsIter() hie.Iter[T]                     { return &optionIter[T]{val: s} }

type none[T any] struct {
}

func (none[T]) isOption()    {} //nolint:unused
func (none[T]) IsNone() bool { return true }
func (none[T]) IsSome() bool { return false }
func (n none[T]) Value() T   { panic(fmt.Sprintf("%T doesn't have a value", n)) }
func (none[T]) ValueOrDefault() T {
	var zero T
	return zero
}
func (none[T]) ValueOr(defaultValue T) T {
	return defaultValue
}
func (none[T]) ValueOrElse(defaultValue Defaulter[T]) T { return defaultValue() }
func (n none[T]) AsIter() hie.Iter[T]                   { return &optionIter[T]{val: n} }

type optionIter[T any] struct {
	val      Option[T]
	consumed bool
}

func (o *optionIter[T]) HasNext() bool {
	return !o.consumed && o.val.IsSome()
}

func (o *optionIter[T]) Next() T {
	if o.consumed {
		panic("next called on a consumed option iter")
	}
	o.consumed = true
	return o.val.Value()
}
