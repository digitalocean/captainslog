package captainslog

// TODO: update on worker should cause all workers to send hlls back to sextant for union
// TODO: move hlls to use the new struct

import (
	"bufio"
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/mynameisfiber/gohll"
)

// Sextant tracks metrics derviced from logs.
type Sextant struct {
	msgChan   chan []byte
	hllChan   <-chan *gohll.HLL
	estimator *Estimator
	stats     *Stats
}

// NewSextant returns a new Sextant.
func NewSextant(namespace string, errorRate float64, numWorkers int) (*Sextant, error) {
	hllChan := make(chan *gohll.HLL)

	s := &Sextant{
		msgChan:   make(chan []byte),
		stats:     NewStats(namespace),
		hllChan:   hllChan,
		estimator: &Estimator{},
	}

	var err error
	s.estimator, err = NewEstimator(errorRate)
	if err != nil {
		return s, err
	}

	for i := 0; i < numWorkers; i++ {
		worker, err := newWorker(s.stats, s.msgChan, hllChan, errorRate)
		if err != nil {
			return s, err
		}
		worker.start()
	}

	return s, err
}

// Update dervices metrics from a log passed to it.
func (s *Sextant) Update(b []byte) {
	s.msgChan <- b
}

// Stop shuts down the Sextant.
func (s *Sextant) Stop() {
	close(s.msgChan)
	s.stats.Unregister()
}

type worker struct {
	stats            *Stats
	estimator        *Estimator
	mutex            *sync.Mutex
	previousKeys     float64
	previousPrograms float64
	ticker           *time.Ticker
	quit             chan struct{}
	msgChan          <-chan []byte
	hllChan          chan<- *gohll.HLL
}

func newWorker(stats *Stats, msgChan chan []byte, hllChan chan *gohll.HLL, errorRate float64) (*worker, error) {

	w := &worker{
		mutex:        &sync.Mutex{},
		previousKeys: 0,
		quit:         make(chan struct{}),
		msgChan:      msgChan,
		stats:        stats,
		ticker:       time.NewTicker(5 * time.Second),
		hllChan:      hllChan,
		estimator:    &Estimator{},
	}

	var err error
	w.estimator, err = NewEstimator(errorRate)
	return w, err
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.snapshot()
			case <-w.quit:
				w.ticker.Stop()
				return
			}
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

func (w *worker) stop() {
	close(w.quit)
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

	w.mutex.Lock()
	for k := range msg.JSONValues {
		w.estimator.KeysHLL.AddWithHasher(k, gohll.MMH3Hash)
	}

	w.estimator.ProgramsHLL.AddWithHasher(msg.Program, gohll.MMH3Hash)
	w.mutex.Unlock()
}

func (w *worker) snapshot() {
	w.mutex.Lock()
	keys := w.estimator.KeysHLL.Cardinality() - w.previousKeys
	programs := w.estimator.ProgramsHLL.Cardinality() - w.previousPrograms
	w.mutex.Unlock()

	w.stats.UniqueKeysTotal.Add(keys)
	w.previousKeys = keys

	w.stats.UniqueProgramsTotal.Add(keys)
	w.previousPrograms = programs
}
