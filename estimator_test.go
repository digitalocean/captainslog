package captainslog_test

import (
	"testing"

	"github.com/digitalocean/captainslog"
	"github.com/prometheus/client_golang/prometheus"
)

func TestEstimator(t *testing.T) {

	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "test",
			Name:      "test_unique_counts",
			Help:      "test counter",
		},
	)

	e1, err := captainslog.NewEstimator(0.1)
	if err != nil {
		t.Error(err)
	}

	err = e1.AddCounter(&captainslog.JSONKeyExtractor{}, counter)
	if err != nil {
		t.Error(err)
	}

	e2, err := captainslog.NewEstimator(0.1)
	if err != nil {
		t.Error(err)
	}

	err = e2.AddCounter(&captainslog.JSONKeyExtractor{}, counter)
	if err != nil {
		t.Error(err)
	}

}
