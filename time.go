package banking

import (
	"context"
	"time"
)

// MillisecondToNanosecondMultiplier is the multiplier for milliseconds to nanoseconds conversion or milliseconds to
// nanoseconds.
const MillisecondToNanosecondMultiplier = 1e6

func MillisecondsToNanoseconds(millis int64) int64 {
	return millis * MillisecondToNanosecondMultiplier
}

func NanosecondsToMilliseconds(ns int64) int64 {
	return ns / MillisecondToNanosecondMultiplier
}

func TimeToMilliseconds(t time.Time) int64 {
	return NanosecondsToMilliseconds(t.UnixNano())
}

func MillisecondsToTime(millis int64) time.Time {
	return time.Unix(0, MillisecondsToNanoseconds(millis))
}

// Timer represents a service for retrieving time data.
type Timer interface {
	// Time returns time data.
	Time(ctx context.Context) (time.Time, error)
}
