package captainslog_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/digitalocean/captainslog"
)

func TestSextant(t *testing.T) {
	s, err := captainslog.NewSextant("test", 0.03, 1)
	if err != nil {
		t.Error(err)
	}

	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel[12]: test\n")
	s.Update(input)
	s.Stop()
}

func makeJSON(numKeys int) []byte {
	keys := make(map[string]string, 0)
	for i := 0; i < numKeys; i++ {
		keys[fmt.Sprintf("%d", i)] = "foo"
	}
	cee, _ := json.Marshal(keys)
	log := fmt.Sprintf("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:%s\n", cee)
	return []byte(log)
}

func BenchmarkSextant(b *testing.B) {
	s, err := captainslog.NewSextant("bench", 0.03, 4)
	if err != nil {
		panic(err)
	}

	m := makeJSON(20)
	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		if err != nil {
			panic(err)
		}
		s.Update(m)
	}

	s.Stop()
}
