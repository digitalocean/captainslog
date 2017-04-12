package captainslog_test

import (
	"strings"
	"testing"

	"github.com/digitalocean/captainslog"
)

func TestSyslogMsgToStringWithPid(t *testing.T) {
	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel[12]: test\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(input)
	if err != nil {
		t.Error(err)
	}

	wanted := "<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel[12]: test\n"
	if want, got := wanted, msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestSyslogMsgPlainWithAddedKeys(t *testing.T) {
	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: [15803005.789011] ------------[ cut here ]------------\n")

	msg := captainslog.NewSyslogMsg()
	msg, err := captainslog.NewSyslogMsgFromBytes(input)
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
	msg, err := captainslog.NewSyslogMsgFromBytes(input)
	if err != nil {
		t.Error(err)
	}

	output, err := msg.JSON()
	if err != nil {
		t.Error(err)
	}
	wanted := `{"syslog_content":" test","syslog_facilitytext":"kern","syslog_host":"host.example.com","syslog_pid":"","syslog_programname":"kernel","syslog_severitytext":"warning","syslog_tag":"kernel:","syslog_time":"2016-03-08T14:59:36.293816Z"}`
	if want, got := wanted, string(output); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestSyslogMsgJSONFromCEE(t *testing.T) {
	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com test[12]: @cee:{\"a\":1}\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(input)
	if err != nil {
		t.Error(err)
	}

	output, err := msg.JSON()
	if err != nil {
		t.Error(err)
	}
	wanted := `{"a":1,"syslog_facilitytext":"kern","syslog_host":"host.example.com","syslog_pid":"12","syslog_programname":"test","syslog_severitytext":"warning","syslog_tag":"test[12]:","syslog_time":"2016-03-08T14:59:36.293816Z"}`
	if want, got := wanted, string(output); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestSyslogMsgJSONFromCEEWithDontParseJSON(t *testing.T) {
	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: @cee:{\"a\":1}\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(input, captainslog.OptionDontParseJSON)
	if err != nil {
		t.Error(err)
	}

	output, err := msg.JSON()
	if err != nil {
		t.Error(err)
	}

	wanted := `{"a":1,"syslog_facilitytext":"kern","syslog_host":"host.example.com","syslog_pid":"","syslog_programname":"kernel","syslog_severitytext":"warning","syslog_tag":"kernel:","syslog_time":"2016-03-08T14:59:36.293816Z"}`

	if want, got := wanted, string(output); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestNginxToSyslogMsgBackToString(t *testing.T) {
	input := []byte("<174>2017-04-12T13:31:11.918068+00:00 www.example.com nginx 192.168.1.1 - - [12/Apr/2017:13:31:11 +0000] \"GET /hello?from=world HTTP/1.1\" 200 18 \"https://something.example.com\" \"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36\"")
	msg, err := captainslog.NewSyslogMsgFromBytes(input)
	if err != nil {
		t.Error(err)
	}
	wanted := string(input) + "\n"
	if want, got := wanted, msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
