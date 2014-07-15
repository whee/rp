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
		t, err := rp.NewWriter(wp.(string))
		if err != nil {
			panic(err)
		}
		defer t.Close()
		io.Copy(t, os.Stdin)
	} else if r, ok := arguments["--read"]; ok && r != nil {
		t, err := rp.NewReader(r.(string))
		if err != nil {
			panic(err)
		}
		defer t.Close()
		io.Copy(os.Stdout, t)
	}
}
