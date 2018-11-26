package twitterTimestamp

import (
	"errors"
	"fmt"
	"time"

	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
)

const (
	// Seconds field of the earliest valid Timestamp.
	// This is time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	minValidSeconds = -62135596800
	// Seconds field just after the latest valid Timestamp.
	// This is time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	maxValidSeconds = 253402300800
)

func validateTimestamp(ts *protocol.Timestamp) error {
	if ts == nil {
		return errors.New("timestamp: nil Timestamp")
	}
	if ts.Seconds < minValidSeconds {
		return fmt.Errorf("timestamp: %v before 0001-01-01", ts)
	}
	if ts.Seconds >= maxValidSeconds {
		return fmt.Errorf("timestamp: %v after 10000-01-01", ts)
	}
	if ts.Nanos < 0 || ts.Nanos >= 1e9 {
		return fmt.Errorf("timestamp: %v: nanos not in range [0, 1e9)", ts)
	}
	return nil
}
func TimestampProto(t time.Time) *protocol.Timestamp {
	seconds := t.Unix()
	nanos := int32(t.Sub(time.Unix(seconds, 0)))
	ts := &protocol.Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}
	if err := validateTimestamp(ts); err != nil {
		return nil
	}
	return ts
}
