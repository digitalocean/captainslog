package captainslog

import (
	"fmt"
	"sync"
	"time"
)

// getMsgID creates a key for a log line from
// it's hostname and program name tag.
func getMsgID(msg *SyslogMsg) string {
	return fmt.Sprintf("%s!%s", msg.Host, msg.Tag)
}

// MutateRangeTransformer is a Transformer implementation that
// mutates log lines that meet a selection criteria and are logged
// between a start and end match. Matches are performed by
// implementations of the Matcher interface.
type MutateRangeTransformer struct {
	selectMatcher Matcher
	startMatcher  Matcher
	endMatcher    Matcher
	mutator       Mutator
	trackingDB    map[string]time.Time
	ttl           time.Duration
	reapInterval  time.Duration
	mutex         *sync.Mutex
}

// NewMutateRangeTransformer creates a new MutateRangeTransformer.
func NewMutateRangeTransformer(selecter, starter, ender Matcher, mutator Mutator, waitTime time.Duration) *MutateRangeTransformer {
	m := &MutateRangeTransformer{
		selectMatcher: selecter,
		startMatcher:  starter,
		endMatcher:    ender,
		mutator:       mutator,
		ttl:           waitTime,
		trackingDB:    make(map[string]time.Time),
		reapInterval:  waitTime / 2,
		mutex:         &sync.Mutex{},
	}

	go func() {
		for {
			time.Sleep(m.reapInterval)
			m.reap()
		}
	}()

	return m
}

// reap reaps expired keys from the trackingDB.
func (m *MutateRangeTransformer) reap() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k, v := range m.trackingDB {
		duration := time.Since(v)
		if duration.Seconds() > m.ttl.Seconds() {
			delete(m.trackingDB, k)
		}
	}
}

// Transform accepts a SyslogMsg and applies the Transformer to it.
func (m *MutateRangeTransformer) Transform(msg SyslogMsg) (SyslogMsg, error) {
	var err error

	if !m.selectMatcher.Match(&msg) {
		return msg, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	logID := getMsgID(&msg)
	var tagIt bool

	if _, ok := m.trackingDB[logID]; ok {
		tagIt = true
		if m.endMatcher.Match(&msg) {
			delete(m.trackingDB, logID)
		}
	} else {
		if m.startMatcher.Match(&msg) {
			tagIt = true
			m.trackingDB[logID] = time.Now()
		}
	}

	if !tagIt {
		return msg, err
	}

	err = m.mutator.Mutate(&msg)
	return msg, err
}
