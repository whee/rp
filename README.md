rp - Redis Pipe
===============

Install
-------

    go get github.com/whee/rp/cmd/rp

Usage
-----

    rp      -r <name>...
    rp [-p] -w <name>...

    Options:
      -r, --read <name>...  Read from the named channel.
      -w, --write <name>...  Write to the named channel.
      -p, --passthrough  Pass written data to standard output.
