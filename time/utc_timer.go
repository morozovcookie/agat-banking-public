package time

import (
	"context"
	"time"

	banking "github.com/morozovcookie/agat-banking"
)

var _ banking.Timer = (*UTCTimer)(nil)

// UTCTimer represents a service for retrieving time data.
type UTCTimer struct{}

// NewUTCTimer returns a new UTCTimer instance.
func NewUTCTimer() *UTCTimer {
	return &UTCTimer{}
}

// Time returns time data.
func (*UTCTimer) Time(_ context.Context) (time.Time, error) {
	return time.Now().UTC(), nil
}
