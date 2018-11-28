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

// validateTimestamp is copied from "google/protobuf/timestamp.proto"
// should have imported this pkg, but failed to set the environment correctly
// so now I'm rewriting the function and using my own Timestamp struct instead
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

// TimestampProto is copied from "google/protobuf/timestamp.proto"
// should have imported this pkg, but failed to set the environment correctly
// so now I'm rewriting the function and using my own Timestamp struct instead
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

// Timestamp convert protoTimestamp to time.Time
// we store protoTimestamp in our storage, but when frontend is displaying time,
// this function will be called to convert protoTimestamp to time.Time in golang,
// time.Time is then converted to readable time and date in front-end implementation.
func Timestamp(ts *protocol.Timestamp) time.Time {
	var t time.Time
	if ts == nil {
		t = time.Unix(0, 0).UTC() // treat nil like the empty Timestamp
	} else {
		t = time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	}
	return t
}
