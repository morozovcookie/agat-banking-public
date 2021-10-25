package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

const (
	MissedInputExitCode      = 1
	WriteInputStringExitCode = 2

	RequiredArgumentsCount = 2
)

func main() {
	if len(os.Args) < RequiredArgumentsCount {
		_, _ = fmt.Fprintln(os.Stderr, "missed input string")

		os.Exit(MissedInputExitCode)
	}

	var (
		sha   = sha256.New()
		input = os.Args[1]
	)

	if _, err := sha.Write(append([]byte{}, input...)); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(WriteInputStringExitCode)
	}

	_, _ = fmt.Fprintln(os.Stdout, hex.EncodeToString(sha.Sum(nil)))
}
