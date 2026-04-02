// SPDX-License-Identifier: MIT
package dmesg

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestDecoder(t *testing.T) {
	raw := strings.Join([]string{
		"6,1,1000,-;first message",
		" continuation line",
		"5,2,2000,-;second message",
	}, "\n")

	rc := io.NopCloser(strings.NewReader(raw))
	s := &KmsgScanner{
		Scanner:    *bufio.NewScanner(rc),
		ReadCloser: rc,
	}

	decoder, err := NewDecoder(s)
	if err != nil {
		t.Fatalf("NewDecoder() error: %v", err)
	}

	var got []*Record
	err = decoder.Follow(func(r *Record) {
		got = append(got, r)
	})
	if err != nil {
		t.Fatalf("Follow() error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}
	if got[0].Message != "first message" {
		t.Fatalf("unexpected first message: %q", got[0].Message)
	}
	if got[1].Message != "second message" {
		t.Fatalf("unexpected second message: %q", got[1].Message)
	}
}
