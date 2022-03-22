// Package hie provides a generic collection library for golang
package hie

type AsSlice[T any] interface {
	AsSlice() []T
}
