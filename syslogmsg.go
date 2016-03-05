package captainslog

import (
	"encoding/json"
	"fmt"
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
