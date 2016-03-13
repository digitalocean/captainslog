package captainslog

import (
	"bufio"
	"bytes"
	"net"
	"sync"
	"testing"
)

func TestOutputChannelerWithTCPOutputter(t *testing.T) {
	testMsg := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	var msg SyslogMsg
	err := Unmarshal(testMsg, &msg)
	if err != nil {
		t.Error(err)
	}

	address := "127.0.0.1:45454"
	retryInterval := 5

	l, err := net.Listen("tcp", address)
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Error(err)
		}

		reader := bufio.NewReader(conn)
		line, _, err := reader.ReadLine()
		if err != nil {
			t.Error(err)
		}

		line = append(line, '\n')
		if want, got := 0, bytes.Compare(testMsg, line); want != got {
			t.Errorf("want '%v', got '%v'", want, got)
			t.Errorf("%s", line)
		}

		wg.Done()
	}()

	a := NewTCPOutputter(address, retryInterval)
	o := NewOutputChanneler(a)

	o.OutChan <- &msg
	o.CmdChan <- CmdStop
	wg.Wait()
}
