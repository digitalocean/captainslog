package captainslog

import "github.com/prometheus/client_golang/prometheus"

type Sextant struct {
	logLinesTotal   prometheus.Counter
	bytesTotal      prometheus.Counter
	parseErrorTotal prometheus.Counter
	jsonLogsTotal   prometheus.Counter
}

func (s *Sextant) register() {
	prometheus.MustRegister(s.bytesTotal)
	prometheus.MustRegister(s.logLinesTotal)
	prometheus.MustRegister(s.parseErrorTotal)
	prometheus.MustRegister(s.jsonLogsTotal)
}

func (s *Sextant) Update(msg *SyslogMsg) {
	s.bytesTotal.Add(float64(len(msg.buf)))
	s.logLinesTotal.Inc()

	if msg.errored {
		s.parseErrorTotal.Inc()
	}

	if msg.IsCee {
		s.jsonLogsTotal.Inc()
	}
}

func NewSextant(namespace string) *Sextant {

	sextant := &Sextant{

		bytesTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "bytes_total",
				Help:      "total bytes read",
			},
		),

		logLinesTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "log_lines_total",
				Help:      "total logs",
			},
		),

		parseErrorTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "parse_errors_total",
				Help:      "total parse errors",
			},
		),

		jsonLogsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "json_logs_total",
				Help:      "total logs that were json",
			},
		),
	}

	sextant.register()
	return sextant
}
