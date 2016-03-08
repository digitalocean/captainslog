package captainslog

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// getMsgID creates a key for a log line from
// it's hostname and program name tag
func getMsgID(msg *SyslogMsg) string {
	return fmt.Sprintf("%s!%s", msg.Host, msg.Tag)
}

// TagMatcher is a matcher that matches on RFC3164 tag
type TagMatcher struct {
	match string
}

// Match applies the match to a SyslogMsg
func (t *TagMatcher) Match(msg *SyslogMsg) bool {
	if msg.Tag == t.match {
		return true
	}
	return false
}

// NewTagMatcher creates a TagMatcher that tries
// to match on the supplied string
func NewTagMatcher(match string) Matcher {
	return &TagMatcher{match: match}
}

// ContentContainsMatcher matches when the Content
// field of a SyslogMsg contains the supplied string
type ContentContainsMatcher struct {
	match string
}

// Match applies the match to a SyslogMsg
func (c *ContentContainsMatcher) Match(msg *SyslogMsg) bool {
	return strings.Contains(msg.Content, c.match)
}

// NewContentContainsMatcher creates a ContentContainsMatcher that
// tries to match on the supplied string
func NewContentContainsMatcher(match string) Matcher {
	return &ContentContainsMatcher{match: match}
}

// TagRangeMutator is a Mutator implementation that tags
// log lines that meet a selection criteria and are logged
// between a start and end match. Matches are performed by
// implementations of the Matcher interface.
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

// NewTagRangeMutator accepts a Matcher to select which logs lines the
// mutator should scan, and a start and end Matcher to denote the lines
// that designate the start and end of a match. All lines that match
// the selection criteria and are either the start and end match or
// log lines in between them will be tagged. The tag is designated by
// the tagValue argument. The tag array will exist at the key tagKey.
// ttlSeconds denotes how long a host / program name combination should
// be watched for the end match after a start match. reapIntervalSeconds
// designates how often the reaper routine that checks for expired
// matches should run.
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

// reap reaps expired keys from the trackingDB
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

// Mutate accepts a SyslogMsg and applies the Mutator to it.
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
