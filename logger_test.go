package captainslog

import "testing"

func TestNewMostlyFeaturelessLogger(t *testing.T) {
	_, err := NewMostlyFeaturelessLogger(Local7)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateLogMessage(t *testing.T) {
	msg, err := createLogMessage(Fields{"hello": "world"})

	if err != nil {
		t.Error(err)
	}

	if want, got := "@cee:{\"hello\":\"world\"}", msg; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}
