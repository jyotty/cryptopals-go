package set1

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	//"fmt"
	"log"
	//"os"
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

	l := len(ba)
	dest := make([]byte, l, l)

	for i := 0; i < l; i++ {
		dest[i] = ba[i] ^ bb[i]
	}

	return hex.EncodeToString(dest)
}
