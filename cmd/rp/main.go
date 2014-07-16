// Copyright (c) 2014 Brian Hetro <whee@smaertness.net>
// Use of this source code is governed by the ISC
// license which can be found in the LICENSE file.

package main

import (
	"io"
	"os"

	docopt "github.com/docopt/docopt-go"
	"github.com/whee/rp"
)

func main() {
	usage := `Redis Pipe.

Usage:
  rp -r <name>
  rp -w <name> [-p]

Options:
  -r, --read PIPE    Read from pipe named PIPE.
  -w, --write PIPE   Write to pipe named PIPE.
  -p, --passthrough  Pass written data to standard output.`

	arguments, _ := docopt.Parse(usage, nil, true, "Redis Pipe 0.1", false)
	if wp, ok := arguments["--write"]; ok && wp != nil {
		writeTo(wp.(string))
	} else if r, ok := arguments["--read"]; ok && r != nil {
		readFrom(r.(string))
	}
}

func writeTo(name string) (int64, error) {
	t, err := rp.NewWriter(name)
	if err != nil {
		return 0, err
	}
	defer t.Close()
	return io.Copy(t, os.Stdin)
}

func readFrom(name string) (int64, error) {
	t, err := rp.NewReader(name)
	if err != nil {
		return 0, err
	}
	defer t.Close()
	return io.Copy(os.Stdout, t)
}
