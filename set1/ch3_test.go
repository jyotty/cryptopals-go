package set1

import (
	"encoding/hex"
	"fmt"
	"log"
)

func ExampleBruteSingleByteXOR() {
	enc := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	encb, err := hex.DecodeString(enc)
	if err != nil {
		log.Fatal(err)
	}
	_, dec := bruteSingleByteXOR(encb)
	fmt.Printf("%s\n", dec)
	// Output: Cooking MC's like a pound of bacon
}
