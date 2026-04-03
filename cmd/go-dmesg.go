// SPDX-License-Identifier: MIT
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/alexjch/go-dmesg/pkg/dmesg"
)

// Primitives are:
// - scanner, a bufio.Scanner that reads from /dev/kmsg
// - decoder, a dmesg.Decoder that wraps a scanner and decodes log records
// 	{decoder{scanner}}}

func main() {
	// Build a context canceled on Ctrl-C.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Create a scanner
	scan, err := dmesg.NewScanner()
	if err != nil {
		panic(err)
	}
	defer scan.Close()
	// Create a decoder
	dec, err := dmesg.NewDecoder(scan)
	if err != nil {
		panic(err)
	}

	if err := dec.Follow(func(r *dmesg.Record) {
		fmt.Println(r.String())
	}, ctx); err != nil {
		panic(err)
	}
}
