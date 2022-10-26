package hie

import (
	"fmt"
	"reflect"
)

type Defaulter[T any] func() T

type Option[T any] interface {
	IsNone() bool
	IsSome() bool
	Value() T
	ValueOrDefault(T) T
	ValueOrElse(Defaulter[T]) T
	isOption()

	AsIter[T]
}

func NewOption[T any](val T) Option[T] {
	tpe := reflect.TypeOf(val)
	switch tpe.Kind() {
	case reflect.Pointer, reflect.UnsafePointer, reflect.Func, reflect.Chan, reflect.Slice, reflect.Interface, reflect.Map:
		vv := reflect.ValueOf(val)
		if vv.IsNil() {
			return none[T]{}
		}
		return some[T]{
			value: val,
		}
	default:
		return some[T]{
			value: val,
		}
	}
}

func Some[T any](val T) Option[T] {
	return some[T]{value: val}
}

func None[T any]() Option[T] {
	return none[T]{}
}

type some[T any] struct {
	value T
}

func (some[T]) isOption()                                 {} //nolint:unused
func (some[T]) IsNone() bool                              { return false }
func (some[T]) IsSome() bool                              { return true }
func (s some[T]) Value() T                                { return s.value }
func (s some[T]) ValueOrDefault(defaultValue T) T         { return s.value }
func (s some[T]) ValueOrElse(defaultValue Defaulter[T]) T { return s.value }
func (s some[T]) AsIter() Iter[T]                         { return Slice(s.value).AsIter() }

type none[T any] struct {
}

func (none[T]) isOption()                               {} //nolint:unused
func (none[T]) IsNone() bool                            { return true }
func (none[T]) IsSome() bool                            { return false }
func (n none[T]) Value() T                              { panic(fmt.Sprintf("%T doesn't have a value", n)) }
func (none[T]) ValueOrDefault(defaultValue T) T         { return defaultValue }
func (none[T]) ValueOrElse(defaultValue Defaulter[T]) T { return defaultValue() }
func (n none[T]) AsIter() Iter[T]                       { return Slice[T]().AsIter() }
