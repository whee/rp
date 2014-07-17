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

type config struct {
	channels    []string
	read        bool
	write       bool
	passthrough bool
	rpConfig    rp.Config
}

func main() {
	usage := `Redis Pipe.

Usage:
	rp      read  <name>...
	rp [-p] write <name>...
	rp      -n <network> -a <address> read  <name>...
	rp [-p] -n <network> -a <address> write <name>...

Options:
  -p, --passthrough  Pass written data to standard output.
  -n, --network <network>  Network of the Redis server [default: tcp]
  -a, --address <address>  Address of the Redis server [default: :6379]`

	arguments, err := docopt.Parse(usage, nil, true, "Redis Pipe 0.9", false)
	if err != nil {
		panic(err)
	}
	conf := config{
		channels:    arguments["<name>"].([]string),
		read:        arguments["read"].(bool),
		write:       arguments["write"].(bool),
		passthrough: arguments["--passthrough"].(bool),
		rpConfig: rp.Config{
			Network: arguments["--network"].(string),
			Address: arguments["--address"].(string)},
	}
	_, err = conf.do()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func (c *config) do() (n int64, err error) {
	if c.write {
		n, err = writeTo(c)
	} else if c.read {
		n, err = readFrom(c)
	}
	return
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

func writeTo(c *config) (int64, error) {
	t, err := rp.NewWriter(&c.rpConfig, c.channels...)
	if err != nil {
		return 0, err
	}
	defer t.Close()

	var w io.Writer
	if c.passthrough {
		w = ptWriter{t}
	} else {
		w = t
	}
	return io.Copy(w, os.Stdin)
}

func readFrom(c *config) (int64, error) {
	t, err := rp.NewReader(&c.rpConfig, c.channels...)
	if err != nil {
		return 0, err
	}
	defer t.Close()
	return io.Copy(os.Stdout, t)
}
