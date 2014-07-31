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
      
Basics
------

rp converts standard input and output to networked pipes via Redis Pub/Sub.
This allows trivial construction of MIMO data flows for processing and analysis.

For example:

    Terminal 1> $ rp read foo | jq --unbuffered '.c' | rp write cs
    Terminal 2> $ rp read foo | jq '.'
    Terminal 3> $ rp read cs
    Terminal 4> $ echo '{"a":1,"b":2,"c":"foo"}' | rp write foo
    Terminal 2> {
    Terminal 2>     "a": 1,
    Terminal 2>     "b": 2,
    Terminal 2>     "c": "foo"
    Terminal 2> }
    Terminal 3> "foo"
