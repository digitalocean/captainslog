package captainslog

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Stats is a collection of prometheus metrics useful
// for gainin insight into log streams.
type Stats struct {
	Namespace           string
	LogLinesTotal       prometheus.Counter
	BytesTotal          prometheus.Counter
	ParseErrorTotal     prometheus.Counter
	JSONLogsTotal       prometheus.Counter
	UniqueKeysTotal     prometheus.Counter
	UniqueProgramsTotal prometheus.Counter
}

// NewStats returns a new Stats reference.
func NewStats(namespace string) *Stats {
	s := &Stats{
		Namespace: namespace,
		BytesTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "bytes_total",
				Help:      "total bytes read",
			},
		),

		LogLinesTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "log_lines_total",
				Help:      "total logs",
			},
		),

		ParseErrorTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "parse_errors_total",
				Help:      "total parse errors",
			},
		),

		JSONLogsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "json_logs_total",
				Help:      "total logs that were json",
			},
		),

		UniqueKeysTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "unique_keys_total",
				Help:      "unique JSON keys",
			},
		),

		UniqueProgramsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "unique_programs_total",
				Help:      "unique program names keys",
			},
		),
	}

	prometheus.MustRegister(s.BytesTotal)
	prometheus.MustRegister(s.LogLinesTotal)
	prometheus.MustRegister(s.ParseErrorTotal)
	prometheus.MustRegister(s.JSONLogsTotal)
	prometheus.MustRegister(s.UniqueKeysTotal)
	prometheus.MustRegister(s.UniqueProgramsTotal)

	return s
}

// Unregister calls unregister for each prometheus
// counter.
func (s *Stats) Unregister() {
	prometheus.Unregister(s.BytesTotal)
	prometheus.Unregister(s.LogLinesTotal)
	prometheus.Unregister(s.ParseErrorTotal)
	prometheus.Unregister(s.JSONLogsTotal)
	prometheus.Unregister(s.UniqueKeysTotal)
	prometheus.Unregister(s.UniqueProgramsTotal)
}