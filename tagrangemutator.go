package captainslog

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func getMsgID(msg *SyslogMsg) string {
	return fmt.Sprintf("%s!%s", msg.Host, msg.Tag)
}

func TagMatcher(tag string) func(msg *SyslogMsg) bool {
	return func(msg *SyslogMsg) bool {
		if msg.Tag == tag {
			return true
		}
		return false
	}
}

func ContentContainsMatcher(contains string) func(msg *SyslogMsg) bool {
	return func(msg *SyslogMsg) bool {
		return strings.Contains(msg.Content, contains)
	}
}

type TagRangeMutator struct {
	selectMatcher func(msg *SyslogMsg) bool
	startMatcher  func(msg *SyslogMsg) bool
	endMatcher    func(msg *SyslogMsg) bool
	tagKey        string
	tagValue      string
	trackingDB    map[string]time.Time
	ttl           time.Duration
	reapInterval  time.Duration
	mutex         *sync.Mutex
}

func NewTagRangeMutator(selectMatcher, startMatcher, endMatcher func(msg *SyslogMsg) bool,
	tagKey, tagValue string) *TagRangeMutator {
	tr := &TagRangeMutator{
		selectMatcher: selectMatcher,
		startMatcher:  startMatcher,
		endMatcher:    endMatcher,
		tagKey:        tagKey,
		tagValue:      tagValue,
		trackingDB:    make(map[string]time.Time),
		ttl:           60 * time.Second,
		reapInterval:  10 * time.Second,
		mutex:         &sync.Mutex{},
	}

	go func() {
		for {
			time.Sleep(time.Second * tr.reapInterval)
			tr.reap()
		}
	}()

	return tr
}

func (m *TagRangeMutator) reap() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k, v := range m.trackingDB {
		duration := time.Since(v)
		if duration.Minutes() > m.ttl.Minutes() {
			delete(m.trackingDB, k)
		}
	}
}

func (m *TagRangeMutator) Mutate(msg SyslogMsg) (SyslogMsg, error) {
	var err error

	if !m.selectMatcher(&msg) {
		return msg, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	logID := getMsgID(&msg)
	var tagIt bool

	if _, ok := m.trackingDB[logID]; ok {
		tagIt = true
		if m.endMatcher(&msg) {
			delete(m.trackingDB, logID)
		}
	} else {
		if m.startMatcher(&msg) {
			tagIt = true
			m.trackingDB[logID] = time.Now()
		}
	}

	if !tagIt {
		return msg, err
	}

	if _, ok := msg.JSONValues[m.tagKey]; !ok {
		msg.JSONValues[m.tagKey] = make([]interface{}, 0)
	}

	switch val := msg.JSONValues[m.tagKey].(type) {
	case []interface{}:
		msg.JSONValues[m.tagKey] = append(val, m.tagValue)
		return msg, err
	default:
		err = fmt.Errorf("tags key in message was not an array")
		return msg, err
	}
}
