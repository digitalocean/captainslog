package captainslog

import "errors"

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

	priStart = '<'
	priEnd   = '>'
	priLen   = 5
)

// Facility represents a syslog facility code
type Facility int

//FacilityTextToFacility accepts a string representation of a syslog
// facility and returns a captainslog.Facility
func FacilityTextToFacility(facilityText string) (Facility, error) {
	var err error
	var facility Facility

	switch facilityText {
	case "KERN":
		facility = Kern
	case "USER":
		facility = User
	case "MAIL":
		facility = Mail
	case "DAEMON":
		facility = Daemon
	case "AUTH":
		facility = Auth
	case "SYSLOG":
		facility = Syslog
	case "LPR":
		facility = LPR
	case "NEWS":
		facility = News
	case "UUCP":
		facility = UUCP
	case "CRON":
		facility = Cron
	case "AUTHPRIV":
		facility = AuthPriv
	case "FTP":
		facility = FTP
	case "LOCAL0":
		facility = Local0
	case "LOCAL1":
		facility = Local1
	case "LOCAL2":
		facility = Local2
	case "LOCAL3":
		facility = Local3
	case "LOCAL4":
		facility = Local4
	case "LOCAL5":
		facility = Local5
	case "LOCAL6":
		facility = Local6
	case "LOCAL7":
		facility = Local7
	default:
		facility = Facility(-1)
		err = ErrBadFacility
	}
	return facility, err
}

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

	timeFormats = []string{
		"2006-01-02T15:04:05.999999-07:00",
		"2006-01-02T15:04:05.999-07:00",
		"2006-01-02T15:04:05-07:00",
		"Mon Jan _2 15:04:05 MST 2006",
		"Mon Jan _2 15:04:05 2006",
		"Mon Jan _2 15:04:05",
	}
)

// Priority represents the PRI of a rfc3164 message.
type Priority struct {
	Priority int
	Facility Facility
	Severity Severity
}

// NewPriority calculates a Priority from a Facility
// and Severity
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
