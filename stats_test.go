package captainslog_test

import (
	"testing"

	"github.com/digitalocean/captainslog"
)

func TestStats(t *testing.T) {
	s := captainslog.NewStats("foo")
	s.Unregister()

}

func BenchmarkStats(b *testing.B) {
	s := captainslog.NewStats("foo")
	for i := 0; i < b.N; i++ {
		s.BytesTotal.Add(800)
		s.LogLinesTotal.Inc()
		s.ParseErrorTotal.Inc()
		s.JSONLogsTotal.Inc()
		s.UniqueKeysTotal.Inc()
		s.UniqueProgramsTotal.Inc()
	}
	s.Unregister()
}
