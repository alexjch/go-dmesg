// SPDX-License-Identifier: MIT
package dmesg

import (
	"bufio"
	"io"
	"os"
	"syscall"
)

const devkmsg = "/dev/kmsg"

// Scanner combines around bufio.Scanner and io.ReadCloser
// to wrap /dev/kmsg in a non-blocking manner.
type Scanner struct {
	bufio.Scanner
	io.ReadCloser
}

func (s Scanner) Close() error {
	return s.ReadCloser.Close()
}

// NewScanner returns a dmesg.Scanner this primitive can
// be used with dmesg.Decoder to read kernel messages as
// dmesg.Record(s).
func NewScanner() (*Scanner, error) {
	// Open /dev/kmsg for reading in a non-blocking manner
	fd, err := syscall.Open(devkmsg, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	// Check if the file descriptor is valid
	if fd < 0 {
		return nil, os.NewSyscallError("open", err)
	}
	// Wrap a file descriptor in an os.File
	f := os.NewFile(uintptr(fd), "")
	// Seek to the end of the file
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return nil, err
	}
	// Wrap an os.File in a bufio.Scanner
	scanner := bufio.NewScanner(f)
	return &Scanner{
		Scanner:    *scanner,
		ReadCloser: f,
	}, nil
}
