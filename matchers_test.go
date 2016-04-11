package captainslog

import "testing"

func TestTagMatcher(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	msg, err := NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	goodMatcher := NewTagMatcher("test:")
	if want, got := true, goodMatcher.Match(&msg); want != got {
		t.Errorf("want %t, got %t", want, got)
	}

	badMatcher := NewTagMatcher("nope:")
	if want, got := false, badMatcher.Match(&msg); want != got {
		t.Errorf("want %t, got %t", want, got)
	}

	prefixMatcher := NewTagMatcher("tes")
	if want, got := true, prefixMatcher.Match(&msg); want != got {
		t.Errorf("want %t, got %t", want, got)
	}
}

func TestContentContainsMatcher(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")
	msg, err := NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	goodMatcher := NewContentContainsMatcher("hello")
	if want, got := true, goodMatcher.Match(&msg); want != got {
		t.Errorf("want %t, got %t", want, got)
	}

	badMatcher := NewContentContainsMatcher("goodbye")
	if want, got := false, badMatcher.Match(&msg); want != got {
		t.Errorf("want %t, got %t", want, got)
	}
}
