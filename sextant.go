package captainslog

import (
	"sync"
	"time"

	"github.com/mynameisfiber/gohll"
	"github.com/prometheus/client_golang/prometheus"
)

// Sextant derives prometheus metrics from syslog logs.
type Sextant struct {
	LogLinesTotal       prometheus.Counter
	BytesTotal          prometheus.Counter
	ParseErrorTotal     prometheus.Counter
	JSONLogsTotal       prometheus.Counter
	UniqueKeysTotal     prometheus.Counter
	UniqueProgramsTotal prometheus.Counter
	keysHLL             *gohll.HLL
	programsHLL         *gohll.HLL
	mutex               *sync.Mutex
	previousKeys        float64
	previousPrograms    float64
	Quit                chan struct{}
	ticker              *time.Ticker
}

func (s *Sextant) register() {
	prometheus.MustRegister(s.BytesTotal)
	prometheus.MustRegister(s.LogLinesTotal)
	prometheus.MustRegister(s.ParseErrorTotal)
	prometheus.MustRegister(s.JSONLogsTotal)
	prometheus.MustRegister(s.UniqueKeysTotal)
	prometheus.MustRegister(s.UniqueProgramsTotal)
}

// Start starts the Sextant's goroutine that periodically
// updates prometheus counters from the HLLs.
func (s *Sextant) Start() {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.updatePrometheus()
			case <-s.Quit:
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the Sextant's goroutine.
func (s *Sextant) Stop() {
	close(s.Quit)
	prometheus.Unregister(s.BytesTotal)
	prometheus.Unregister(s.LogLinesTotal)
	prometheus.Unregister(s.ParseErrorTotal)
	prometheus.Unregister(s.JSONLogsTotal)
	prometheus.Unregister(s.UniqueKeysTotal)
	prometheus.Unregister(s.UniqueProgramsTotal)
}

// Update updates the prometheus counters within Sextant based on
// data from the passed in captainslog.SyslogMsg reference
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

	for k := range msg.JSONValues {
		s.keysHLL.AddWithHasher(k, gohll.MMH3Hash)
	}

	s.programsHLL.AddWithHasher(msg.Program, gohll.MMH3Hash)
	s.mutex.Unlock()
}

func (s *Sextant) updatePrometheus() {
	s.mutex.Lock()

	keys := s.keysHLL.Cardinality() - s.previousKeys
	s.UniqueKeysTotal.Add(keys)
	s.previousKeys = keys

	programs := s.programsHLL.Cardinality() - s.previousPrograms
	s.UniqueProgramsTotal.Add(keys)
	s.previousPrograms = programs

	s.mutex.Unlock()
}

// NewSextant creates a new instance of a Sextant, which provides
// prometheus metrics from syslog logs.
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

		UniqueProgramsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "unique_programs_total",
				Help:      "unique program names keys",
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
	sextant.programsHLL, err = gohll.NewHLLByError(0.05)
	return sextant, err
}
