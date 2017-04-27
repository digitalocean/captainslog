package captainslog_test

import (
	"testing"

	"github.com/digitalocean/captainslog"
	"github.com/prometheus/client_golang/prometheus"
)

func TestSextant(t *testing.T) {
	s, err := captainslog.NewSextant("test")
	if err != nil {
		t.Error(err)
	}
	defer close(s.Quit)

	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel[12]: test\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(input)
	if err != nil {
		t.Error(err)
	}

	s.Update(&msg)
}

func BenchmarkSextant(b *testing.B) {
	s, err := captainslog.NewSextant("bench")
	defer close(s.Quit)

	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n")
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		s.Update(&msg)
	}

	prometheus.Unregister(s.BytesTotal)
	prometheus.Unregister(s.LogLinesTotal)
	prometheus.Unregister(s.ParseErrorTotal)
	prometheus.Unregister(s.JSONLogsTotal)
	prometheus.Unregister(s.UniqueKeysTotal)

}
