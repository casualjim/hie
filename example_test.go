package hie_test

import (
	"fmt"

	"github.com/casualjim/hie"
)

func ExampleFlatMap() {
	slices := hie.Slice(hie.Slice(1, 2).AsIter(), hie.Slice(3, 4).AsIter(), hie.Slice(5, 6).AsIter())

	iter := hie.FlatMap(slices.AsIter(), func(i hie.Iter[int]) hie.Iter[string] {
		return hie.Map(i, func(i int) string { return fmt.Sprintf("%d", i) })
	})

	fmt.Printf("%#v\n", hie.Collect(iter))
	// Output: []string{"1", "2", "3", "4", "5", "6"}
}
