package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrTimeout = errors.New("connection timeout")

type TelnetClient interface {
	// Place your code here
	Connect() error
	Close() error
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("%v, %w", err, ErrTimeout)
	}
	t.conn = conn
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn == nil {
		return nil
	}
	return t.conn.Close()
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return nil
	}
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return nil
	}
	_, err := io.Copy(t.out, t.conn)
	return err
}
