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

// NewTagRangeTransformer starts the construction of a TagRangeTransformer.
func NewTagRangeTransformer() *TagRangeTransformer {
	return &TagRangeTransformer{}
}

// Select accepts a MatchType and a matchValue string. If a given SyslogMsg does
// not match the Select critera, than the TagRangeTransformer will return the
// original message without processing it.
func (t *TagRangeTransformer) Select(matchType MatchType, matchValue string) *TagRangeTransformer {
	switch matchType {
	case Program:
		t.selectMatcher = NewTagMatcher(matchValue)
	case Contains:
		t.selectMatcher = NewContentContainsMatcher(matchValue)
	default:
	}
	return t
}

// StartMatch accept a MatchType and a matchValue string. If the SyslogMsg being
// processed matches, then it will be tagged, and every following message that
// matches the Select criteria will be tagged until the first message after
// the EndMatch.
func (t *TagRangeTransformer) StartMatch(matchType MatchType, matchValue string) *TagRangeTransformer {
	switch matchType {
	case Program:
		t.startMatcher = NewTagMatcher(matchValue)
	case Contains:
		t.startMatcher = NewContentContainsMatcher(matchValue)
	default:
	}
	return t
}

// EndMatch accepts a MatchType and a matchValue string. If the SyslogMsg being
// processed matches, then it will be tagged, its key will be removed
// from the tracking db and subsequent messages that match the Select
// will not be tagged.
func (t *TagRangeTransformer) EndMatch(matchType MatchType, matchValue string) *TagRangeTransformer {
	switch matchType {
	case Program:
		t.endMatcher = NewTagMatcher(matchValue)
	case Contains:
		t.endMatcher = NewContentContainsMatcher(matchValue)
	default:
	}
	return t
}

// WaitDuration sets the ammount of time the TagRangeTransformer
// will wait to see an EndMatch after seeing the StartMatch.
func (t *TagRangeTransformer) WaitDuration(duration time.Duration) *TagRangeTransformer {
	t.ttl = duration
	t.reapInterval = t.ttl / 2
	return t
}

// AddTag specifies the tag to be added.
func (t *TagRangeTransformer) AddTag(key string, value string) *TagRangeTransformer {
	t.tagger = NewTagArrayMutator(key, value)
	return t
}

// Do starts the TagRangeTransformer.
func (t *TagRangeTransformer) Do() (*TagRangeTransformer, error) {
	if t.selectMatcher == nil ||
		t.startMatcher == nil ||
		t.endMatcher == nil ||
		t.tagger == nil {
		return nil, fmt.Errorf("argument error")
	}

	t.mutex = &sync.Mutex{}
	t.trackingDB = make(map[string]time.Time)

	go func() {
		for {
			time.Sleep(t.reapInterval)
			t.reap()
		}
	}()

	return t, nil
}

// reap reaps expired keys from the trackingDB.
func (t *TagRangeTransformer) reap() {
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
func (t *TagRangeTransformer) Transform(msg SyslogMsg) (SyslogMsg, error) {
	var err error

	if !t.selectMatcher.Match(&msg) {
		return msg, err
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	logID := getMsgID(&msg)
	var tagIt bool

	if _, ok := t.trackingDB[logID]; ok {
		tagIt = true
		if t.endMatcher.Match(&msg) {
			delete(t.trackingDB, logID)
		}
	} else {
		if t.startMatcher.Match(&msg) {
			tagIt = true
			t.trackingDB[logID] = time.Now()
		}
	}

	if !tagIt {
		return msg, err
	}

	err = t.tagger.Mutate(&msg)
	return msg, err
}
