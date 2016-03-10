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

// TagArrayMutator is a Mutator that modifies a syslog message
// by adding a tag value to an array of tags.
type TagArrayMutator struct {
	tagKey   string
	tagValue string
}

// NewTagArrayMutator constructs a new TagArrayMutator from a supplied
// key and value. The tag will be added to the array at the key if
// it exists - if the key does not exist in the SyslogMsg's JSONValues
// map, it will be  created.
func NewTagArrayMutator(tagKey, tagValue string) Mutator {
	return &TagArrayMutator{
		tagKey:   tagKey,
		tagValue: tagValue,
	}
}

// Mutate modifies the SyslogMsg passed by reference
func (t *TagArrayMutator) Mutate(msg *SyslogMsg) error {
	var err error

	if _, ok := msg.JSONValues[t.tagKey]; !ok {
		msg.JSONValues[t.tagKey] = make([]interface{}, 0)
	}

	switch val := msg.JSONValues[t.tagKey].(type) {
	case []interface{}:
		msg.JSONValues[t.tagKey] = append(val, t.tagValue)
		return err
	default:
		err = fmt.Errorf("tags key in message was not an array")
		return err
	}
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

// TagRangeTransformer is a Transformer implementation that tags
// log lines that meet a selection criteria and are logged
// between a start and end match. Matches are performed by
// implementations of the Matcher interface.
type TagRangeTransformer struct {
	selectMatcher Matcher
	startMatcher  Matcher
	endMatcher    Matcher
	tagger        Mutator
	trackingDB    map[string]time.Time
	ttl           time.Duration
	reapInterval  time.Duration
	mutex         *sync.Mutex
}

// NewTagRangeTransformer accepts a Matcher to select which logs lines the
// mutator should scan, and a start and end Matcher to denote the lines
// that designate the start and end of a match. All lines that match
// the selection criteria and are either the start and end match or
// log lines in between them will be tagged.
func NewTagRangeTransformer(selectMatcher, startMatcher, endMatcher Matcher,
	tagger Mutator, ttlSeconds, reapIntervalSeconds int) *TagRangeTransformer {
	tr := &TagRangeTransformer{
		selectMatcher: selectMatcher,
		startMatcher:  startMatcher,
		endMatcher:    endMatcher,
		tagger:        tagger,
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
func (m *TagRangeTransformer) reap() {
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
func (m *TagRangeTransformer) Transform(msg SyslogMsg) (SyslogMsg, error) {
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

	err = m.tagger.Mutate(&msg)
	return msg, err
}
