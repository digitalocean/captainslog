package captainslog

import (
	"bufio"
	"net"
)

// TCPInputter receives RFC3164 bytes over TPC
// and emits *SyslogMsg struct pointers over a channel.
type TCPInputter struct {
	listener net.Listener
	address  string
	pipe     chan *SyslogMsg
}

// NewTCPInputter accepts a tcp address ("127.0.0.1:31337")
// and returns a new TCPInputter.
func NewTCPInputter(address string) (*TCPInputter, error) {
	i := &TCPInputter{
		address: address,
		pipe:    make(chan *SyslogMsg),
	}

	var err error
	i.listener, err = net.Listen("tcp", i.address)
	return i, err
}

// Listen starts listening at the address
func (i *TCPInputter) Listen() <-chan *SyslogMsg {
	go func() {
		for {
			conn, err := i.listener.Accept()
			if err != nil {
				continue
			}

			go func() {
				reader := bufio.NewReader(conn)
				for {
					b, err := reader.ReadBytes('\n')
					if err != nil {
						conn.Close()
						break
					}
					msg, err := NewSyslogMsgFromBytes(b)
					if err != nil {
						continue
					}
					i.pipe <- &msg
				}
			}()
		}
	}()

	return i.pipe
}

// Close closes the underlying listener.
func (i *TCPInputter) Close() {
	i.listener.Close()
}
