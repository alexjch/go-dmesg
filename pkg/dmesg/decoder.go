// SPDX-License-Identifier: MIT
package dmesg

import (
	"context"
	"errors"
)

// Decoder is a dmesg.Decoder that wraps a dmesg.KmsgScanner
// and decodes each line returned by dmesg.KmsgScanner as a
// dmesg.Record
type Decoder struct {
	*MonotonicDecoder
	scanner *KmsgScanner
}

// NewDecoder returns a dmesg.Decoder that wraps a dmesg.KmsgScanner
func NewDecoder(s *KmsgScanner) (*Decoder, error) {
	md, err := NewMonotonicDecoder()
	if err != nil {
		return nil, err
	}
	return &Decoder{
		scanner:          s,
		MonotonicDecoder: md,
	}, nil
}

func (d *Decoder) decodeLine(line string) (*Record, error) {
	matches := kmsgMatches(line)

	switch {
	case continuation(line):
		return nil, nil
	case len(matches) != 5:
		return nil, errors.New("invalid line: " + line)
	}

	p, s, mon, err := parseLogPrefix(matches)
	if err != nil {
		return nil, err
	}
	ts := d.RealTime(float64(mon))

	return &Record{
		Priority:  p,
		Sequence:  s,
		Timestamp: ts,
		Message:   matches[4],
	}, nil
}

// scan reads from dmesg.KmsgScanner and sends each line to the
// channel line. In case of an error the error is sent to the
// channel line. The type of the value sent to the channel line
// should be asserted by the caller.
func (d *Decoder) scan(line chan interface{}, ctx context.Context) {
	for d.scanner.Scan() {
		text := d.scanner.Text()
		select {
		case <-ctx.Done():
			// Stop signal received; exit without blocking on send.
			return
		case line <- text:
		}
	}
	if err := d.scanner.Err(); err != nil {
		select {
		case <-ctx.Done():
			// Stop signal received; exit without blocking on send.
			return
		case line <- err:
		}
	}
	close(line)
}

// decode converts a scanned item into a dmesg.Record.
// It returns nil, nil for continuation lines.
func (d *Decoder) decode(line interface{}) (*Record, error) {
	switch v := line.(type) {
	case string:
		record, err := d.decodeLine(v)
		if err != nil {
			return nil, err
		}
		// Skip continuation lines
		if record != nil {
			return record, nil
		}
		return nil, nil
	case error:
		return nil, v
	default:
		return nil, errors.New("invalid type")
	}
}

// Follow receives one line at a time from dmesg.KmsgScanner, decodes
// it as a dmesg.Record, and passes the dmesg.Record to the callback
// function f passed as argument. It returns an error if scanning or
// decoding fails. The function returns when the scanner is exhausted
// or when ctx is canceled.
func (d *Decoder) Follow(f func(*Record), ctx context.Context) error {
	c := make(chan interface{})
	// Scan in a separate goroutine
	go d.scan(c, ctx)
	// Loop to read stop signal or notify
	// the callback function with a dmesg.Record
	for {
		select {
		case <-ctx.Done():
			return nil
		case line, ok := <-c:
			if !ok {
				return nil
			}
			record, err := d.decode(line)
			if err != nil {
				return err
			}
			if record != nil {
				f(record)
			}
		}
	}
}
