package captainslog

import (
	"bufio"
	"bytes"
	"net"
	"sync"
	"testing"
)

func TestTCPOutputter(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	msg, err := NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	address := "127.0.0.1:33333"
	retryInterval := 1

	l, err := net.Listen("tcp", address)
	if err != nil {
		t.Error(err)
	}
	defer l.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()

		reader := bufio.NewReader(conn)
		line, _, err := reader.ReadLine()
		if err != nil {
			t.Error(err)
		}

		line = append(line, '\n')
		if want, got := 0, bytes.Compare(b, line); want != got {
			t.Errorf("want '%v', got '%v'", want, got)
			t.Errorf("%s", line)
		}

		wg.Done()
	}()

	a := NewTCPOutputter(address, retryInterval)
	defer a.Close()

	err = a.Connect()
	if err != nil {
		t.Error(err)
	}

	n, err := a.Output(&msg)
	if err != nil {
		t.Error(err)
	}

	if want, got := len(b), n; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}
