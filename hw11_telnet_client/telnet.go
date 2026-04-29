package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (tc *telnetClient) Connect() error {
	var err error
	tc.conn, err = net.DialTimeout("tcp", tc.address, tc.timeout)
	return err
}

func (tc *telnetClient) Close() error {
	if tc.conn == nil {
		return nil
	}

	return tc.conn.Close()
}

func (tc *telnetClient) Send() error {
	scanner := bufio.NewScanner(tc.in)

	for scanner.Scan() {
		line := scanner.Bytes()
		_, err := tc.conn.Write(append(line, '\n'))
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (tc *telnetClient) Receive() error {
	r := bufio.NewReader(tc.conn)

	for {
		b, err := r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		_, err = tc.out.Write([]byte{b})
		if err != nil {
			return err
		}
	}
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
