package main

import (
	"fmt"
	"log"
	"set1"
)

func main() {
	key, decoded, err := set1.BruteForceRKXOR("src/set1/6.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Key: %s\nLyrics: %s", key, decoded)
}
