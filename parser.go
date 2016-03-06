package captainslog

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

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
	JSONValues map[string]interface{}
}

// String returns the SyslogMsg as an RFC3164 string
func (s *SyslogMsg) String() string {
	var content string
	if s.IsCee {
		b, err := json.Marshal(s.JSONValues)
		if err != nil {
			panic(err)
		}
		content = string(b)
	} else {
		content = s.Content
	}
	return fmt.Sprintf("<%d>%s %s %s%s%s\n", s.Pri.Priority, s.Time.Format(s.timeFormat), s.Host, s.Tag, s.Cee, content)
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
	var err error

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

	if p.msg.IsCee {
		err = json.Unmarshal(p.buf[p.tokenStart:p.tokenEnd], &p.msg.JSONValues)
		if err != nil {
			p.msg.IsCee = false
		}
	}
	p.msg.Content = string(p.buf[p.tokenStart:p.tokenEnd])
	return err
}
