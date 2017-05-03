package captainslog

// TODO: finish snapshotting

import (
	"bufio"
	"bytes"
	"io"
	"sync"
	"time"
)

// Sextant tracks metrics derviced from logs.
type Sextant struct {
	msgChan       chan []byte
	estimator     *Estimator
	estimatorChan chan *Estimator
	workers       []*worker
	stats         *Stats
	ticker        *time.Ticker
	quitChan      chan struct{}
}

// NewSextant returns a new Sextant.
func NewSextant(namespace string, errorRate float64, numWorkers int) (*Sextant, error) {
	estimatorChan := make(chan *Estimator)

	s := &Sextant{
		msgChan:       make(chan []byte),
		estimator:     &Estimator{},
		estimatorChan: estimatorChan,
		workers:       make([]*worker, numWorkers),
		stats:         NewStats(namespace),
		ticker:        time.NewTicker(5 * time.Second),
		quitChan:      make(chan struct{}),
	}

	var err error
	s.estimator, err = NewEstimator(errorRate)
	if err != nil {
		return s, err
	}

	err = s.estimator.AddCounter(&JSONKeyExtractor{}, s.stats.UniqueKeysTotal)

	if err != nil {
		return s, err
	}

	for i := 0; i < numWorkers; i++ {
		var err error
		s.workers[i], err = newWorker(s.stats, s.msgChan, estimatorChan, errorRate)
		if err != nil {
			return s, err
		}
		s.workers[i].start()
	}

	go func() {
		for e := range s.estimatorChan {
			s.estimator.Union(e)
		}
		s.estimator.UpdateCounters()
	}()

	go func() {
		var signal struct{}
		for {
			select {
			case <-s.ticker.C:
				for _, w := range s.workers {
					w.snapChan <- signal
				}
			case <-s.quitChan:
				s.ticker.Stop()
				return
			}
		}
	}()

	return s, err
}

// Update derives metrics from a log passed to it.
func (s *Sextant) Update(b []byte) {
	s.msgChan <- b
}

// Stop shuts down the Sextant.
func (s *Sextant) Stop() {
	close(s.msgChan)
	close(s.quitChan)
	s.stats.Unregister()
}

type worker struct {
	stats         *Stats
	estimator     *Estimator
	mutex         *sync.Mutex
	snapChan      chan struct{}
	msgChan       <-chan []byte
	estimatorChan chan<- *Estimator
}

func newWorker(stats *Stats, msgChan chan []byte, estimatorChan chan *Estimator, errorRate float64) (*worker, error) {

	w := &worker{
		mutex:         &sync.Mutex{},
		msgChan:       msgChan,
		stats:         stats,
		snapChan:      make(chan struct{}),
		estimatorChan: estimatorChan,
		estimator:     &Estimator{},
	}

	var err error
	w.estimator, err = NewEstimator(errorRate)
	return w, err
}

func (w *worker) start() {
	go func() {
		for _ = range w.snapChan {
			w.estimatorChan <- w.estimator
		}
	}()

	go func() {
		for b := range w.msgChan {
			reader := bufio.NewReader(bytes.NewBuffer(b))
			var err error
			for err == nil {
				var line []byte
				line, err = reader.ReadBytes('\n')
				if err != nil || err == io.EOF {
					msg, _ := NewSyslogMsgFromBytes(line)
					w.update(&msg)
				}
			}
		}
	}()
}

func (w *worker) update(msg *SyslogMsg) {
	w.stats.BytesTotal.Add(float64(len(msg.buf)))
	w.stats.LogLinesTotal.Inc()

	if msg.errored {
		w.stats.ParseErrorTotal.Inc()
	}

	if msg.IsCee {
		w.stats.JSONLogsTotal.Inc()
	}

	w.estimator.Estimate(msg)
}
