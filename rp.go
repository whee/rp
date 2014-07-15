package rp

import redis "github.com/garyburd/redigo/redis"

type Writer struct {
	conn redis.Conn
	name string
}

func NewWriter(name string) (*Writer, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}
	return &Writer{conn, name}, nil
}

func (w *Writer) Close() error {
	return w.conn.Close()
}

func (w *Writer) Write(p []byte) (n int, err error) {
	n = len(p)
	_, err = w.conn.Do("PUBLISH", w.name, p)
	return
}

type Reader struct {
	conn redis.Conn
	name string
	psc  redis.PubSubConn
}

func NewReader(name string) (*Reader, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}
	psc := redis.PubSubConn{conn}
	err = psc.Subscribe(name)
	if err != nil {
		return nil, err
	}
	return &Reader{conn, name, psc}, nil
}

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
	return n, nil
}

func (r *Reader) Close() error {
	err := r.psc.Unsubscribe()
	if err != nil {
		return err
	}
	return r.conn.Close()
}
