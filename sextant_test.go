package captainslog_test

import (
	"testing"

	"github.com/digitalocean/captainslog"
)

func TestSextant(t *testing.T) {
	s := captainslog.NewSextant("test")

	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel[12]: test\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(input)
	if err != nil {
		t.Error(err)
	}

	s.Update(&msg)
}
