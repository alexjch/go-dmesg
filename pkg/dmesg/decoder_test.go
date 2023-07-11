// SPDX-License-Identifier: MIT
package dmesg

import (
	"fmt"
	"testing"
)

func TestDecoder(t *testing.T) {
	s, err := NewScanner()
	if err != nil {
		t.Fatalf("NewScanner() error: %v", err)
	}
	decoder, err := NewDecoder(s)
	if err != nil {
		t.Fatalf("NewDecoder() error: %v", err)
	}
	go func() {
		decoder.Stop()
	}()
	err = decoder.Follow(func(r *Record) {
		fmt.Println(r.String())
	})
	if err != nil {
		t.Fatalf("Follow() error: %v", err)
	}
}
