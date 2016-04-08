package captainslog

import (
	"sync"
	"time"
)

// TimeSinceTransformer is a transformer implementation that adds a "since" tag with
// pointing to duration in seconds since the last time a log line that matched the
// selectors was seen.
type TimeSinceTransformer struct {
	selectMatchers []Matcher
	trackingDB     map[string]time.Time
	ttl            time.Duration
	reapInterval   time.Duration
	mutex          *sync.Mutex
}

// NewTimeSinceTransformer creates a new TimeSinceTransformer.
func NewTimeSinceTransformer(waitTime time.Duration, selecters ...Matcher) *TimeSinceTransformer {
	t := &TimeSinceTransformer{
		selectMatchers: selecters,
		ttl:            waitTime * time.Second,
		trackingDB:     make(map[string]time.Time),
		reapInterval:   waitTime / 2,
		mutex:          &sync.Mutex{},
	}

	go func() {
		for {
			time.Sleep(t.reapInterval)
			t.reap()
		}
	}()

	return t
}

// reap reaps expired keys from the trackingDB.
func (t *TimeSinceTransformer) reap() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for k, v := range t.trackingDB {
		duration := time.Since(v)
		if duration.Seconds() > t.ttl.Seconds() {
			delete(t.trackingDB, k)
		}
	}
}

// Transform accepts a SyslogMsg and applies the Transformer to it.
func (t *TimeSinceTransformer) Transform(msg SyslogMsg) (SyslogMsg, error) {
	var err error

	for _, m := range t.selectMatchers {
		if !m.Match(&msg) {
			return msg, nil
		}
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	logID := getMsgID(&msg)

	if _, ok := t.trackingDB[logID]; !ok {
		t.trackingDB[logID] = time.Now()
	}

	duration := time.Since(t.trackingDB[logID])
	t.trackingDB[logID] = time.Now()
	msg.AddTag("since", duration.Seconds())
	return msg, err
}
