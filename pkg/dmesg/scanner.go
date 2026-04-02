// SPDX-License-Identifier: MIT
package dmesg

import (
	"bufio"
	"io"
	"os"
	"syscall"
)

const devkmsg = "/dev/kmsg"

// readSeekCloser narrows the dependency used by NewScanner so tests can
// inject a mock and exercise open/seek/close error paths.
type readSeekCloser interface {
	io.ReadSeeker
	io.Closer
}

var openKmsg = func() (readSeekCloser, error) {
	fd, err := syscall.Open(devkmsg, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), ""), nil
}

// KmsgScanner wraps /dev/kmsg in a non-blocking manner and exposes
// both scanning and close operations.
//
// The function that creates a KmsgScanner is responsible for calling Close
// when scanning is done.
type KmsgScanner struct {
	bufio.Scanner
	io.ReadCloser
}

// Close closes the underlying /dev/kmsg file descriptor.
func (s KmsgScanner) Close() error {
	return s.ReadCloser.Close()
}

// NewScanner creates a KmsgScanner that tails /dev/kmsg from the current end
// in non-blocking mode.
//
// The caller that instantiates the KmsgScanner should defer Close to release
// the underlying file descriptor.
func NewScanner() (*KmsgScanner, error) {
	// Open /dev/kmsg for reading in a non-blocking manner
	f, err := openKmsg()
	if err != nil {
		return nil, err
	}

	// Seek to the end of the file
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		_ = f.Close()
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	return &KmsgScanner{
		Scanner:    *scanner,
		ReadCloser: f,
	}, nil
}
