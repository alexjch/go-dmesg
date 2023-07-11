// SPDX-License-Identifier: MIT
package dmesg

import (
	"testing"
)

func TestScanner(t *testing.T) {
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
