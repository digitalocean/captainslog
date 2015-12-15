# captainslog [![Build Status](https://travis-ci.org/digitalocean/captainslog.svg?branch=master)](https://travis-ci.org/digitalocean/captainslog) [![Doc Status](https://godoc.org/github.com/digitalocean/captainslog?status.png)](https://godoc.org/github.com/digitalocean/captainslog)

--
    import "github.com/digitalocean/captainslog"


## Usage

```go
var (
	//ErrBadPriority is returned when the priority of a message is malformed
	ErrBadPriority = errors.New("Priority not found")

	//ErrBadFacility is returned when a facility is not within allowed values
	ErrBadFacility = errors.New("Facility not found")

	//ErrBadSeverity is returned when a severity is not within allowed values
	ErrBadSeverity = errors.New("Severity not found")

	//ErrBadTime is returned when the time of a message is malformed
	ErrBadTime = errors.New("Time not found")

	//ErrBadHost is returned when the host of a message is malformed
	ErrBadHost = errors.New("Host not found")

	//ErrBadTag is returned when the tag of a message is malformed
	ErrBadTag = errors.New("Tag not found")

	//ErrBadContent is returned when the content of a message is malformed
	ErrBadContent = errors.New("Content not found")
)
```

```go
var (
	// ErrMutate is returned by a Mutator when it cannot
	// perform its function
	ErrMutate = errors.New("mutate error")
)
```

#### func  Unmarshal

```go
func Unmarshal(b []byte, msg *SyslogMsg) error
```
Unmarshal accepts a byte array containing an rfc3164 message and a pointer to a
SyslogMsg struct, and attempts to parse the message and fill in the struct.

#### type Facility

```go
type Facility int
```

Facility represents a syslog facility code

```go
const (
	//Kern is the kernel rfc3164 facility
	Kern Facility = 0

	//User is the user rfc3164 facility
	User Facility = 1

	// Mail is the mail rfc3164 facility
	Mail Facility = 2

	// Daemon is the daemon rfc3164 facility
	Daemon Facility = 3

	// Auth is the auth rfc3164 facility
	Auth Facility = 4

	// Syslog is the syslog rfc3164 facility
	Syslog Facility = 5

	// LPR is the printer rfc3164 facility
	LPR Facility = 6

	// News is a news rfc3164 facility
	News Facility = 7

	// UUCP is the UUCP rfc3164 facility
	UUCP Facility = 8

	// Cron is the cron rfc3164 facility
	Cron Facility = 9

	//AuthPriv is the authpriv rfc3164 facility
	AuthPriv Facility = 10

	// FTP is the ftp rfc3164 facility
	FTP Facility = 11

	// Local0 is the local0 rfc3164 facility
	Local0 Facility = 16

	// Local1 is the local1 rfc3164 facility
	Local1 Facility = 17

	// Local2  is the local2 rfc3164 facility
	Local2 Facility = 18

	// Local3 is the local3 rfc3164 facility
	Local3 Facility = 19

	// Local4 is the local4 rfc3164 facility
	Local4 Facility = 20

	// Local5 is the local5 rfc3164 facility
	Local5 Facility = 21

	// Local6 is the local6 rfc3164 facility
	Local6 Facility = 22

	// Local7 is the local7 rfc3164 facility
	Local7 Facility = 23
)
```

#### func  FacilityTextToFacility

```go
func FacilityTextToFacility(facilityText string) (Facility, error)
```
FacilityTextToFacility accepts a string representation of a syslog facility and
returns a captainslog.Facility

#### type Fields

```go
type Fields map[string]interface{}
```

Fields are a map of key value pairs for a log line that will be output as JSON

#### type JSONKeyMutator

```go
type JSONKeyMutator struct {
}
```

JSONKeyMutator is a Mutator implementation that finds periods in JSON keys in
CEE syslog messages and replaces them. This can be used in conjunction with
systems such as Elasticsearch 2.x which do not fully support ECMA-404 (for
instance, Elasticsearch 2.x does not allow periods in key names, which ECMA-404
does)

