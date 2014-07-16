// Copyright (c) 2014 Brian Hetro <whee@smaertness.net>
// Use of this source code is governed by the ISC
// license which can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"

	docopt "github.com/docopt/docopt-go"
	"github.com/whee/rp"
)

func main() {
	usage := `Redis Pipe.

Usage:
  rp      -r <name>...
  rp [-p] -w <name>...

Options:
  -r, --read <name>...  Read from the named channel.
  -w, --write <name>...  Write to the named channel.
  -p, --passthrough  Pass written data to standard output.`

	arguments, err := docopt.Parse(usage, nil, true, "Redis Pipe 0.1", false)
	if err != nil {
		panic(err)
	}
	var names []string
	reading, writing, passthrough := false, false, false

	if w, ok := arguments["--write"]; ok {
		if ns := w.([]string); len(ns) > 0 {
			names = ns
			writing = true
			passthrough = arguments["--passthrough"].(bool)
		}
	}
	if r, ok := arguments["--read"]; ok {
		if ns := r.([]string); len(ns) > 0 {
			names = ns
			reading = true
		}
	}

	if writing {
		_, err = writeTo(names, passthrough)
	} else if reading {
		_, err = readFrom(names)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

type ptWriter struct {
	w io.Writer
}

func (w ptWriter) Write(p []byte) (n int, err error) {
	n, err = os.Stdout.Write(p)
	if err != nil {
		return
	}
	return w.w.Write(p)
}

func writeTo(names []string, passthrough bool) (int64, error) {
	t, err := rp.NewWriter(nil, names...)
	if err != nil {
		return 0, err
	}
	defer t.Close()

	var w io.Writer
	if passthrough {
		w = ptWriter{t}
	} else {
		w = t
	}
	return io.Copy(w, os.Stdin)
}

func readFrom(names []string) (int64, error) {
	t, err := rp.NewReader(nil, names...)
	if err != nil {
		return 0, err
	}
	defer t.Close()
	return io.Copy(os.Stdout, t)
}
