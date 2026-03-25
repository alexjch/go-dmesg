// SPDX-License-Identifier: MIT
package dmesg

import (
	"os"
	"testing"
)

func TestScanner(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("skipping test: requires root privileges to read /dev/kmsg")
	}
	s, e := NewScanner()
	if e != nil {
		t.Fatalf("NewScanner() error: %v", e)
	}
	s.Scan()
	if s.Text() == "" {
		t.Fatal("Scanner.Text() returned empty string")
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Scanner.Close() error: %v", err)
	}
}
