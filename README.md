rp - Redis Pipe
===============

Install
-------

    go get github.com/whee/rp/cmd/rp

Usage
-----

    rp      read  <name>...
    rp [-p] write <name>...
    rp      -n <network> -a <address> read  <name>...
    rp [-p] -n <network> -a <address> write <name>...

    Options:
      -p, --passthrough  Pass written data to standard output.
      -n, --network <network>  Network of the Redis server [default: tcp]
      -a, --address <address>  Address of the Redis server [default: :6379]
