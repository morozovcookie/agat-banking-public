package main

import (
	"context"
	"fmt"
	"os"

	"github.com/morozovcookie/agat-banking/time"
)

const MillisecondsMultiplier = 1e6

func main() {
	var (
		ctx   = context.Background()
		timer = time.NewUTCTimer()
	)

	now, err := timer.Time(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	_, _ = fmt.Fprintln(os.Stdout, now.UnixNano()/MillisecondsMultiplier)
}
