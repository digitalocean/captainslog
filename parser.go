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

// ParseBytes accepts a []byte and tries to parse it into a SyslogMsg.
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
	var offset int

	offset, p.msg.Pri, err = ParsePri(p.buf)
	if err != nil {
		return err
	}
	p.cur = p.cur + offset

	var msgTime Time
	offset, msgTime, err = ParseTime(p.buf[p.cur:])
	if err != nil {
		return err
	}
	p.cur = p.cur + offset

	p.msg.Time = msgTime.Time
	p.msg.timeFormat = msgTime.TimeFormat

	if p.optionNoHostname {
		host, err := os.Hostname()
		if err != nil {
			return ErrBadHost
		}
		p.msg.Host = host
	} else {
		offset, p.msg.Host, err = ParseHost(p.buf[p.cur:])
		if err != nil {
			return err
		}
		p.cur = p.cur + offset
	}

	var msgTag Tag
	offset, msgTag, err = ParseTag(p.buf[p.cur:])
	if err != nil {
		return err
	}
	p.cur = p.cur + offset
	p.msg.Tag = msgTag.Tag
	p.msg.Program = msgTag.Program
	p.msg.Pid = msgTag.Pid

	var cee string
	offset, cee, err = ParseCEE(p.buf[p.cur:])
	if err != nil {
		return err
	}
	p.cur = p.cur + offset

	if cee != "" {
		p.msg.Cee = cee
		p.msg.IsCee = true
	}

	copts := make([]func(*contentOpts), 0)
	if !p.optionDontParseJSON {
		copts = append(copts, ContentOptionParseJSON)
	}

	if p.requireTerminator {
		copts = append(copts, ContentOptionRequireTerminator)
	}

	var content Content
	_, content, err = ParseContent(p.buf[p.cur:], copts...)
	p.msg.Content = content.Content
	p.msg.JSONValues = content.JSONValues
	return err
}

