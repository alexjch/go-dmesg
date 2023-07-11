// SPDX-License-Identifier: MIT
package main

import (
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
	// Handle Ctrl-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
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
	// Stop the decoder when Ctrl-C is pressed
	go func() {
		<-c
		dec.Stop()
	}()
	// Follow the decoder until Ctrl-C is pressed
	if err := dec.Follow(func(r *dmesg.Record) {
		fmt.Println(r.String())
	}); err != nil {
		panic(err)
	}
}
