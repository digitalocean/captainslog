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
	p.cur, p.msg.Pri, err = ParsePri(p.cur, p.buf)
	if err != nil {
		return err
	}

	var msgTime Time
	p.cur, msgTime, err = ParseTime(p.cur, p.buf)
	if err != nil {
		return err
	}

	p.msg.Time = msgTime.Time
	p.msg.timeFormat = msgTime.TimeFormat

	if p.optionNoHostname {
		host, err := os.Hostname()
		if err != nil {
			return ErrBadHost
		}
		p.msg.Host = host
	} else {
		p.cur, p.msg.Host, err = ParseHost(p.cur, p.buf)
		if err != nil {
			return err
		}
	}

	var msgTag Tag
	p.cur, msgTag, err = ParseTag(p.cur, p.buf)
	if err != nil {
		return err
	}
	p.msg.Tag = msgTag.Tag
	p.msg.Program = msgTag.Program
	p.msg.Pid = msgTag.Pid

	var cee string
	p.cur, cee, err = ParseCEE(p.cur, p.buf)
	if err != nil {
		return err
	}

	if cee != "" {
		p.msg.Cee = cee
		p.msg.IsCee = true
	}

	parseJSON := true
	if p.optionDontParseJSON {
		parseJSON = false
	}

	var content Content
	p.cur, content, err = ParseContent(p.cur, p.requireTerminator, p.msg.IsCee, parseJSON, p.buf)
	p.msg.Content = content.Content
	p.msg.JSONValues = content.JSONValues
	return err
}

func ParsePri(cur int, buf []byte) (int, Priority, error) {
	var err error
	var pri Priority

	if len(buf) == 0 || (cur+priLen) > len(buf)-1 {
		return cur, pri, ErrBadPriority
	}

	if buf[cur] != priStart {
		return cur, pri, ErrBadPriority
	}

	cur++
	tokenStart := cur

	if buf[cur] == priEnd {
		return cur, pri, ErrBadPriority
	}

	for buf[cur] != priEnd {
		if !(buf[cur] >= '0' && buf[cur] <= '9') {
			return cur, pri, ErrBadPriority
		}

		cur++

		if cur > (priLen - 1) {
			return cur, pri, ErrBadPriority
		}
	}

	pVal, _ := strconv.Atoi(string(buf[tokenStart:cur]))

	pri = Priority{
		Priority: pVal,
		Facility: Facility(pVal / 8),
		Severity: Severity(pVal % 8),
	}

	cur++
	return cur, pri, err
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

func ParseTime(cur int, buf []byte) (int, Time, error) {
	var err error
	var foundTime bool
	var msgTime Time

	// no timestamp format is shorter than YYYY-MM-DD, so if buffer is shorter
	// than this it is safe to assume we don't have a valid datetime.
	if cur+dateStampLen > len(buf)-1 {
		return cur, msgTime, ErrBadTime
	}

	if CheckForLikelyDateTime(buf[cur : cur+dateStampLen]) {
		tokenStart := cur
		tokenEnd := cur

		for buf[tokenEnd] != ' ' {
			tokenEnd++
			if tokenEnd > len(buf)-1 {
				return cur, msgTime, ErrBadTime
			}
		}

		timeStr := string(buf[tokenStart:tokenEnd])
		msgTime.Time, err = time.Parse(rsyslogTimeFormat, timeStr)
		if err == nil {
			cur = tokenEnd
			msgTime.TimeFormat = rsyslogTimeFormat
		}
		return cur, msgTime, err
	}

	for _, timeFormat := range timeFormats {
		tLen := len(timeFormat)
		if cur+tLen > len(buf) {
			continue
		}

		timeStr := string(buf[cur : cur+tLen])
		msgTime.Time, err = time.Parse(timeFormat, timeStr)
		if err == nil {
			cur = cur + tLen
			msgTime.TimeFormat = timeFormat
			foundTime = true
			break
		}
	}
	if !foundTime {
		err = ErrBadTime
	}
	return cur, msgTime, err
}

func ParseHost(cur int, buf []byte) (int, string, error) {
	var err error
	var host string

	for buf[cur] == ' ' {
		cur++
		if cur > len(buf)-1 {
			return cur, host, ErrBadHost
		}
	}

	tokenStart := cur

	for buf[cur] != ' ' {
		cur++
		if cur > len(buf)-1 {
			return cur, host, ErrBadHost
		}
	}

	host = string(buf[tokenStart:cur])
	return cur, host, err
}

func ParseTag(cur int, buf []byte) (int, Tag, error) {
	var err error
	var hasPid bool
	var hasColon bool
	var tag Tag
	var tokenEnd int

	for buf[cur] == ' ' {
		cur++
		if cur > len(buf)-1 {
			return cur, tag, ErrBadTag
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
			tag.Program = string(buf[tokenStart:cur])
			cur++
			if cur > len(buf)-1 {
				return cur, tag, ErrBadTag
			}
			pidStart := cur
			tokenEnd = cur
			for buf[cur] != ']' {
				cur++
				if cur > len(buf)-1 {
					return cur, tag, ErrBadTag
				}
			}
			pidEnd := cur
			tag.Pid = string(buf[pidStart:pidEnd])
			hasPid = true
		}
		cur++
		if cur > len(buf)-1 {
			return cur, tag, ErrBadTag
		}
	}
FoundEndOfTag:
	tag.Tag = string(buf[tokenStart:tokenEnd])
	if !hasPid {
		if !hasColon {
			tag.Program = tag.Tag
		} else {
			tag.Program = string(buf[tokenStart : tokenEnd-1])
		}
	}
	return cur, tag, err
}

func ParseCEE(cur int, buf []byte) (int, string, error) {
	var err error
	var cee string

	if cur >= len(buf)-1 {
		return cur, cee, err
	}

	tokenStart := cur
	tokenEnd := cur

	for buf[tokenEnd] == ' ' {
		tokenEnd++
		if tokenEnd >= len(buf)-1 {
			return cur, cee, err
		}
	}

	if tokenEnd+4 > len(buf)-1 {
		return cur, cee, err
	}

	if buf[tokenEnd] != '@' {
		return cur, cee, err
	}

	tokenEnd++
	if buf[tokenEnd] != 'c' {
		return cur, cee, err
	}

	tokenEnd++
	if buf[tokenEnd] != 'e' {
		return cur, cee, err
	}

	tokenEnd++
	if buf[tokenEnd] != 'e' {
		return cur, cee, err
	}

	tokenEnd++
	if buf[tokenEnd] != ':' {
		return cur, cee, err
	}

	tokenEnd++
	cur = tokenEnd
	cee = string(buf[tokenStart:tokenEnd])
	return cur, cee, err
}

func ParseContent(cur int, requireTerminator bool, isCee bool, parseJSON bool, buf []byte) (int, Content, error) {
	content := Content{JSONValues: make(map[string]interface{}, 0)}
	var err error
	tokenStart := cur

	if cur >= len(buf)-1 {
		return cur, content, ErrBadContent
	}

	for buf[cur] != '\n' {
		cur++
		if cur > len(buf)-1 {
			if requireTerminator {
				return cur, content, ErrBadContent
			}
			goto exitContentSearch
		}
	}
exitContentSearch:
	tokenEnd := cur

	if parseJSON && isCee {
		decoder := json.NewDecoder(bytes.NewBuffer(buf[tokenStart:tokenEnd]))
		decoder.UseNumber()
		decoder.Decode(&content.JSONValues)
	}
	content.Content = string(buf[tokenStart:tokenEnd])
	return cur, content, err
}