#### func  NewJSONKeyMutator

```go
func NewJSONKeyMutator(replacer *strings.Replacer) *JSONKeyMutator
```
NewJSONKeyMutator applies a strings.Replacer to all keys in a JSON document in a
CEE syslog message.

#### func (*JSONKeyMutator) Mutate

```go
func (m *JSONKeyMutator) Mutate(msg SyslogMsg) (SyslogMsg, error)
```
Mutate accepts a SyslogMsg, and if it is a CEE syslog message, "fixes" the JSON
keys to be compatible with Elasticsearch 2.x

#### type MostlyFeaturelessLogger

```go
type MostlyFeaturelessLogger struct {
}
```

MostlyFeaturelessLogger is a mostly featureless logger created for simple
structured logging of Notice and Err messages from daemons created with
captainslog to syslog. If you need something more than that this probably is not
something that will make you happy.

#### func  NewMostlyFeaturelessLogger

```go
func NewMostlyFeaturelessLogger(f Facility) (*MostlyFeaturelessLogger, error)
```
NewMostlyFeaturelessLogger returns a new MostlyFeaturelessLogger for the given
Facility

#### func (*MostlyFeaturelessLogger) ErrorWithFields

```go
func (l *MostlyFeaturelessLogger) ErrorWithFields(fields Fields) error
```
ErrorWithFields accepts Fields and logs a @cee structured log to syslog at level
Err

#### func (*MostlyFeaturelessLogger) InfoWithFields

```go
func (l *MostlyFeaturelessLogger) InfoWithFields(fields Fields) error
```
InfoWithFields accepts Fields and logs a @cee structured log to syslog at level
Notice

#### type Mutator

```go
type Mutator interface {
	Mutate(SyslogMsg) (SyslogMsg, error)
}
```

Mutator accept a SyslogMsg, and return a modified SyslogMsg

#### type Priority

```go
type Priority struct {
	Priority int
	Facility Facility
	Severity Severity
}
```

Priority represents the PRI of a rfc3164 message.

#### func  NewPriority

```go
func NewPriority(f Facility, s Severity) (*Priority, error)
```
NewPriority calculates a Priority from a Facility and Severity

#### type Queue

```go
type Queue struct {
}
```


#### func  NewQueue

```go
func NewQueue(max int) *Queue
```

#### func (*Queue) Dequeue

```go
func (q *Queue) Dequeue() SyslogMsg
```

#### func (*Queue) Enqueue

```go
func (q *Queue) Enqueue(s SyslogMsg) error
```

#### func (*Queue) Len

```go
func (q *Queue) Len() int
```

#### type Severity

```go
type Severity int
```

Severity represents a syslog severity code

```go
const (
	// Emerg is an emergency rfc3164 severity
	Emerg Severity = 0

	// Alert is an alert rfc3164 severity
	Alert Severity = 1

	// Crit is a critical level rfc3164 severity
	Crit Severity = 2

	// Err is an error level rfc3164 severity
	Err Severity = 3

	// Warning is a warning level rfc3164 severity
	Warning Severity = 4

	// Notice is a notice level rfc3164 severity
	Notice Severity = 5

	// Info is an info level rfc3164 severity
	Info Severity = 6

	// Debug is a debug level rfc3164 severity
	Debug Severity = 7
)
```

#### type SyslogMsg

```go
type SyslogMsg struct {
	Pri     Priority
	Time    time.Time
	Host    string
	Tag     string
	Cee     string
	IsCee   bool
	Content string
}
```

SyslogMsg holds an Unmarshaled rfc3164 message.

#### func (*SyslogMsg) Bytes

```go
func (s *SyslogMsg) Bytes() []byte
```
Bytes returns the SyslogMsg as RFC3164 []byte

#### func (*SyslogMsg) String

```go
func (s *SyslogMsg) String() string
```
String returns the SyslogMsg as an RFC3164 string
