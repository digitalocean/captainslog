package captainslog

import "net"

// OutputGadgetTCP is a gadget for
// sending syslog messages over TCP
type OutputGadgetTCP struct {
	conn          net.Conn
	address       string
	retryInterval int
}

// NewOutputGadgetTCP accepts an address and a retry interval
// and returns a new running gadget
func NewOutputGadgetTCP(address string, retry int) *OutputGadgetTCP {
	return &OutputGadgetTCP{
		address:       address,
		retryInterval: retry,
	}
}

// Connect tries to connect the gadget to an address
func (o *OutputGadgetTCP) Connect() error {
	var err error
	o.conn, err = net.Dial("tcp", o.address)
	return err
}

// Output accepts a SyslogMsg pointer and sends an RFC3164
// []byte representation of it over TCP
func (o *OutputGadgetTCP) Output(s *SyslogMsg) (int, error) {
	return o.conn.Write(s.Bytes())
}

// RetryInterval returns the retry interval of the gadget
func (o *OutputGadgetTCP) RetryInterval() int {
	return o.retryInterval
}

// Close closes the underlying connection in the gadget
func (o *OutputGadgetTCP) Close() {
	o.conn.Close()
}
