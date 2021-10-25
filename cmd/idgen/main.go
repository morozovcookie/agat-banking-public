package main

import (
	"context"
	"fmt"
	"os"

	"github.com/morozovcookie/agat-banking/nanoid"
)

func main() {
	var (
		ctx = context.Background()
		gen = nanoid.NewIdentifierGenerator()
	)

	id, err := gen.GenerateIdentifier(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	_, _ = fmt.Fprintln(os.Stdout, id)
}
