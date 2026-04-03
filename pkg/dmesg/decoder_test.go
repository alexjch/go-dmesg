// SPDX-License-Identifier: MIT
package dmesg

import (
	"bufio"
	"context"
	"io"
	"strings"
	"testing"
	"time"
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
	}, context.Background())
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

func TestDecoderFollowEarlyCancellation(t *testing.T) {
	raw := strings.Join([]string{
		"6,1,1000,-;first message",
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

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	called := false
	done := make(chan error, 1)
	go func() {
		done <- decoder.Follow(func(r *Record) {
			called = true
		}, ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Follow() error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Follow() did not return after context cancellation")
	}

	if called {
		t.Fatal("callback should not be called when context is canceled before Follow")
	}
}
