package mock

import (
	"context"
	"time"

	banking "github.com/morozovcookie/agat-banking"
	"github.com/stretchr/testify/mock"
)

var _ banking.Timer = (*Timer)(nil)

// Timer represents a service for retrieving time data.
type Timer struct {
	mock.Mock
}

// NewTimer returns a new Timer instance.
func NewTimer() *Timer {
	return &Timer{}
}

// Time returns time data.
func (t *Timer) Time(_ context.Context) (time.Time, error) {
	args := t.Called()

	return args.Get(0).(time.Time), args.Error(1)
}
