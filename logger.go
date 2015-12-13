package captainslog

import (
	"encoding/json"
	"fmt"
	"log"
	"log/syslog"
)

// Fields are a map of key value pairs for a
// log line that will be output as JSON
type Fields map[string]interface{}

// MostlyFeaturelessLogger is a mostly featureless logger created for simple
// structured logging of Notice and Err messages from daemons created with
// captainslog to syslog. If you need something more than that this probably
// is not something that will make you happy.
type MostlyFeaturelessLogger struct {
	errorLogger *log.Logger
	infoLogger  *log.Logger
}

// NewMostlyFeaturelessLogger returns a new MostlyFeaturelessLogger for
// the given Facility
func NewMostlyFeaturelessLogger(f Facility) (*MostlyFeaturelessLogger, error) {
	l := &MostlyFeaturelessLogger{}
	var p *Priority
	var err error

	p, err = NewPriority(f, Err)
	if err != nil {
		return l, err
	}

	l.errorLogger, err = syslog.NewLogger(syslog.Priority(p.Priority), 0)
	if err != nil {
		return l, err
	}

	p, err = NewPriority(f, Notice)
	if err != nil {
		return l, err
	}

	l.infoLogger, err = syslog.NewLogger(syslog.Priority(p.Priority), 0)
	if err != nil {
		return l, err
	}

	return l, err
}

func createLogMessage(fields Fields) (string, error) {
	b, err := json.Marshal(fields)
	return fmt.Sprintf("@cee:%s", string(b)), err
}

// ErrorWithFields accepts Fields and logs a @cee structured log
// to syslog at level Err
func (l *MostlyFeaturelessLogger) ErrorWithFields(fields Fields) error {
	var err error
	var payload string
	payload, err = createLogMessage(fields)
	if err != nil {
		return err
	}
	l.errorLogger.Printf("%s", payload)
	return err
}

// InfoWithFields accepts Fields and logs a @cee structured log
// to syslog at level Notice
func (l *MostlyFeaturelessLogger) InfoWithFields(fields Fields) error {
	var err error
	var payload string
	payload, err = createLogMessage(fields)
	if err != nil {
		return err
	}

	l.infoLogger.Printf("%s", payload)
	return err
}
