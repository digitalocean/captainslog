package captainslog_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/digitalocean/captainslog"
)

func TestNewParserWithPidAndParse(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[12]: hello world\n")
	p := captainslog.NewParser()

	msg, err := p.ParseBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := captainslog.Local7, msg.Pri.Facility; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := captainslog.Debug, msg.Pri.Severity; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	ts := msg.Time

	if want, got := 2006, ts.Year(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := time.Month(1), ts.Month(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 2, ts.Day(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 15, ts.Hour(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 4, ts.Minute(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 5, ts.Second(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 999999, ts.Nanosecond()/1000; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	_, zoneOffsetSecs := ts.Zone()
	if want, got := -25200, zoneOffsetSecs; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "host.example.org", msg.Host; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test", msg.Program; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test[12]:", msg.Tag; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "12", msg.Pid; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := false, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := " hello world", msg.Content; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := 0, bytes.Compare(b, msg.Bytes()); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := string(b), msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

}

func TestNewParserAndParse(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	p := captainslog.NewParser()

	msg, err := p.ParseBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := captainslog.Local7, msg.Pri.Facility; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := captainslog.Debug, msg.Pri.Severity; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	ts := msg.Time

	if want, got := 2006, ts.Year(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := time.Month(1), ts.Month(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 2, ts.Day(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 15, ts.Hour(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 4, ts.Minute(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 5, ts.Second(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 999999, ts.Nanosecond()/1000; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	_, zoneOffsetSecs := ts.Zone()
	if want, got := -25200, zoneOffsetSecs; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "host.example.org", msg.Host; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test:", msg.Tag; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test", msg.Program; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := false, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := " hello world", msg.Content; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := 0, bytes.Compare(b, msg.Bytes()); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := string(b), msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

}

func TestParserOptionNoHostname(t *testing.T) {
	b := []byte("<86>Jul 24 11:53:47 sudo: pam_unix(sudo:session): session opened for user root by (uid=0)")
	p := captainslog.NewParser(captainslog.OptionNoHostname)

	host, err := os.Hostname()
	if err != nil {
		t.Error(err)
	}

	msg, err := p.ParseBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := host, msg.Host; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestParseOptionDontParseJSON(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n")
	p := captainslog.NewParser(captainslog.OptionDontParseJSON)

	msg, err := p.ParseBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := true, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := " @cee:", msg.Cee; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "{\"a\":\"b\"}", msg.Content; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := string(b), msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestNewSyslogMsgFromBytes(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := captainslog.Local7, msg.Pri.Facility; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := captainslog.Debug, msg.Pri.Severity; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	ts := msg.Time

	if want, got := 2006, ts.Year(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := time.Month(1), ts.Month(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 2, ts.Day(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 15, ts.Hour(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 4, ts.Minute(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 5, ts.Second(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 999999, ts.Nanosecond()/1000; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	_, zoneOffsetSecs := ts.Zone()
	if want, got := -25200, zoneOffsetSecs; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "host.example.org", msg.Host; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test", msg.Program; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test:", msg.Tag; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := false, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := " hello world", msg.Content; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := 0, bytes.Compare(b, msg.Bytes()); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := string(b), msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalDateNoMicros(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999-07:00 host.example.org test: hello world\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := captainslog.Local7, msg.Pri.Facility; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := captainslog.Debug, msg.Pri.Severity; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	ts := msg.Time

	if want, got := 2006, ts.Year(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := time.Month(1), ts.Month(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 2, ts.Day(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 15, ts.Hour(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 4, ts.Minute(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 5, ts.Second(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 999000, ts.Nanosecond()/1000; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	_, zoneOffsetSecs := ts.Zone()
	if want, got := -25200, zoneOffsetSecs; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "host.example.org", msg.Host; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test", msg.Program; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test:", msg.Tag; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := false, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := " hello world", msg.Content; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := 0, bytes.Compare(b, msg.Bytes()); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := string(b), msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalDateNoMillis(t *testing.T) {
	b := []byte("<171>2015-12-18T18:08:17+00:00 host.example.org test: hello world\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := captainslog.Local5, msg.Pri.Facility; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := captainslog.Err, msg.Pri.Severity; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	ts := msg.Time

	if want, got := 2015, ts.Year(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := time.Month(12), ts.Month(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 18, ts.Day(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 18, ts.Hour(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 8, ts.Minute(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 17, ts.Second(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 0, ts.Nanosecond()/1000; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	_, zoneOffsetSecs := ts.Zone()
	if want, got := 0, zoneOffsetSecs; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "host.example.org", msg.Host; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test", msg.Program; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test:", msg.Tag; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := false, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := " hello world", msg.Content; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := 0, bytes.Compare(b, msg.Bytes()); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := string(b), msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalCeeSpace(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := true, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := " @cee:", msg.Cee; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "{\"a\":\"b\"}", msg.Content; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := string(b), msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalCeeNoSpace(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:@cee:{\"a\":\"b\"}\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := true, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := "@cee:", msg.Cee; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := string(b), msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalCeeEarlyBufferBeforeColon(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:@cee\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := false, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := "", msg.Cee; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "@cee", msg.Content; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalCeeEarlyBufferAfterColon(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:@cee:\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadContent, err; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := true, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func unmarshalCeeButNotCee(t *testing.T, b []byte) {
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := false, msg.IsCee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestUnmarshalCeeButNotCee(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:@cee{\"a\":\"b\"}\n")
	unmarshalCeeButNotCee(t, b)

	b = []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:@ce{\"a\":\"b\"}\n")
	unmarshalCeeButNotCee(t, b)

	b = []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:@c{\"a\":\"b\"}\n")
	unmarshalCeeButNotCee(t, b)

	b = []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:@{\"a\":\"b\"}\n")
	unmarshalCeeButNotCee(t, b)
}

func TestUnmarshalNoContent(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:\n")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadContent, err; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestUnmarshalTagEndHandling(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := "test", msg.Program; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test:", msg.Tag; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	b = []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test hello world\n")
	msg, err = captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := "test", msg.Program; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := "test", msg.Tag; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if want, got := string(b), msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestHandlingTruncatedSubseconds(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.99999-07:00 host.example.org test: hello world\n")
	_, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}
}

func TestUnmarshalUnixTime(t *testing.T) {
	b := []byte("<38>Mon Jan  2 15:04:05 host.example.org test: hello world\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	ts := msg.Time

	if want, got := time.Now().Year(), ts.Year(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := time.Month(1), ts.Month(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 2, ts.Day(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 15, ts.Hour(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 4, ts.Minute(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 5, ts.Second(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}
}

func TestUnmarshalTimeANSIC(t *testing.T) {
	b := []byte("<38>Mon Jan  2 15:04:05 2006 host.example.org test: hello world\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	ts := msg.Time

	if want, got := time.Month(1), ts.Month(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 2, ts.Day(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 15, ts.Hour(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 4, ts.Minute(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 5, ts.Second(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := string(b), msg.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalTimeUnixDate(t *testing.T) {
	b := []byte("<38>Mon Jan  2 15:04:05 MST 2006 host.example.org test: hello world\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	ts := msg.Time

	if want, got := time.Month(1), ts.Month(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 2, ts.Day(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 15, ts.Hour(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 4, ts.Minute(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 5, ts.Second(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	zone, _ := ts.Zone()
	if want, got := "MST", zone; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalTimeNoYear(t *testing.T) {
	b := []byte("<38>Mon Jan  2 15:04:05 host.example.org test: hello world\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	ts := msg.Time

	if want, got := time.Month(1), ts.Month(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 2, ts.Day(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 15, ts.Hour(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 4, ts.Minute(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := 5, ts.Second(); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}
}

func TestUnmarshalNoPriority(t *testing.T) {
	b := []byte("2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadPriority, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalNoPriorityEnd(t *testing.T) {
	b := []byte("<1912006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadPriority, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalPriorityTooLong(t *testing.T) {
	b := []byte("<9999>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadPriority, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalPriorityTruncated(t *testing.T) {
	b := []byte("<99\n")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadPriority, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalDateTruncated(t *testing.T) {
	b := []byte("<191>2006-01-02T15:0")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadTime, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalHostTruncated(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.examp")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadHost, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalNoHost(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 ")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadHost, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalTagTruncated(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org tes")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadTag, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalNoTag(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org ")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadTag, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestUnmarshalContentNotTerminated(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello wo")
	_, err := captainslog.NewSyslogMsgFromBytes(b)
	// 2016-07-22: changed this test to pass as default captainslog.NewSyslogMsgFromBytes now
	// treats end buffer as end of content.
	if err != nil {
		t.Error(err)
	}
}

func TestUnmarshalPriNotNumber(t *testing.T) {
	b := []byte("<1a1>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := captainslog.ErrBadPriority, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestTagWithColonNoPid(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	p := captainslog.NewParser()

	msg, err := p.ParseBytes(b)
	if err != nil {
		t.Error(err)
	}

	if want, got := "test", msg.Program; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
	if want, got := "test:", msg.Tag; want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func testFuzzFindings(fuzzData string, t *testing.T) {
	b := []byte(fuzzData)
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := false, err == nil; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestFuzzFindings(t *testing.T) {
	inputs := []string{
		"<0>Mon Jan 00 00:00:000 0 ",
		"<0>Mon Jan 00 00:00:000 :",
	}

	for _, fuzzData := range inputs {
		testFuzzFindings(fuzzData, t)
	}
}

func TestBadPids(t *testing.T) {
	inputs := []string{
		"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[1: hello world\n",
		"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[1 hello world\n",
		"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[\n",
		"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[]\n",
	}

	for _, fuzzData := range inputs {
		testFuzzFindings(fuzzData, t)
	}
}

func BenchmarkParserParse(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Content != " hello world" {
			panic("unexpected msg.Content")
		}
	}
}

func BenchmarkParserParseCEE(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Host != "host.example.org" {
			panic("unexpected msg.Host")
		}
	}
}

func BenchmarkParserParseCEEWithOptionDontParseJSON(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m, captainslog.OptionDontParseJSON)
		if err != nil {
			panic(err)
		}
		if msg.Host != "host.example.org" {
			panic("unexpected msg.Host")
		}
	}
}

func BenchmarkParserParseLeastLikelyTime(b *testing.B) {
	m := []byte("<38>Mon Jan  2 15:04:05 host.example.org test: hello world\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Content != " hello world" {
			panic("unexpected msg.Content")
		}
	}
}

func BenchmarkParserParseAndString(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Content != " hello world" {
			panic("unexpected msg.Content")
		}
	}
}

func BenchmarkParserParseAndBytes(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Content != " hello world" {
			panic("unexpected msg.Content")
		}
	}
}

func BenchmarkJSONParseNoJSON(b *testing.B) {
	m := []byte("hello world")
	val := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(m, &val)
		if err == nil {
			panic("wat")
		}
	}
}

func BenchmarkJSONCheckFirstChar(b *testing.B) {
	m := []byte("hello world")
	for i := 0; i < b.N; i++ {
		if bytes.Compare(m[0:0], []byte("{")) == 0 {
			panic("wee")
		}
	}
}

func BenchmarkParserParseInvalidDate(b *testing.B) {
	m := []byte("<191>2006-02-30T15:04:05.999999-07:00 host.example.org test: hello world\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		_, err := captainslog.NewSyslogMsgFromBytes(m)
		if err == nil {
			panic(err)
		}
	}
}

func BenchmarkParserParseInvaliSyslog(b *testing.B) {
	m := []byte("Hello I am not a syslog message\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		_, err := captainslog.NewSyslogMsgFromBytes(m)
		if err == nil {
			panic(err)
		}
	}
}
