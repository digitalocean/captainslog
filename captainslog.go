// Package captainslog provides a syslog parser and
// tools for working with syslog messages.
package captainslog

// MatchType represents types of matches that can
// be made against a SyslogMsg.
type MatchType int

const (
	// Program match is an exact match against syslog program name.
	Program MatchType = iota

	// Contains match checks if the syslog content contains a string.
	Contains
)

// Transformer accept a SyslogMsg, and return a modified SyslogMsg.
type Transformer interface {
	Transform(SyslogMsg) (SyslogMsg, error)
}

// Mutator accepts a pointer to a SyslogMsg and modifies it in place.
type Mutator interface {
	Mutate(*SyslogMsg) error
}

// Matcher accepts a SyslogMsg and returns true or false.
type Matcher interface {
	Match(msg *SyslogMsg) bool
}
