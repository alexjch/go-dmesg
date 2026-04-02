// SPDX-License-Identifier: MIT
package dmesg

import (
	"errors"
	"strings"
	"testing"
)

type mockReadSeekCloser struct {
	reader      *strings.Reader
	seekErr     error
	closeErr    error
	seekCalled  bool
	closeCalled bool
}

func (m *mockReadSeekCloser) Read(p []byte) (int, error) {
	return m.reader.Read(p)
}

func (m *mockReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	m.seekCalled = true
	if m.seekErr != nil {
		return 0, m.seekErr
	}
	return m.reader.Seek(offset, whence)
}

func (m *mockReadSeekCloser) Close() error {
	m.closeCalled = true
	return m.closeErr
}

func TestNewScannerOpenError(t *testing.T) {
	orig := openKmsg
	t.Cleanup(func() { openKmsg = orig })

	expectedErr := errors.New("open failed")
	openKmsg = func() (readSeekCloser, error) {
		return nil, expectedErr
	}

	_, err := NewScanner()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected open error %v, got %v", expectedErr, err)
	}
}

func TestNewScannerSeekErrorClosesFile(t *testing.T) {
	orig := openKmsg
	t.Cleanup(func() { openKmsg = orig })

	expectedErr := errors.New("seek failed")
	mock := &mockReadSeekCloser{
		reader:  strings.NewReader("line one\n"),
		seekErr: expectedErr,
	}

	openKmsg = func() (readSeekCloser, error) {
		return mock, nil
	}

	_, err := NewScanner()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected seek error %v, got %v", expectedErr, err)
	}
	if !mock.seekCalled {
		t.Fatal("expected Seek to be called")
	}
	if !mock.closeCalled {
		t.Fatal("expected Close to be called on seek error")
	}
}

func TestNewScannerSuccessAndClose(t *testing.T) {
	orig := openKmsg
	t.Cleanup(func() { openKmsg = orig })

	mock := &mockReadSeekCloser{
		reader: strings.NewReader("line one\n"),
	}

	openKmsg = func() (readSeekCloser, error) {
		return mock, nil
	}

	s, err := NewScanner()
	if err != nil {
		t.Fatalf("NewScanner() error: %v", err)
	}
	if !mock.seekCalled {
		t.Fatal("expected Seek to be called")
	}

	if err := s.Close(); err != nil {
		t.Fatalf("Close() error: %v", err)
	}
	if !mock.closeCalled {
		t.Fatal("expected Close to be called")
	}
}

func TestKmsgScannerClosePropagatesError(t *testing.T) {
	expectedErr := errors.New("close failed")
	mock := &mockReadSeekCloser{
		reader:   strings.NewReader(""),
		closeErr: expectedErr,
	}

	s := KmsgScanner{
		ReadCloser: mock,
	}

	err := s.Close()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected close error %v, got %v", expectedErr, err)
	}
}
