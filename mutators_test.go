package captainslog

import "testing"

func TestTagArrayMutator(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n")
	msg, err := NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	mutator := NewTagArrayMutator("tags", "youareit")
	err = mutator.Mutate(&msg)
	if err != nil {
		t.Error(err)
	}

	if _, ok := msg.JSONValues["tags"]; !ok {
		t.Error("expected 'tags' key in ms.JSONValues")
	}

}
