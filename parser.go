package captainslog

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

const (
	priStart = '<'
	priEnd   = '>'
	priLen   = 5
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

	timeFormats = []string{
		"2006-01-02T15:04:05.999999-07:00",
		"2006-01-02T15:04:05.999-07:00",
		"2006-01-02T15:04:05-07:00",
		"Mon Jan _2 15:04:05 MST 2006",
		"Mon Jan _2 15:04:05 2006",
		"Mon Jan _2 15:04:05",
		"Jan _2 15:04:05",
		"Jan 02 15:04:05",
	}
)

type parser struct {
	tokenStart        int
	tokenEnd          int
	buf               []byte
	bufLen            int
	bufEnd            int
	cur               int
	requireTerminator bool
	msg               *SyslogMsg
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

	err = p.parseCee()
	if err != nil {
		return err
	}

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

	if p.buf[p.cur] == priEnd {
		return ErrBadPriority
	}

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

func (p *parser) parseCee() error {
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

func (p *parser) parseContent() error {
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
			} else {
				goto exitContentSearch
			}
		}
	}
exitContentSearch:
	p.tokenEnd = p.cur

	if p.msg.IsCee {
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
