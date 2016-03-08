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

type TagMatcher struct {
	match string
}

func (t *TagMatcher) Match(msg *SyslogMsg) bool {
	if msg.Tag == t.match {
		return true
	}
	return false
}

func NewTagMatcher(match string) Matcher {
	return &TagMatcher{match: match}
}

type ContentContainsMatcher struct {
	match string
}

func (c *ContentContainsMatcher) Match(msg *SyslogMsg) bool {
	return strings.Contains(msg.Content, c.match)
}

func NewContentContainsMatcher(match string) Matcher {
	return &ContentContainsMatcher{match: match}
}

type TagRangeMutator struct {
	selectMatcher Matcher
	startMatcher  Matcher
	endMatcher    Matcher
	tagKey        string
	tagValue      string
	trackingDB    map[string]time.Time
	ttl           time.Duration
	reapInterval  time.Duration
	mutex         *sync.Mutex
}

func NewTagRangeMutator(selectMatcher, startMatcher, endMatcher Matcher,
	tagKey, tagValue string, ttlSeconds, reapIntervalSeconds int) *TagRangeMutator {
	tr := &TagRangeMutator{
		selectMatcher: selectMatcher,
		startMatcher:  startMatcher,
		endMatcher:    endMatcher,
		tagKey:        tagKey,
		tagValue:      tagValue,
		trackingDB:    make(map[string]time.Time),
		ttl:           time.Duration(ttlSeconds) * time.Second,
		reapInterval:  time.Duration(reapIntervalSeconds) * time.Second,
		mutex:         &sync.Mutex{},
	}

	go func() {
		for {
			time.Sleep(tr.reapInterval)
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
