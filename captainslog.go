package captainslog

import "errors"

var (
	// ErrMutate is returned by a Mutator when it cannot
	// perform its function
	ErrMutate = errors.New("mutate error")
)

// Mutator accept a SyslogMsg, and return a modified SyslogMsg
type Mutator interface {
	Mutate(SyslogMsg) (SyslogMsg, error)
}

// Matcher accepts a SyslogMsg and returns true of it matches
type Matcher interface {
	Match(msg *SyslogMsg) bool
}

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
