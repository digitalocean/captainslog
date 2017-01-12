package captainslog

import (
	"strings"
	"testing"
)

func TestSyslogMsgPlainWithAddedKeys(t *testing.T) {
	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: [15803005.789011] ------------[ cut here ]------------\n")

	msg := NewSyslogMsg()
	msg, err := NewSyslogMsgFromBytes(input)
	if err != nil {
		t.Error(err)
	}

	msg.JSONValues["tags"] = []string{"trace"}
	rfc3164 := msg.String()

	if want, got := true, strings.Contains(rfc3164, "tags"); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := true, strings.Contains(rfc3164, "msg"); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

}

func TestSyslogMsgJSONFromPlain(t *testing.T) {
	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: test\n")
	msg, err := NewSyslogMsgFromBytes(input)
	if err != nil {
		t.Error(err)
	}

	output, err := msg.JSON()
	if err != nil {
		t.Error(err)
	}

	wanted := `{"syslog_content":" test","syslog_facilitytext":"kern","syslog_host":"host.example.com","syslog_program":"kernel:","syslog_severitytext":"warning","syslog_time":"2016-03-08T14:59:36.293816Z"}`
	if want, got := wanted, string(output); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
func TestSyslogMsgJSONFromCEE(t *testing.T) {
	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: @cee:{\"a\":1}\n")
	msg, err := NewSyslogMsgFromBytes(input)
	if err != nil {
		t.Error(err)
	}

	output, err := msg.JSON()
	if err != nil {
		t.Error(err)
	}

	wanted := `{"a":1,"syslog_facilitytext":"kern","syslog_host":"host.example.com","syslog_program":"kernel:","syslog_severitytext":"warning","syslog_time":"2016-03-08T14:59:36.293816Z"}`
	if want, got := wanted, string(output); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestSyslogMsgJSONFromCEEWithDontParseJSON(t *testing.T) {
	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: @cee:{\"a\":1}\n")
	msg, err := NewSyslogMsgFromBytes(input, OptionDontParseJSON)
	if err != nil {
		t.Error(err)
	}

	output, err := msg.JSON()
	if err != nil {
		t.Error(err)
	}

	wanted := `{"a":1,"syslog_facilitytext":"kern","syslog_host":"host.example.com","syslog_program":"kernel:","syslog_severitytext":"warning","syslog_time":"2016-03-08T14:59:36.293816Z"}`
	if want, got := wanted, string(output); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
