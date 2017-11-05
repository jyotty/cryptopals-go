package set1

import (
	"fmt"
)

func ExampleHammingDistance() {
	fmt.Println(hamming_distance([]byte("this is a test"), []byte("wokka wokka!!!")))
	// Output: 37
}
