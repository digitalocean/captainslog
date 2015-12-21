package captainslog

// OutputAdapter is an interface for adapters that provide
// specific functionality to OutputChannelers. They
// are transport adapters - for instance, TCPOutputAdapter
// converts *Syslog messages received off a channeler
// to RFC3164 []byte encoded syslog messages sent over TCP.
type OutputAdapter interface {
	Output(s *SyslogMsg) (int, error)
	Connect() error
	RetryInterval() int
	Close()
}
