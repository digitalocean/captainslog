package captainslog

import "net"

// TCPOutputAdapter sends *SyslogMsg as RFC3164
// encoded bytes over TCP to a destination
type TCPOutputAdapter struct {
	conn          net.Conn
	address       string
	retryInterval int
}

// NewTCPOutputAdapter accepts a tcp address ("127.0.0.1:31337")
// and a retry interval and returns a new running
// OutputChanneler
func NewTCPOutputAdapter(address string, retry int) *TCPOutputAdapter {
	return &TCPOutputAdapter{
		address:       address,
		retryInterval: retry,
	}
}

// Connect tries to connect to the address
func (o *TCPOutputAdapter) Connect() error {
	var err error
	o.conn, err = net.Dial("tcp", o.address)
	return err
}

// Output accepts a *SyslogMsg and sends an RFC3164
// []byte representation of it over TCP
func (o *TCPOutputAdapter) Output(s *SyslogMsg) (int, error) {
	return o.conn.Write(s.Bytes())
}

// RetryInterval returns the retry interval of the OutputAdapter
func (o *TCPOutputAdapter) RetryInterval() int {
	return o.retryInterval
}

// Close closes the underlying connection
func (o *TCPOutputAdapter) Close() {
	o.conn.Close()
}
