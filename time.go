package banking

import (
	"context"
	"time"
)

// Timer represents a service for retrieving time data.
type Timer interface {
	// Time returns time data.
	Time(ctx context.Context) (time.Time, error)
}
