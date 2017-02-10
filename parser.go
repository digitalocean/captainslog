package captainslog

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	priStart     = '<'
	priEnd       = '>'
	priLen       = 5
	dateStampLen = 10
)

var (
	//ErrBadTime is returned when the time of a message is malformed.
	ErrBadTime = errors.New("Time not found")

	//ErrBadHost is returned when the host of a message is malformed.
	ErrBadHost = errors.New("Host not found")

	//ErrBadTag is returned when the tag of a message is malformed.
	ErrBadTag = errors.New("Tag not found")

	//ErrBadContent is returned when the content of a message is malformed.
	ErrBadContent = errors.New("Content not found")

	dateStampFormat   = "2006-01-02"
	rsyslogTimeFormat = "2006-01-02T15:04:05.999999-07:00"

	timeFormats = []string{
		"Mon Jan _2 15:04:05 MST 2006",
		"Mon Jan _2 15:04:05 2006",
		"Mon Jan _2 15:04:05",
		"Jan _2 15:04:05",
		"Jan 02 15:04:05",
	}
)

// Parser is a parser for syslog messages.
type Parser struct {
	tokenStart          int
	tokenEnd            int
	buf                 []byte
	bufLen              int
	bufEnd              int
	cur                 int
	requireTerminator   bool
	optionNoHostname    bool
	optionDontParseJSON bool
	msg                 *SyslogMsg
}

// NewParser returns a new parser
func NewParser(options ...func(*Parser)) *Parser {
	p := Parser{}
	for _, option := range options {
		option(&p)
	}
	return &p
}

// OptionNoHostname sets the parser to not expect the hostname
// as part of the syslog message, and instead ask the host
// for its hostname.
func OptionNoHostname(p *Parser) {
	p.optionNoHostname = true
}

// OptionDontParseJSON sets the parser to not parse JSON in
// the content field of the message. A subsequent call to SyslogMsg.String()
// or SyslogMsg.Bytes() will then use SyslogMsg.Content for the content field,
// unless SyslogMsg.JSONValues have been added since the message was
// originally parsed. If SyslogMsg.JSONValues have been added, the call to
// SyslogMsg.String() or SyslogMsg.Bytes() will then parse the JSON, and
// merge the results with the keys in SyslogMsg.JSONVaues.
func OptionDontParseJSON(p *Parser) {
	p.optionDontParseJSON = true
}

// ParseBytes accepts a []byte and tries to parse it into a SyslogMsg
func (p *Parser) ParseBytes(b []byte) (SyslogMsg, error) {
	p.buf = b
	p.bufLen = len(b)
	p.bufEnd = len(b) - 1
	p.cur = 0
	msg := NewSyslogMsg()
	msg.optionDontParseJSON = p.optionDontParseJSON
	p.msg = &msg

	err := p.parse()
	if p.msg.Time.Year() == 0 {
		p.msg.Time = p.msg.Time.AddDate(time.Now().Year(), 0, 0)
	}
	return *p.msg, err
}

func (p *Parser) parse() error {
	err := p.parsePri()
	if err != nil {
		return err
	}
	err = p.parseTime()
	if err != nil {
		return err
	}

	if p.optionNoHostname {
		host, err := os.Hostname()
		if err != nil {
			return ErrBadHost
		}
		p.msg.Host = host
	} else {
		err = p.parseHost()
		if err != nil {
			return err
		}
	}

	err = p.parseTag()
	if err != nil {
		return err
	}

	err = p.parseCee()
	if err != nil {
		return err
	}

	err = p.parseContent()
	return err
}

