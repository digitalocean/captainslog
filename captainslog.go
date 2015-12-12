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
