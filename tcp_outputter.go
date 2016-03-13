package captainslog

import "net"

// TCPOutputter sends *SyslogMsg as RFC3164
// encoded bytes over TCP to a destination
type TCPOutputter struct {
	conn          net.Conn
	address       string
	retryInterval int
}

// NewTCPOutputter accepts a tcp address ("127.0.0.1:31337")
// and a retry interval and returns a new running
// OutputChanneler
func NewTCPOutputter(address string, retry int) *TCPOutputter {
	return &TCPOutputter{
		address:       address,
		retryInterval: retry,
	}
}

// Connect tries to connect to the address
func (o *TCPOutputter) Connect() error {
	var err error
	o.conn, err = net.Dial("tcp", o.address)
	return err
}

// Output accepts a *SyslogMsg and sends an RFC3164
// []byte representation of it over TCP
func (o *TCPOutputter) Output(s *SyslogMsg) (int, error) {
	return o.conn.Write(s.Bytes())
}

// RetryInterval returns the retry interval of the Outputter
func (o *TCPOutputter) RetryInterval() int {
	return o.retryInterval
}

// Close closes the underlying connection
func (o *TCPOutputter) Close() {
	o.conn.Close()
}
