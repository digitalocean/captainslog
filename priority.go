package captainslog

import (
	"errors"
	"fmt"
)

var (
	//ErrBadPriority is returned when the priority of a message is malformed.
	ErrBadPriority = errors.New("Priority not found")

	//ErrBadFacility is returned when a facility is not within allowed values.
	ErrBadFacility = errors.New("Facility not found")

	//ErrBadSeverity is returned when a severity is not within allowed values.
	ErrBadSeverity = errors.New("Severity not found")
)

// Severity represents a syslog severity code
type Severity int

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

func (s Severity) String() string {
	var severityText string

	switch s {
	case Emerg:
		severityText = "emerg"
	case Alert:
		severityText = "alert"
	case Crit:
		severityText = "crit"
	case Err:
		severityText = "err"
	case Warning:
		severityText = "warning"
	case Notice:
		severityText = "notice"
	case Info:
		severityText = "info"
	case Debug:
		severityText = "debug"
	default:
	}

	return severityText
}

// Facility represents a syslog facility code
type Facility int

const (
	//Kern is the kernel rfc3164 facility.
	Kern Facility = 0

	//User is the user rfc3164 facility.
	User Facility = 1

	// Mail is the mail rfc3164 facility.
	Mail Facility = 2

	// Daemon is the daemon rfc3164 facility.
	Daemon Facility = 3

	// Auth is the auth rfc3164 facility.
	Auth Facility = 4

	// Syslog is the syslog rfc3164 facility.
	Syslog Facility = 5

	// LPR is the printer rfc3164 facility.
	LPR Facility = 6

	// News is a news rfc3164 facility.
	News Facility = 7

	// UUCP is the UUCP rfc3164 facility.
	UUCP Facility = 8

	// Cron is the cron rfc3164 facility.
	Cron Facility = 9

	//AuthPriv is the authpriv rfc3164 facility.
	AuthPriv Facility = 10

	// FTP is the ftp rfc3164 facility.
	FTP Facility = 11

	// Local0 is the local0 rfc3164 facility.
	Local0 Facility = 16

	// Local1 is the local1 rfc3164 facility.
	Local1 Facility = 17

	// Local2  is the local2 rfc3164 facility.
	Local2 Facility = 18

	// Local3 is the local3 rfc3164 facility.
	Local3 Facility = 19

	// Local4 is the local4 rfc3164 facility.
	Local4 Facility = 20

	// Local5 is the local5 rfc3164 facility.
	Local5 Facility = 21

	// Local6 is the local6 rfc3164 facility.
	Local6 Facility = 22

	// Local7 is the local7 rfc3164 facility.
	Local7 Facility = 23
)

func (f Facility) String() string {
	var faciliyText string

	switch f {
	case Kern:
		faciliyText = "kern"
	case User:
		faciliyText = "user"
	case Mail:
		faciliyText = "mail"
	case Daemon:
		faciliyText = "daemon"
	case Auth:
		faciliyText = "auth"
	case Syslog:
		faciliyText = "syslog"
	case LPR:
		faciliyText = "lpr"
	case News:
		faciliyText = "news"
	case UUCP:
		faciliyText = "uucp"
	case Cron:
		faciliyText = "cron"
	case AuthPriv:
		faciliyText = "authpriv"
	case FTP:
		faciliyText = "ftp"
	case Local0:
		faciliyText = "local0"
	case Local1:
		faciliyText = "local1"
	case Local2:
		faciliyText = "local2"
	case Local3:
		faciliyText = "local3"
	case Local4:
		faciliyText = "local4"
	case Local5:
		faciliyText = "local5"
	case Local6:
		faciliyText = "local6"
	case Local7:
		faciliyText = "local7"
	default:
	}

	return faciliyText
}

// Priority represents the PRI of a rfc3164 message.
type Priority struct {
	Priority int
	Facility Facility
	Severity Severity
}

// String converts the given Priority to a string.
func (p Priority) String() string {
	return fmt.Sprintf("%d", p.Priority)
}

func (p *Priority) SetFacility(f Facility) error {
	if int(f) < 0 || int(f) > 23 {
		return ErrBadFacility
	}

	p.Facility = f
	p.calculatePriority()

	return nil
}

func (p *Priority) SetSeverity(s Severity) error {
	if int(s) < 0 || int(s) > 7 {
		return ErrBadSeverity
	}

	p.Severity = s
	p.calculatePriority()

	return nil
}

func (p *Priority) calculatePriority() {
	p.Priority = int(p.Facility)*8 + int(p.Severity)
}

// NewPriority calculates a Priority from a Facility
// and Severity.
func NewPriority(f Facility, s Severity) (*Priority, error) {
	p := &Priority{}

	if err := p.SetFacility(f); err != nil {
		return nil, err
	}

	if err := p.SetSeverity(s); err != nil {
		return nil, err
	}

	p.calculatePriority()

	return p, nil
}
