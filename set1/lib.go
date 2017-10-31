package set1

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"log"
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
