package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const ArgumentsRequiredCount = 2

func main() {
	if len(os.Args) < ArgumentsRequiredCount {
		log.Fatalln("missed bytes count parameter")
	}

	size, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	bb := make([]byte, size)

	if _, err := io.ReadFull(rand.Reader, bb); err != nil {
		log.Fatalln(err)
	}

	_, _ = fmt.Fprintln(os.Stdout, hex.EncodeToString(bb))
}
