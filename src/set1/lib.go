package set1

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/bits"
	"os"
	"regexp"
	"sort"
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

	infreq := []byte("`'~@")
	uncommonLetters := make(map[byte]bool)

	for i := 0; i < len(infreq); i++ {
		uncommonLetters[infreq[i]] = true
	}

	score := 0
	for i := 0; i < len(s); i++ {
		if s[i] == 0x0D {
			return 0
		} else if _, ok := commonLetters[s[i]]; ok {
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

func encodeRepeatingKeyXOR(key []byte, text []byte) []byte {
	// assume we aren't dealing with gigs of text here
	textlen := len(text)
	repeat := (textlen / len(key)) + 1 // one more than needed...

	fullkey := bytes.Repeat(key, repeat)
	return XOR(fullkey[:textlen], text) // ... then sliced off
}

func hamming_distance(a, b []byte) int {
	val := XOR(a, b)
	sum := 0
	for _, bt := range val {
		sum += bits.OnesCount8(bt)
	}
	return sum
}

func readB64File(filename string) []byte {
	handle, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	b, err := ioutil.ReadAll(handle)
	if err != nil {
		log.Fatal(err)
	}
	decoded := make([]byte, len(b))
	declen, err := base64.StdEncoding.Decode(decoded, b)
	// hack: this is too big by 4/3, so
	decoded = bytes.Trim(decoded, "\x00")
	if err != nil {
		log.Fatal(declen, err)
	}
	return decoded
}

func findProbableKeyLengths(encoded_blob []byte) []int {
	type score struct {
		Keylen int
		Value  float64
	}

	var ham_scores []score
	for length := 2; length <= 40; length++ {
		sample_dist := 0
		for i := 0; i < 3*length; i += length {
			sample_a := encoded_blob[i:(i + length)]
			sample_b := encoded_blob[(i + length):(i + length*2)]

			sample_dist += hamming_distance(sample_a, sample_b)
		}

		ham_scores = append(ham_scores, score{length, float64(sample_dist) / float64(length)})
	}

	sort.Slice(ham_scores, func(i, j int) bool {
		return ham_scores[i].Value < ham_scores[j].Value
	})

	fmt.Fprintln(os.Stderr, ham_scores)

	var probable_keylengths []int
	for i, score := range ham_scores {
		probable_keylengths = append(probable_keylengths, score.Keylen)
		if i >= 4 {
			break
		}
	}

	return probable_keylengths
}

func BruteForceRKXOR(filename string) ([]byte, []byte, error) {
	encoded_blob := readB64File(filename)
	probable_keylengths := findProbableKeyLengths(encoded_blob)
	fmt.Fprintln(os.Stderr, probable_keylengths)

	for _, keylength := range probable_keylengths {
		slices := make([][]byte, keylength)

		for i := 0; i < len(encoded_blob); i++ {
			slices[i%keylength] = append(slices[i%keylength], encoded_blob[i])
		}

		var probable_key bytes.Buffer
		for _, slice := range slices {
			top_score := 0
			var c byte
			for i := byte(0x20); i <= 0x7E; i++ {
				mask := bytes.Repeat([]byte{i}, len(slice))
				result := XOR(slice, mask)
				score := scoreText(result)
				if score > top_score {
					top_score = score
					c = i
				}
			}
			if c == 0 {
				continue
			}
			err := probable_key.WriteByte(c)
			if err != nil {
				log.Fatal(err)
			}
		}
		if probable_key.Len() != keylength {
			fmt.Fprintf(os.Stderr, "Found no suitable ASCII bytes to decode with key length %s\n", keylength)
		} else {
			key := probable_key.Bytes()
			return key, encodeRepeatingKeyXOR(key, encoded_blob), nil
		}
	}

	return nil, nil, fmt.Errorf("Found no suitable keys, tried lengths %s\n", probable_keylengths)
}
