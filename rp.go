// Copyright (c) 2014 Brian Hetro <whee@smaertness.net>
// Use of this source code is governed by the ISC
// license which can be found in the LICENSE file.

//
// Package rp implements Reader and Writer interfaces to Redis Pub/Sub channels.
//
package rp

import redis "github.com/garyburd/redigo/redis"

// Writer implements writing to a Redis Pub/Sub channel.
type Writer struct {
	conn  redis.Conn
	names []string
}

// NewWriter returns a new Writer backed by the named Redis Pub/Sub
// channels. Writes are sent to each named channel.
// If unable to connect to Redis, the connection error is returned.
func NewWriter(names ...string) (*Writer, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}
	return &Writer{conn, names}, nil
}

// Close closes the Redis connection.
func (w *Writer) Close() error {
	return w.conn.Close()
}

// Write writes the contents of p to the Redis Pub/Sub channel via PUBLISH.
// It returns the number of bytes written.
// If the PUBLISH command fails, it returns why.
func (w *Writer) Write(p []byte) (n int, err error) {
	n = len(p)
	for _, c := range w.names {
		_, err = w.conn.Do("PUBLISH", c, p)
		if err != nil {
			return
		}
	}
	return
}

// Reader implements reading from a Redis Pub/Sub channel.
type Reader struct {
	conn  redis.Conn
	names []string
	psc   redis.PubSubConn
}

// NewReader returns a new Reader backed by the named Redis Pub/Sub
// channels. Reads may return content from any channel.
// If unable to connect to Redis or subscribe to the named channel,
// the error is returned.
func NewReader(names ...string) (r *Reader, err error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return
	}
	psc := redis.PubSubConn{conn}
	r = &Reader{conn, names, psc}
	for _, c := range names {
		err = psc.Subscribe(c)
		if err != nil {
			return
		}
	}
	return
}

// Read reads data from the Redis Pub/Sub channel into p.
// It returns the number of bytes read into p.
// If unable to receive from the Redis Pub/Sub channel, the error is returned.
func (r *Reader) Read(p []byte) (n int, err error) {
	switch no := r.psc.Receive().(type) {
	case redis.Message:
		n = len(no.Data)
		copy(p, no.Data)
	case redis.PMessage:
	case redis.Subscription:
	case error:
		err = no
		return
	}
	return
}

// Close unsubscribes from the Redis Pub/Sub channel and closes
// the Redis connection.
func (r *Reader) Close() error {
	err := r.psc.Unsubscribe()
	if err != nil {
		return err
	}
	return r.conn.Close()
}
