package captainslog

import (
	"sync"

	"github.com/mynameisfiber/gohll"
	"github.com/prometheus/client_golang/prometheus"
)

// JSONKeyExtractor is an Extractor that gets all keys
// from a CEE Syslog Message
type JSONKeyExtractor struct{}

// Extract implements the Extractor interface
func (j *JSONKeyExtractor) Extract(msg *SyslogMsg) ([]string, error) {
	var keys []string
	for k := range msg.JSONValues {
		keys = append(keys, k)
	}
	return keys, nil
}

type Extractor interface {
	Extract(msg *SyslogMsg) ([]string, error)
}

// Estimator holds a set of hyperloglogs that it uses to
// perform cardidnality estimation of keyspaces derived
// from logs.
type Estimator struct {
	mutex      *sync.Mutex
	errorRate  float64
	db         []*gohll.HLL
	extractors []Extractor
	counters   []prometheus.Counter
}

// NewEstimator creates a new Estimator.
func NewEstimator(errorRate float64) (*Estimator, error) {
	e := &Estimator{
		mutex:      &sync.Mutex{},
		errorRate:  errorRate,
		db:         make([]*gohll.HLL, 0),
		extractors: make([]Extractor, 0),
		counters:   make([]prometheus.Counter, 0),
	}
	var err error
	return e, err
}

// AddCounter adds another extractor function + prometheus counter to the
// Estimator, and makes a new hll to hold extracted keys.
func (e *Estimator) AddCounter(ex Extractor, c prometheus.Counter) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	hll, err := gohll.NewHLLByError(e.errorRate)
	if err != nil {
		return err
	}

	e.db = append(e.db, hll)
	e.extractors = append(e.extractors, ex)
	e.counters = append(e.counters, c)
	return err
}

// Estimate ranges over the hlls and updates them by applying each hlls
// associated extractor function to the SyslogMsg.
func (e *Estimator) Estimate(msg *SyslogMsg) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for i, _ := range e.db {
		keys, err := e.extractors[i].Extract(msg)
		if err != nil {
			continue
		}
		for _, k := range keys {
			e.db[i].Add(k)
		}
	}
}

// UpdateCounters updates the prometheus counters from the hll data for each hll.
// We don't do this every time Estimate is called because retrieving the
// cardinality is much more expensive than adding keys.
func (e *Estimator) UpdateCounters() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for i, _ := range e.db {
		e.counters[i].Add(e.db[i].Cardinality())
	}

}

// Union unions this Estimator with another Estimator.
func (e *Estimator) Union(other *Estimator) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for i, o := range other.db {
		e.db[i].Union(o)
	}
}
