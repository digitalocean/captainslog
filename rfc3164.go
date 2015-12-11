package captainslog

import (
	"errors"
	"fmt"
	"strconv"
	"time"
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

	priStart = '<'
	priEnd   = '>'
	priLen   = 5
)

// Facility represents a syslog facility code
type Facility int

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

// SyslogMsg holds an Unmarshaled rfc3164 message.
type SyslogMsg struct {
	Pri        Priority
	Time       time.Time
	Host       string
	Tag        string
	Cee        string
	IsCee      bool
	Content    string
	timeFormat string
}

// String returns the SyslogMsg as an RFC3164 string
func (s *SyslogMsg) String() string {
	return fmt.Sprintf("<%d>%s %s %s%s%s\n", s.Pri.Priority, s.Time.Format(s.timeFormat), s.Host, s.Tag, s.Cee, s.Content)
}

// Bytes returns the SyslogMsg as RFC3164 []byte
func (s *SyslogMsg) Bytes() []byte {
	return []byte(s.String())
}

type parser struct {
	tokenStart int
	tokenEnd   int
	buf        []byte
	bufLen     int
	bufEnd     int
	cur        int
	msg        *SyslogMsg
}

// Unmarshal accepts a byte array containing an rfc3164 message
// and a pointer to a SyslogMsg struct, and attempts to parse
// the message and fill in the struct.
func Unmarshal(b []byte, msg *SyslogMsg) error {

	p := &parser{
		buf:    b,
		bufLen: len(b),
		bufEnd: len(b) - 1,
		cur:    0,
		msg:    msg,
	}

	err := p.parse()
	return err
}

func (p *parser) parse() error {
	err := p.parsePri()
	if err != nil {
		return err
	}
	err = p.parseTime()
	if err != nil {
		return err
	}

	err = p.parseHost()
	if err != nil {
		return err
	}

	err = p.parseTag()
	if err != nil {
		return err
	}

	p.parseCee()

	err = p.parseContent()
	return err
}

func isNum(c byte) bool {
	return c >= '0' && c <= '9'
}

func (p *parser) parsePri() error {
	if p.bufLen == 0 || (p.cur+priLen) > p.bufEnd {
		return ErrBadPriority
	}

	if p.buf[p.cur] != priStart {
		return ErrBadPriority
	}

	p.cur++
	p.tokenStart = p.cur

	for p.buf[p.cur] != priEnd {
		if !isNum(p.buf[p.cur]) {
			return ErrBadPriority
		}

		p.cur++

		if p.cur > p.bufEnd {
			return ErrBadPriority
		}

		if p.cur > (priLen - 1) {
			return ErrBadPriority
		}
	}

	p.tokenEnd = p.cur
	pVal, err := strconv.Atoi(string(p.buf[p.tokenStart:p.tokenEnd]))
	if err != nil {
		return err
	}

	p.msg.Pri = Priority{
		Priority: pVal,
		Facility: Facility(pVal / 8),
		Severity: Severity(pVal % 8),
	}

	p.cur++
	return err
}

func (p *parser) parseTime() error {
	var err error
	var foundTime bool

	p.tokenStart = p.cur
	for _, timeFormat := range timeFormats {
		tLen := len(timeFormat)
		if p.cur+tLen > p.bufEnd {
			continue
		}

		timeStr := string(p.buf[p.cur : p.cur+tLen])
		p.msg.Time, err = time.Parse(timeFormat, timeStr)
		if err == nil {
			p.cur = p.cur + tLen
			p.tokenEnd = p.cur
			p.msg.timeFormat = timeFormat
			foundTime = true
			break
		}
	}
	if !foundTime {
		err = ErrBadTime
	}
	return err
}

func (p *parser) parseHost() error {
	var err error
	for p.buf[p.cur] == ' ' {
		p.cur++
		if p.cur > p.bufEnd {
			return ErrBadHost
		}
	}

	p.tokenStart = p.cur

	for p.buf[p.cur] != ' ' {
		p.cur++
		if p.cur > p.bufEnd {
			return ErrBadHost
		}
	}

	p.tokenEnd = p.cur
	p.msg.Host = string(p.buf[p.tokenStart:p.tokenEnd])
	return err
}

func (p *parser) parseTag() error {
	var err error

	for p.buf[p.cur] == ' ' {
		p.cur++
		if p.cur > p.bufEnd {
			return ErrBadTag
		}
	}

	p.tokenStart = p.cur

	for p.buf[p.cur] != ':' && p.buf[p.cur] != ' ' {
		p.cur++
		if p.cur > p.bufEnd {
			return ErrBadTag
		}
	}

	if p.buf[p.cur] == ':' {
		p.cur++
	}

	p.tokenEnd = p.cur
	p.msg.Tag = string(p.buf[p.tokenStart:p.tokenEnd])
	return err
}

func (p *parser) parseCee() {
	p.tokenStart = p.cur
	cur := p.cur

	for p.buf[cur] == ' ' {
		cur++
		if cur > p.bufEnd {
			return
		}
	}

	if cur+4 > p.bufEnd {
		return
	}

	if p.buf[cur] != '@' {
		return
	}

	cur++
	if p.buf[cur] != 'c' {
		return
	}

	cur++
	if p.buf[cur] != 'e' {
		return
	}

	cur++
	if p.buf[cur] != 'e' {
		return
	}

	cur++
	if p.buf[cur] != ':' {
		return
	}

	cur++
	p.cur = cur

	p.tokenEnd = cur
	p.msg.IsCee = true
	p.msg.Cee = string(p.buf[p.tokenStart:p.tokenEnd])

	return
}

func (p *parser) parseContent() error {
	var err error

	p.tokenStart = p.cur

	for p.buf[p.cur] != '\n' {
		p.cur++
		if p.cur > p.bufEnd {
			return ErrBadContent
		}
	}

	p.tokenEnd = p.cur
	p.msg.Content = string(p.buf[p.tokenStart:p.tokenEnd])
	return err
}
