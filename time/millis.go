package time

// MillisecondToNanosecondMultiplier is the multiplier for milliseconds to nanoseconds conversion or milliseconds to
// nanoseconds.
const MillisecondToNanosecondMultiplier = 1e6

func MillisecondsToNanoseconds(millis int64) int64 {
	return millis * MillisecondToNanosecondMultiplier
}

func NanosecondsToMilliseconds(ns int64) int64 {
	return ns / MillisecondToNanosecondMultiplier
}
