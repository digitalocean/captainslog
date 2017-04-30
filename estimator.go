package captainslog

import (
	"github.com/mynameisfiber/gohll"
)

// Estimator holds a set of hyperloglogs that it uses to
// perform cardidnality estimation of keyspaces derived
// from logs.
type Estimator struct {
	KeysHLL     *gohll.HLL
	ProgramsHLL *gohll.HLL
}

// NewEstimator creates a new Estimator.
func NewEstimator(errorRate float64) (*Estimator, error) {
	e := &Estimator{}
	var err error
	e.KeysHLL, err = gohll.NewHLLByError(errorRate)
	if err != nil {
		return e, err
	}
	e.ProgramsHLL, err = gohll.NewHLLByError(errorRate)
	return e, err
}
