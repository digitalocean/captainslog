package captainslog

import "errors"

var (
	// ErrTransform is returned by a Transformer when it cannot
	// perform its function
	ErrTransform = errors.New("transform error")
)

// Transformer accept a SyslogMsg, and return a modified SyslogMsg
type Transformer interface {
	Transform(SyslogMsg) (SyslogMsg, error)
}

// Mutator accepts a pointer to a SyslogMsg and modifies it in place
type Mutator interface {
	Mutate(*SyslogMsg) error
}

// Matcher accepts a SyslogMsg and returns true if it matches
type Matcher interface {
	Match(msg *SyslogMsg) bool
}

// Outputter is an interface that provides
// specific functionality to OutputChannelers. They
// are transport adapters - for instance, TCPOutputter
// converts *Syslog messages received off a channeler
// to RFC3164 []byte encoded syslog messages sent over TCP.
type Outputter interface {
	Output(s *SyslogMsg) (int, error)
	Connect() error
	RetryInterval() int
	Close()
}
