package captainslog

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
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
	mutex      *sync.Mutex
}

// NewSyslogMsg creates a new empty SyslogMsg.
func NewSyslogMsg() SyslogMsg {
	return SyslogMsg{
		JSONValues: make(map[string]interface{}),
		mutex:      &sync.Mutex{},
	}
}

// NewSyslogMsgFromBytes accepts a []byte containing an RFC3164
// message and returns a SyslogMsg. If the original RFC3164
// message is a CEE enhanced message, the JSON will be
// parsed into the JSONValues map[string]inferface{}
func NewSyslogMsgFromBytes(b []byte) (SyslogMsg, error) {
	msg := NewSyslogMsg()
	err := Unmarshal(b, &msg)
	return msg, err
}

// AddTagArray adds a tag to an array of tags at the key. If the key
// does not already exist, it will create the key and initially it
// to a []interface{}.
func (s *SyslogMsg) AddTagArray(key string, value interface{}) error {
	if _, ok := s.JSONValues[key]; !ok {
		s.JSONValues[key] = make([]interface{}, 0)
	}

	switch val := s.JSONValues[key].(type) {
	case []interface{}:
		s.JSONValues[key] = append(val, value)
		if !s.IsCee {
			s.IsCee = true
			s.Cee = " @cee:"
			s.JSONValues["msg"] = s.Content[1:]
		}
		return nil
	default:
		return fmt.Errorf("tags key in message was not an array")
	}
}

// AddTag adds a tag to the value at key. If the key exists,
// the value currently at the key will be overwritten.
func (s *SyslogMsg) AddTag(key string, value interface{}) {
	s.JSONValues[key] = value
}

// String returns the SyslogMsg as an RFC3164 string.
func (s *SyslogMsg) String() string {
	var content string
	if s.IsCee {
		b, err := json.Marshal(s.JSONValues)
		if err != nil {
			panic(err)
		}
		content = string(b)
	} else {
		if len(s.JSONValues) > 0 {
			s.JSONValues["msg"] = strings.TrimLeft(s.Content, " ")
			s.IsCee = true
			s.Cee = " @cee:"
			b, err := json.Marshal(s.JSONValues)
			if err != nil {
				panic(err)
			}
			content = string(b)
		} else {
			content = s.Content
		}
	}
	return fmt.Sprintf("<%d>%s %s %s%s%s\n", s.Pri.Priority, s.Time.Format(s.timeFormat), s.Host, s.Tag, s.Cee, content)
}

// Bytes returns the SyslogMsg as RFC3164 []byte.
func (s *SyslogMsg) Bytes() []byte {
	return []byte(s.String())
}

func (s *SyslogMsg) JSON() ([]byte, error) {
	content := make(map[string]interface{})
	for key, value := range s.JSONValues {
		content[key] = value
	}

	content["syslog_time"] = s.Time
	content["syslog_host"] = s.Host
	content["syslog_program"] = s.Tag

	b, err := json.Marshal(content)
	return b, err
}
