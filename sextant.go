package captainslog

import (
	"sync"
	"time"

	"github.com/mynameisfiber/gohll"
	"github.com/prometheus/client_golang/prometheus"
)

type Sextant struct {
	LogLinesTotal   prometheus.Counter
	BytesTotal      prometheus.Counter
	ParseErrorTotal prometheus.Counter
	JSONLogsTotal   prometheus.Counter
	UniqueKeysTotal prometheus.Counter
	keysHLL         *gohll.HLL
	mutex           *sync.Mutex
	previousKeys    float64
	Quit            chan struct{}
	ticker          *time.Ticker
}

func (s *Sextant) register() {
	prometheus.MustRegister(s.BytesTotal)
	prometheus.MustRegister(s.LogLinesTotal)
	prometheus.MustRegister(s.ParseErrorTotal)
	prometheus.MustRegister(s.JSONLogsTotal)
	prometheus.MustRegister(s.UniqueKeysTotal)
}

func (s *Sextant) Update(msg *SyslogMsg) {
	s.mutex.Lock()
	s.BytesTotal.Add(float64(len(msg.buf)))
	s.LogLinesTotal.Inc()

	if msg.errored {
		s.ParseErrorTotal.Inc()
	}

	if msg.IsCee {
		s.JSONLogsTotal.Inc()
	}

	for k, _ := range msg.JSONValues {
		s.keysHLL.AddWithHasher(k, gohll.MMH3Hash)
	}
	s.mutex.Unlock()
}

func (s *Sextant) updatePrometheus() {
	s.mutex.Lock()
	keys := s.keysHLL.Cardinality() - s.previousKeys
	s.UniqueKeysTotal.Add(keys)
	s.previousKeys = keys
	s.mutex.Unlock()
}

func NewSextant(namespace string) (*Sextant, error) {

	sextant := &Sextant{

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

		mutex:        &sync.Mutex{},
		previousKeys: 0,
		Quit:         make(chan struct{}),
		ticker:       time.NewTicker(5 * time.Second),
	}

	sextant.register()

	var err error
	sextant.keysHLL, err = gohll.NewHLLByError(0.05)

	go func() {
		for {
			select {
			case <-sextant.ticker.C:
				sextant.updatePrometheus()
			case <-sextant.Quit:
				sextant.ticker.Stop()
				return
			}
		}
	}()
	return sextant, err
}
