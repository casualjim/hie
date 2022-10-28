package hie_test

import (
	"fmt"

	"github.com/casualjim/hie"
)

func ExampleFlatMap() {
	asIter := hie.Slice(hie.Slice(1, 2), hie.Slice(3, 4), hie.Slice(5, 6))

	iter := hie.FlatMap(asIter, func(i hie.Iter[int]) hie.Iter[string] {
		return hie.Map(i, func(i int) string { return fmt.Sprintf("%d", i) })
	})

	fmt.Printf("%#v\n", hie.Collect(iter))
	// Output: []string{"1", "2", "3", "4", "5", "6"}
}
