package captainslog

import (
	"fmt"
	"testing"
	"time"
)

func TestNewLog(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
	if err != nil {
		t.Error(err)
	}

	if want, got := Local7, msg.Pri.Facility; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := Debug, msg.Pri.Severity; want != got {
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

	// host
	if want, got := "host.example.org", msg.Host; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	// tag
	if want, got := "test:", msg.Tag; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	// cee
	if want, got := false, msg.Cee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	// content
	if want, got := " hello world", msg.Content; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestNewLogCeeSpace(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
	if err != nil {
		t.Error(err)
	}

	if want, got := true, msg.Cee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	if want, got := "{\"a\":\"b\"}", msg.Content; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestNewLogCeeNoSpace(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:@cee:{\"a\":\"b\"}\n")
	var msg SyslogMsg
	err := Unmarshal(b, &msg)
	if err != nil {
		t.Error(err)
	}

	if want, got := true, msg.Cee; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestNoContent(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
	if err != nil {
		t.Error(err)
	}

	if want, got := "", msg.Content; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestTagEndHandling(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
	if err != nil {
		t.Error(err)
	}

	if want, got := "test:", msg.Tag; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	b = []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test hello world\n")
	err = Unmarshal(b, &msg)
	if err != nil {
		t.Error(err)
	}

	if want, got := "test", msg.Tag; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestParseUnixTime(t *testing.T) {
	b := []byte("<38>Mon Jan  2 15:04:05 host.example.org test: hello world\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
	if err != nil {
		t.Error(err)
	}

	ts := msg.Time

	if want, got := 0, ts.Year(); want != got {
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

func TestParseTimeANSIC(t *testing.T) {
	b := []byte("<38>Mon Jan  2 15:04:05 2006 host.example.org test: hello world\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
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

func TestParseTimeUnixDate(t *testing.T) {
	b := []byte("<38>Mon Jan  2 15:04:05 MST 2006 host.example.org test: hello world\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
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

func TestParseTimeNoYear(t *testing.T) {
	b := []byte("<38>Mon Jan  2 15:04:05 host.example.org test: hello world\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
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

func TestNoPriority(t *testing.T) {
	b := []byte("2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)

	if want, got := ErrBadPriority, err; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestNoPriorityEnd(t *testing.T) {
	b := []byte("<1912006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)

	if want, got := ErrBadPriority, err; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestPriorityTooLong(t *testing.T) {
	b := []byte("<9999>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)

	if want, got := ErrBadPriority, err; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestPriorityTruncated(t *testing.T) {
	b := []byte("<99\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)

	if want, got := ErrBadPriority, err; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestDateTruncated(t *testing.T) {
	b := []byte("<191>2006-01-02T15:0")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)

	if want, got := ErrBadTime, err; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestHostTruncated(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.examp")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)

	if want, got := ErrBadHost, err; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestTagTruncated(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org tes")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)

	if want, got := ErrBadTag, err; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestContentNotTerminated(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello wo")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)

	if want, got := ErrBadContent, err; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func ExampleUnmarshal() {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Syslog message was from host '%s'", msg.Host)
	// Output: Syslog message was from host 'host.example.org'

}

func BenchmarkParserParse(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	var msg SyslogMsg

	for i := 0; i < b.N; i++ {
		err := Unmarshal(m, &msg)
		if err != nil {
			panic(err)
		}
	}
}