// ParsePri will try to find a syslog priority at the
// beginning of the passed in []byte. It will return the offset
// from the start of the []byte to the end of the priority string,
// a captainslog.Priority, and an error.
func ParsePri(buf []byte) (int, Priority, error) {
	var err error
	var pri Priority
	var offset int

	if len(buf) == 0 || (offset+priLen) > len(buf)-1 {
		return offset, pri, ErrBadPriority
	}

	if buf[offset] != priStart {
		return offset, pri, ErrBadPriority
	}

	offset++
	tokenStart := offset

	if buf[offset] == priEnd {
		return offset, pri, ErrBadPriority
	}

	for buf[offset] != priEnd {
		if !(buf[offset] >= '0' && buf[offset] <= '9') {
			return offset, pri, ErrBadPriority
		}

		offset++

		if offset > (priLen - 1) {
			return offset, pri, ErrBadPriority
		}
	}

	pVal, _ := strconv.Atoi(string(buf[tokenStart:offset]))

	pri = Priority{
		Priority: pVal,
		Facility: Facility(pVal / 8),
		Severity: Severity(pVal % 8),
	}

	offset++
	return offset, pri, err
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

// ParseTime will try to find a syslog time at the beginning of the
// passed in []byte. It returns the offset from the start of the []byte
// to the end of the time string, a captainslog.Time, and an error.
func ParseTime(buf []byte) (int, Time, error) {
	var err error
	var foundTime bool
	var msgTime Time
	var offset int

	// no timestamp format is shorter than YYYY-MM-DD, so if buffer is shorter
	// than this it is safe to assume we don't have a valid datetime.
	if offset+dateStampLen > len(buf)-1 {
		return offset, msgTime, ErrBadTime
	}

	if CheckForLikelyDateTime(buf[offset : offset+dateStampLen]) {
		tokenStart := offset
		tokenEnd := offset

		for buf[tokenEnd] != ' ' {
			tokenEnd++
			if tokenEnd > len(buf)-1 {
				return offset, msgTime, ErrBadTime
			}
		}

		timeStr := string(buf[tokenStart:tokenEnd])
		msgTime.Time, err = time.Parse(rsyslogTimeFormat, timeStr)
		if err == nil {
			offset = tokenEnd
			msgTime.TimeFormat = rsyslogTimeFormat
		}
		return offset, msgTime, err
	}

	for _, timeFormat := range timeFormats {
		tLen := len(timeFormat)
		if offset+tLen > len(buf) {
			continue
		}

		timeStr := string(buf[offset : offset+tLen])
		msgTime.Time, err = time.Parse(timeFormat, timeStr)
		if err == nil {
			offset = offset + tLen
			msgTime.TimeFormat = timeFormat
			foundTime = true
			break
		}
	}
	if !foundTime {
		err = ErrBadTime
	}
	return offset, msgTime, err
}

// ParseHost will try to find a host at the
// beginning of the passed in []byte. It will return the offset
// from the start of the []byte to the end of the host string,
// a captainslog.Priority, and an error.
func ParseHost(buf []byte) (int, string, error) {
	var err error
	var host string
	var offset int

	if offset > len(buf)-1 {
		return offset, host, ErrBadHost
	}

	for buf[offset] == ' ' {
		offset++
		if offset > len(buf)-1 {
			return offset, host, ErrBadHost
		}
	}

	tokenStart := offset

	for buf[offset] != ' ' {
		offset++
		if offset > len(buf)-1 {
			return offset, host, ErrBadHost
		}
	}

	host = string(buf[tokenStart:offset])
	return offset, host, err
}

// ParseTag will try to find a syslog tag at the beginning of the
// passed in []byte. It returns the offset from the start of the []byte
// to the end of the tag string, a captainslog.Tag, and an error.
func ParseTag(buf []byte) (int, Tag, error) {
	var err error
	var hasPid bool
	var hasColon bool
	var tag Tag
	var tokenEnd int
	var offset int

	for buf[offset] == ' ' {
		offset++
		if offset > len(buf)-1 {
			return offset, tag, ErrBadTag
		}
	}

	tokenStart := offset
	for {
		switch buf[offset] {
		case ':':
			offset++
			tokenEnd = offset
			hasColon = true
			goto FoundEndOfTag
		case ' ':
			tokenEnd = offset
			goto FoundEndOfTag
		case '[':
			tag.Program = string(buf[tokenStart:offset])
			offset++
			if offset > len(buf)-1 {
				return offset, tag, ErrBadTag
			}
			pidStart := offset
			tokenEnd = offset
			for buf[offset] != ']' {
				offset++
				if offset > len(buf)-1 {
					return offset, tag, ErrBadTag
				}
			}
			pidEnd := offset
			tag.Pid = string(buf[pidStart:pidEnd])
			hasPid = true
		}
		offset++
		if offset > len(buf)-1 {
			return offset, tag, ErrBadTag
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
	return offset, tag, err
}

// ParseCEE will try to find a syslog cee cookie  at the beginning of the
// passed in []byte. It returns the offset from the start of the []byte
// to the end of the cee string, the string, and an error.
func ParseCEE(buf []byte) (int, string, error) {
	var err error
	var cee string
	var offset int

	if offset >= len(buf)-1 {
		return offset, cee, err
	}

	tokenStart := offset
	tokenEnd := offset

	for buf[tokenEnd] == ' ' {
		tokenEnd++
		if tokenEnd >= len(buf)-1 {
			return offset, cee, err
		}
	}

	if tokenEnd+4 > len(buf)-1 {
		return offset, cee, err
	}

	if buf[tokenEnd] != '@' {
		return offset, cee, err
	}

	tokenEnd++
	if buf[tokenEnd] != 'c' {
		return offset, cee, err
	}

	tokenEnd++
	if buf[tokenEnd] != 'e' {
		return offset, cee, err
	}

	tokenEnd++
	if buf[tokenEnd] != 'e' {
		return offset, cee, err
	}

	tokenEnd++
	if buf[tokenEnd] != ':' {
		return offset, cee, err
	}

	tokenEnd++
	offset = tokenEnd
	cee = string(buf[tokenStart:tokenEnd])
	return offset, cee, err
}

type contentOpts struct {
	requireTerminator bool
	parseJSON         bool
}

// ContentOptionRequireTerminator sets ParseContent to require a \n terminator
func ContentOptionRequireTerminator(opts *contentOpts) {
	opts.requireTerminator = true
}

// ContentOptionParseJSON will treat the content as a CEE message
func ContentOptionParseJSON(opts *contentOpts) {
	opts.parseJSON = true
}

// ParseContent will try to find syslog content at the beginning of the
// passed in []byte. It returns the offset from the start of the []byte
// to the end of the content, a captainslog.Content, and an error. It
// accepts two options:
//
// ContentOptionRequireTerminator: if true, if the syslog message does not
//		contain a '\n' terminator it will be treated as invalid.
//
// ContentOptionParseJSON: if true, it will treat the content field of the
//		syslog message as a CEE message and parse the JSON.
func ParseContent(buf []byte, options ...func(*contentOpts)) (int, Content, error) {
	var o contentOpts
	for _, option := range options {
		option(&o)
	}

	content := Content{JSONValues: make(map[string]interface{}, 0)}

	var err error
	var offset int
	var probablyJSON bool

	if offset >= len(buf)-1 {
		return offset, content, ErrBadContent
	}

	tokenStart := offset
	for buf[offset] == ' ' {
		offset++
		if offset >= len(buf)-1 {
			if o.requireTerminator {
				return offset, content, ErrBadContent
			}
			break
		}
	}

	if buf[offset] == '{' {
		probablyJSON = true
	}

	for buf[offset] != '\n' {
		offset++
		if offset > len(buf)-1 {
			if o.requireTerminator {
				return offset, content, ErrBadContent
			}
			break
		}
	}

	if o.parseJSON && probablyJSON {
		decoder := json.NewDecoder(bytes.NewBuffer(buf[tokenStart:offset]))
		decoder.UseNumber()
		decoder.Decode(&content.JSONValues)
	}
	content.Content = string(buf[tokenStart:offset])
	return offset, content, err
}
