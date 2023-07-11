// SPDX-License-Identifier: MIT
package dmesg

import (
	"fmt"
	"time"
)

type Record struct {
	Priority  uint8
	Sequence  uint64
	Timestamp time.Time
	Message   string
}

func (r *Record) String() string {
	return fmt.Sprintf("%d %d %s %s", r.Priority, r.Sequence, r.Timestamp.String(), r.Message)
}
