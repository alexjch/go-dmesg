// SPDX-License-Identifier: MIT
package dmesg

import (
	"errors"
	"sync"
)

// Decoder is a dmesg.Decoder that wraps a dmesg.KmsgScanner
// and decodes each line returned by dmesg.KmsgScanner as a
// dmesg.Record
type Decoder struct {
	*MonotonicDecoder
	scanner *KmsgScanner
	stop    chan struct{}
	once    sync.Once
}

// NewDecoder returns a dmesg.Decoder that wraps a dmesg.KmsgScanner
func NewDecoder(s *KmsgScanner) (*Decoder, error) {
	md, err := NewMonotonicDecoder()
	if err != nil {
		return nil, err
	}
	return &Decoder{
		scanner:          s,
		stop:             make(chan struct{}),
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
func (d *Decoder) scan(line chan interface{}) {
	for d.scanner.Scan() {
		text := d.scanner.Text()
		select {
		case <-d.stop:
			// Stop signal received; exit without blocking on send.
			return
		case line <- text:
		}
	}
	if err := d.scanner.Err(); err != nil {
		select {
		case <-d.stop:
			// Stop signal received; exit without blocking on send.
			return
		case line <- err:
		}
	}
	close(line)
}

// notify receives a line from dmesg.KmsgScanner and notifies
// the callback function f passed as argument with a dmesg.Record.
// This function returns an error when the decoder fails to decode
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

// Follow receives one line at a time from dmesg.KmsgScanner, decodes
// it as a dmesg.Record, and passes the dmesg.Record to the callback
// function f passed as argument. It returns an error if the scanner
// returns an error. The function blocks until dmesg.Decoder.Stop()
// is called.
func (d *Decoder) Follow(f func(*Record)) error {
	c := make(chan interface{})
	// Scan in a separate goroutine
	go d.scan(c)
	// Ensure that the scan goroutine is signalled to stop when Follow exits.
	defer d.Stop()
	// Loop to read stop signal or notify
	// the callback function with a dmesg.Record
	for {
		select {
		case <-d.stop:
			return nil
		case line, ok := <-c:
			if !ok {
				return nil
			}
			if err := d.notify(line, f); err != nil {
				return err
			}
		}
	}
}

// Stop stops the dmesg.Decoder Follow() blocking call
func (d *Decoder) Stop() {
	d.once.Do(func() {
		close(d.stop)
	})
}
