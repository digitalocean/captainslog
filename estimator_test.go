package captainslog_test

import (
	"testing"

	"github.com/digitalocean/captainslog"
)

func TestEstimator(t *testing.T) {
	e1, err := captainslog.NewEstimator(0.1)
	if err != nil {
		t.Error(err)
	}

	e2, err := captainslog.NewEstimator(0.1)
	if err != nil {
		t.Error(err)
	}

	e2.Add("hello")
	if want, got := 1, int(e2.Cardinality()); want != got {
		t.Errorf("want %d, got %d", want, got)
	}

	e1.Union(e2)
	if want, got := 1, int(e1.Cardinality()); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func BenchmarkEstimator(b *testing.B) {
	e, err := captainslog.NewEstimator(0.1)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		e.Add("hello")
	}
}
