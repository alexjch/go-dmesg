// SPDX-License-Identifier: MIT
package dmesg

import (
	"time"
)

type MonotonicDecoder struct {
	boot time.Time
}

// Takes a monotonic timestamp and returns the real time.
// the monotonic timestamp should be in microseconds.
func (m *MonotonicDecoder) RealTime(timestamp float64) time.Time {
	return m.boot.Add(time.Duration(timestamp) * time.Microsecond)
}

func NewMonotonicDecoder() (*MonotonicDecoder, error) {
	uptime, err := osUptime()
	if err != nil {
		return nil, err
	}
	return &MonotonicDecoder{
		boot: time.Now().Add(uptime * -1).UTC(),
	}, nil
}
