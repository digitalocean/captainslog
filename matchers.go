package captainslog

import "strings"

// TagMatcher is a matcher that matches on RFC3164 program tag.
type TagMatcher struct {
	match string
}

// Match applies the match to a SyslogMsg. It returns
// true of the tag matches and false if not.
func (t *TagMatcher) Match(msg *SyslogMsg) bool {
	return strings.HasPrefix(msg.Tag, t.match)
}

// NewTagMatcher creates a TagMatcher that tries
// to match on the supplied string.
func NewTagMatcher(match string) Matcher {
	return &TagMatcher{match: match}
}

// ContentContainsMatcher matches when the Content
// field of a SyslogMsg contains the supplied string.
type ContentContainsMatcher struct {
	match string
}

// Match applies the match to a SyslogMsg.
func (c *ContentContainsMatcher) Match(msg *SyslogMsg) bool {
	return strings.Contains(msg.Content, c.match)
}

// NewContentContainsMatcher creates a ContentContainsMatcher that
// tries to match on the supplied string.
func NewContentContainsMatcher(match string) Matcher {
	return &ContentContainsMatcher{match: match}
}