func (p *Parser) parsePri() error {
	var err error

	if p.bufLen == 0 || (p.cur+priLen) > p.bufEnd {
		return ErrBadPriority
	}

	if p.buf[p.cur] != priStart {
		return ErrBadPriority
	}

	p.cur++
	p.tokenStart = p.cur

	if p.buf[p.cur] == priEnd {
		return ErrBadPriority
	}

	for p.buf[p.cur] != priEnd {
		if !(p.buf[p.cur] >= '0' && p.buf[p.cur] <= '9') {
			return ErrBadPriority
		}

		p.cur++

		if p.cur > (priLen - 1) {
			return ErrBadPriority
		}
	}

	p.tokenEnd = p.cur
	pVal, _ := strconv.Atoi(string(p.buf[p.tokenStart:p.tokenEnd]))

	p.msg.Pri = Priority{
		Priority: pVal,
		Facility: Facility(pVal / 8),
		Severity: Severity(pVal % 8),
	}

	p.cur++
	return err
}

// checkForDateTime checks for a YYYY-MM-DD date. This routine could be modified
// to be more efficient, as constructing a time.Time is not necessary. For now
// this is a simple way to avoid the complexity of checking for a valid
// datetime, which is more complicated than it might first appear to be.
func checkForDateTime(timeStr string) bool {
	_, err := time.Parse(dateStampFormat, timeStr)
	if err != nil {
		return false
	}
	return true
}

func (p *Parser) parseTime() error {
	var err error
	var foundTime bool

	p.tokenStart = p.cur

	// no timestamp format is shorter than YYYY-MM-DD, so if buffer is shorter
	// than this it is safe to assume we don't have a valid datetime.
	if p.cur+dateStampLen > p.bufEnd {
		return ErrBadTime
	}

	if checkForDateTime(string(p.buf[p.cur : p.cur+dateStampLen])) {
		tokenStart := p.cur
		tokenEnd := p.cur

		for p.buf[tokenEnd] != ' ' {
			tokenEnd++
			if tokenEnd > p.bufEnd {
				return ErrBadTime
			}
		}

		timeStr := string(p.buf[tokenStart:tokenEnd])
		p.msg.Time, err = time.Parse(rsyslogTimeFormat, timeStr)
		if err == nil {
			p.cur = tokenEnd
			p.tokenEnd = p.cur
			p.msg.timeFormat = rsyslogTimeFormat
		}
		return err
	}

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

func (p *Parser) parseHost() error {
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

func (p *Parser) parseTag() error {
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

func (p *Parser) parseCee() error {
	if p.cur >= len(p.buf)-1 {
		return ErrBadContent
	}

	p.tokenStart = p.cur
	cur := p.cur

	for p.buf[cur] == ' ' {
		cur++
		if cur >= len(p.buf)-1 {
			return nil
		}
	}

	if cur+4 > p.bufEnd {
		return nil
	}

	if p.buf[cur] != '@' {
		return nil
	}

	cur++
	if p.buf[cur] != 'c' {
		return nil
	}

	cur++
	if p.buf[cur] != 'e' {
		return nil
	}

	cur++
	if p.buf[cur] != 'e' {
		return nil
	}

	cur++
	if p.buf[cur] != ':' {
		return nil
	}

	cur++
	p.cur = cur

	p.tokenEnd = cur
	p.msg.IsCee = true
	p.msg.Cee = string(p.buf[p.tokenStart:p.tokenEnd])

	return nil
}

func (p *Parser) parseContent() error {
	if p.cur >= len(p.buf)-1 {
		return ErrBadContent
	}

	var err error
	p.tokenStart = p.cur

	for p.buf[p.cur] != '\n' {
		p.cur++
		if p.cur > p.bufEnd {
			if p.requireTerminator {
				return ErrBadContent
			}
			goto exitContentSearch
		}
	}
exitContentSearch:
	p.tokenEnd = p.cur

	if p.msg.IsCee && !p.optionDontParseJSON {
		decoder := json.NewDecoder(bytes.NewBuffer(p.buf[p.tokenStart:p.tokenEnd]))
		decoder.UseNumber()
		err = decoder.Decode(&p.msg.JSONValues)
		if err != nil {
			p.msg.IsCee = false
		}
	}
	p.msg.Content = string(p.buf[p.tokenStart:p.tokenEnd])
	return err
}
