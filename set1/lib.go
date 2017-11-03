package set1

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	//	"fmt"
	"log"
	"os"
	"regexp"
)

func hexToB64(hs string) string {
	b, err := hex.DecodeString(hs)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	encoder.Write(b)
	encoder.Close()
	return buf.String()
}

func XOR(a, b []byte) []byte {
	if len(a) != len(b) {
		log.Fatal("can't xor unequal length strings, pad")
	}

	l := len(a)
	dest := make([]byte, l, l)

	for i := 0; i < l; i++ {
		dest[i] = a[i] ^ b[i]
	}

	return dest
}
func hexXOR(a, b string) string {
	if len(a) != len(b) {
		log.Fatal("can't xor unequal length strings, pad")
	}

	ba, err := hex.DecodeString(a)
	if err != nil {
		log.Fatal(err)
	}

	bb, err := hex.DecodeString(b)
	if err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(XOR(ba, bb))
}

func singleByteXOR(char byte, test []byte) []byte {
	l := len(test)
	a := make([]byte, l, l)

	for i := 0; i < l; i++ {
		a[i] = char
	}

	return XOR(a, test)
}

func scoreText(s []byte) int {
	// are there any non-printable bytes?
	match, err := regexp.Match("[^[:print:][:space:]]", s)
	if err != nil {
		log.Fatal(err)
	}

	// if so, return 0
	if match {
		return 0
	}

	freq := []byte("ETAOINSHRDLUetaoinshrdlu ")
	commonLetters := make(map[byte]bool)

	for i := 0; i < len(freq); i++ {
		commonLetters[freq[i]] = true
	}

	infreq := []byte("`'")
	uncommonLetters := make(map[byte]bool)

	for i := 0; i < len(infreq); i++ {
		uncommonLetters[infreq[i]] = true
	}

	score := 0
	for i := 0; i < len(s); i++ {
		if _, ok := commonLetters[s[i]]; ok {
			score++
		} else if _, ok := uncommonLetters[s[i]]; ok {
			score--
		}
	}

	return score
}

func bruteSingleByteXOR(s []byte) (int, []byte) {
	max := 0
	bestCandidate := make([]byte, len(s), len(s))

	for i := 0; i < 256; i++ {
		candidate := singleByteXOR(byte(i), s)
		score := scoreText(candidate)

		if score > max {
			//fmt.Fprintf(os.Stderr, "%v: %s\n", score, candidate)
			max = score
			bestCandidate = candidate
		}
	}

	return max, bestCandidate
}

func bruteForceLines(filename string) []byte {
	handle, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	var lines [][]byte
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		encb, err := hex.DecodeString(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		lines = append(lines, encb)
	}

	max := 0
	bestCandidate := make([]byte, len(lines[0]), len(lines[0]))

	for _, line := range lines {
		score, dec := bruteSingleByteXOR(line)

		if score > max {
			max = score
			bestCandidate = dec
		}
	}
	return bestCandidate
}

func encodeRepeatingKeyXOR(key string, text string) []byte {
	textb := []byte(text)
	keyb := []byte(key)

	// assume we aren't dealing with gigs of text here
	textlen := len(textb)
	repeat := (textlen / len(keyb)) + 1 // one more than needed...

	fullkey := bytes.Repeat(keyb, repeat)
	return XOR(fullkey[:textlen], textb) // ... then sliced off
}
