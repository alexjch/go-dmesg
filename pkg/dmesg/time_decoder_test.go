// SPDX-License-Identifier: MIT
package dmesg

import (
	"fmt"
	"testing"
	"time"
)

const mSec = 1000

func TestOsUptime(t *testing.T) {
	uptime, err := osUptime()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	fmt.Println(uptime)
}

func TestTimeDecoder(t *testing.T) {
	cases := map[string]struct {
		input    string
		ts       float64
		validate func(time.Time, time.Time) bool
	}{
		"boot time": {input: "14100511.99", ts: 0 * mSec,
			validate: func(boot, decoded time.Time) bool {
				return boot.Equal(decoded)
			}},
		"1 sec after": {input: "14100512.99", ts: 1 * mSec,
			validate: func(boot, decoded time.Time) bool {
				return decoded.Sub(boot) == time.Duration(1*time.Millisecond)
			}},
	}

	for name, c := range cases {
		uptime, _ := time.ParseDuration(c.input)
		// Use the beginning of times for testing, having a
		// knonw boot time helps with repitability.
		boot := time.Unix(0, 0).Add(uptime * -1).UTC()
		md := &MonotonicDecoder{
			boot: boot,
		}
		decoded := md.RealTime(c.ts)
		// Check if the boot time is decoded correctly.
		if !c.validate(boot, decoded) || boot == time.Unix(0, 0) {
			t.Fatalf("in test %s expected %v, got %v", name, boot, decoded)
		}
	}
}
