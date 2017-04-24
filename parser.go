package captainslog

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"
	"unicode"
)

const (
	priStart     = '<'
	priEnd       = '>'
	priLen       = 5
	dateStampLen = 10

	// the following are used for checking for likely YYYY-MM-DD datestamps
	// in checkForLikelyDateTime
	yearLen     = 4
	startMonth  = 5
	monthLen    = 2
	dayLen      = 2
	startDay    = 8
	datePartSep = "-"
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
	var err error
	p.msg.Pri, p.cur, err = ParsePri(p.cur, p.buf)
	if err != nil {
		return err
	}

	p.msg.Time, p.msg.timeFormat, p.cur, err = ParseTime(p.cur, p.buf)
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
		p.msg.Host, p.cur, err = ParseHost(p.cur, p.buf)
		if err != nil {
			return err
		}
	}

	p.msg.Tag, p.msg.Program, p.msg.Pid, p.cur, err = ParseTag(p.cur, p.buf)
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

func ParsePri(cur int, buf []byte) (Priority, int, error) {
	var err error
	var pri Priority

	if len(buf) == 0 || (cur+priLen) > len(buf)-1 {
		return pri, cur, ErrBadPriority
	}

	if buf[cur] != priStart {
		return pri, cur, ErrBadPriority
	}

	cur++
	tokenStart := cur

	if buf[cur] == priEnd {
		return pri, cur, ErrBadPriority
	}

	for buf[cur] != priEnd {
		if !(buf[cur] >= '0' && buf[cur] <= '9') {
			return pri, cur, ErrBadPriority
		}

		cur++

		if cur > (priLen - 1) {
			return pri, cur, ErrBadPriority
		}
	}

	pVal, _ := strconv.Atoi(string(buf[tokenStart:cur]))

	pri = Priority{
		Priority: pVal,
		Facility: Facility(pVal / 8),
		Severity: Severity(pVal % 8),
	}

	cur++
	return pri, cur, err
}

// CheckForLikelyDateTime checks for a YYYY-MM-DD string. If one is found,
// we use this to decide that trying to parse a full rsyslog style timestamp
// is worth the cpu time.
func CheckForLikelyDateTime(buf []byte) bool {
	for i := 0; i < yearLen; i++ {
		if !unicode.IsDigit(rune(buf[i])) {
			return false
		}
	}

	if string(buf[startMonth-1]) != datePartSep {
		return false
	}

	for i := startMonth; i < startMonth+dayLen; i++ {
		if !unicode.IsDigit(rune(buf[i])) {
			return false
		}
	}

	if string(buf[startDay-1]) != datePartSep {
		return false
	}

	for i := startDay; i < startDay+dayLen; i++ {
		if !unicode.IsDigit(rune(buf[i])) {
			return false
		}
	}

	return true
}

func ParseTime(cur int, buf []byte) (time.Time, string, int, error) {
	var err error
	var foundTime bool
	var t time.Time
	var tf string

	// no timestamp format is shorter than YYYY-MM-DD, so if buffer is shorter
	// than this it is safe to assume we don't have a valid datetime.
	if cur+dateStampLen > len(buf)-1 {
		return t, tf, cur, ErrBadTime
	}

	if CheckForLikelyDateTime(buf[cur : cur+dateStampLen]) {
		tokenStart := cur
		tokenEnd := cur

		for buf[tokenEnd] != ' ' {
			tokenEnd++
			if tokenEnd > len(buf)-1 {
				return t, tf, cur, ErrBadTime
			}
		}

		timeStr := string(buf[tokenStart:tokenEnd])
		t, err = time.Parse(rsyslogTimeFormat, timeStr)
		if err == nil {
			cur = tokenEnd
			tf = rsyslogTimeFormat
		}
		return t, tf, cur, err
	}

	for _, timeFormat := range timeFormats {
		tLen := len(timeFormat)
		if cur+tLen > len(buf) {
			continue
		}

		timeStr := string(buf[cur : cur+tLen])
		t, err = time.Parse(timeFormat, timeStr)
		if err == nil {
			cur = cur + tLen
			tf = timeFormat
			foundTime = true
			break
		}
	}
	if !foundTime {
		err = ErrBadTime
	}
	return t, tf, cur, err
}

func ParseHost(cur int, buf []byte) (string, int, error) {
	var err error
	var host string

	for buf[cur] == ' ' {
		cur++
		if cur > len(buf)-1 {
			return host, cur, ErrBadHost
		}
	}

	tokenStart := cur

	for buf[cur] != ' ' {
		cur++
		if cur > len(buf)-1 {
			return host, cur, ErrBadHost
		}
	}

	host = string(buf[tokenStart:cur])
	return host, cur, err
}

func ParseTag(cur int, buf []byte) (string, string, string, int, error) {
	var err error
	var hasPid bool
	var hasColon bool
	var program string
	var pid string
	var tag string
	var tokenEnd int

	for buf[cur] == ' ' {
		cur++
		if cur > len(buf)-1 {
			return tag, program, pid, cur, ErrBadTag
		}
	}

	tokenStart := cur
	for {
		switch buf[cur] {
		case ':':
			cur++
			tokenEnd = cur
			hasColon = true
			goto FoundEndOfTag
		case ' ':
			tokenEnd = cur
			goto FoundEndOfTag
		case '[':
			program = string(buf[tokenStart:cur])
			cur++
			if cur > len(buf)-1 {
				return tag, program, pid, cur, ErrBadTag
			}
			pidStart := cur
			tokenEnd = cur
			for buf[cur] != ']' {
				cur++
				if cur > len(buf)-1 {
					return tag, program, pid, cur, ErrBadTag
				}
			}
			pidEnd := cur
			pid = string(buf[pidStart:pidEnd])
			hasPid = true
		}
		cur++
		if cur > len(buf)-1 {
			return tag, program, pid, cur, ErrBadTag
		}
	}
FoundEndOfTag:
	tag = string(buf[tokenStart:tokenEnd])
	if !hasPid {
		if !hasColon {
			program = tag
		} else {
			program = string(buf[tokenStart : tokenEnd-1])
		}
	}
	return tag, program, pid, cur, err
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
