// SPDX-License-Identifier: MIT
package dmesg

import (
	"errors"
)

// Decoder is a dmesg.Decoder that wraps a dmesg.Scanner
// and decodes each line returned by dmesg.Scanner as a
// dmesg.Record
type Decoder struct {
	*MonotonicDecoder
	scanner *Scanner
	stop    chan bool
}

// NewDecoder returns a dmesg.Decoder that wraps a dmesg.Scanner
func NewDecoder(s *Scanner) (*Decoder, error) {
	md, err := NewMonotonicDecoder()
	if err != nil {
		return nil, err
	}
	return &Decoder{
		scanner:          s,
		stop:             make(chan bool),
		MonotonicDecoder: md,
	}, nil
}

func (d *Decoder) decode(line string) (*Record, error) {
	matches := kmsgMatches(line)

	switch {
	case continuation(line):
		return nil, nil
	case len(matches) != 5:
		return nil, errors.New("invalid line: " + line)
	}

	p, s, mon, err := parseLogPrefix(matches)
	ts := d.RealTime(float64(mon))

	if err != nil {
		return nil, err
	}

	return &Record{
		Priority:  p,
		Sequence:  s,
		Timestamp: ts,
		Message:   matches[4],
	}, nil
}

// scan reads from dmesg.Scanner and sends each line to the
// channel line. In case of an error the error is sent to the
// channel line. The type of the value sent to the channel line
// should be asserted by the caller.
func (d *Decoder) scan(line chan interface{}) {
	for d.scanner.Scan() {
		line <- d.scanner.Text()
	}
	line <- d.scanner.Err()
}

// notify receives a line from dmesg.Scanner and notifies
// the callback function f passed as argument with a dmesg.Record.
// This function returns an error whjen the decoder fails to decode
func (d *Decoder) notify(line interface{}, f func(*Record)) error {
	switch v := line.(type) {
	case string:
		record, err := d.decode(v)
		if err != nil {
			return err
		}
		// Skip continuation lines
		if record != nil {
			f(record)
		}
	case error:
		return v
	default:
		return errors.New("invalid type")
	}
	return nil
}

// Follow receives one line at a time from dmesg.Scanner ,decodes
// it as a dmesg.Record, and passes the dmesg.Record to the callback
// function f passed as argument. It returns an error if the scanner
// returns an error. The function blocks until dmesg.Decoder.Stop()
// is called.
func (d *Decoder) Follow(f func(*Record)) error {
	c := make(chan interface{})
	// Scan in a separate goroutine
	go d.scan(c)
	// Loop to read stop signal or notify
	// the callback function with a dmesg.Record
	for {
		select {
		case <-d.stop:
			return nil
		case line := <-c:
			if err := d.notify(line, f); err != nil {
				return err
			}
		}
	}
}

// Stop stops the dmesg.Decoder Follow() blocking call
func (d *Decoder) Stop() {
	d.stop <- true
}
