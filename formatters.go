package captainslog

import "encoding/json"

// ToJSONBytes converts the SyslogMsg to pure JSON
func ToJSONBytes(s *SyslogMsg) []byte {
	m := make(map[string]interface{})
	for k, v := range s.JSONValues {
		m[k] = v
	}

	m["syslog_host"] = s.Host
	m["syslog_program"] = s.Tag
	m["@timestamp"] = s.Time.String()
	m["content"] = s.Content
	b, _ := json.Marshal(m)
	return b
}
