package captainslog

import (
	"sync"

	"github.com/mynameisfiber/gohll"
)

// Estimator holds a set of hyperloglogs that it uses to
// perform cardidnality estimation of keyspaces derived
// from logs.
type Estimator struct {
	keysHLL *gohll.HLL
	mutex   *sync.Mutex
}

// NewEstimator creates a new Estimator.
func NewEstimator(errorRate float64) (*Estimator, error) {
	e := &Estimator{mutex: &sync.Mutex{}}
	var err error
	e.keysHLL, err = gohll.NewHLLByError(errorRate)
	return e, err
}

// Union unions this Estimator with another Estimator.
func (e *Estimator) Union(other *Estimator) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.keysHLL.Union(other.keysHLL)
}

// Add adds a key.
func (e *Estimator) Add(k string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.keysHLL.AddWithHasher(k, gohll.MMH3Hash)
}

// Cardinality returns the etsimated cardinality.
func (e *Estimator) Cardinality() float64 {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.keysHLL.Cardinality()
}
