package captainslog

import "errors"

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
	var severity_text string

	switch s {
	case Emerg:
		severity_text = "emerg"
	case Alert:
		severity_text = "alert"
	case Crit:
		severity_text = "crit"
	case Err:
		severity_text = "err"
	case Warning:
		severity_text = "warning"
	case Notice:
		severity_text = "notice"
	case Info:
		severity_text = "info"
	case Debug:
		severity_text = "debug"
	default:
	}

	return severity_text
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
	var facility_text string

	switch f {
	case Kern:
		facility_text = "kern"
	case User:
		facility_text = "user"
	case Mail:
		facility_text = "mail"
	case Daemon:
		facility_text = "daemon"
	case Auth:
		facility_text = "auth"
	case Syslog:
		facility_text = "syslog"
	case LPR:
		facility_text = "lpr"
	case News:
		facility_text = "news"
	case UUCP:
		facility_text = "uucp"
	case Cron:
		facility_text = "cron"
	case AuthPriv:
		facility_text = "authpriv"
	case FTP:
		facility_text = "ftp"
	case Local0:
		facility_text = "local0"
	case Local1:
		facility_text = "local1"
	case Local2:
		facility_text = "local2"
	case Local3:
		facility_text = "local3"
	case Local4:
		facility_text = "local4"
	case Local5:
		facility_text = "local5"
	case Local6:
		facility_text = "local6"
	case Local7:
		facility_text = "local7"
	default:
	}

	return facility_text
}

// Priority represents the PRI of a rfc3164 message.
type Priority struct {
	Priority int
	Facility Facility
	Severity Severity
}

// NewPriority calculates a Priority from a Facility
// and Severity.
func NewPriority(f Facility, s Severity) (*Priority, error) {
	p := &Priority{
		Facility: f,
		Severity: s,
		Priority: (int(f) * 8) + int(s),
	}

	var err error

	if int(f) < 0 || int(f) > 23 {
		return p, ErrBadFacility
	}

	if int(s) < 0 || int(s) > 7 {
		return p, ErrBadSeverity
	}

	return p, err
}
