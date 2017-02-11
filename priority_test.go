package captainslog_test

import (
	"testing"

	"github.com/digitalocean/captainslog"
)

func TestNewPriority(t *testing.T) {
	_, err := captainslog.NewPriority(captainslog.Local0, captainslog.Err)
	if err != nil {
		t.Error(err)
	}
}

func TestNewPriorityBadFacility(t *testing.T) {
	_, err := captainslog.NewPriority(captainslog.Facility(30), captainslog.Err)
	if want, got := captainslog.ErrBadFacility, err; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestNewPriorityBadSeverity(t *testing.T) {
	_, err := captainslog.NewPriority(captainslog.Local0, captainslog.Severity(50))
	if want, got := captainslog.ErrBadSeverity, err; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestFacilityToString(t *testing.T) {
	var f captainslog.Facility

	f = captainslog.Kern
	if want, got := "kern", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.User
	if want, got := "user", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Mail
	if want, got := "mail", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Daemon
	if want, got := "daemon", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Auth
	if want, got := "auth", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Syslog
	if want, got := "syslog", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.LPR
	if want, got := "lpr", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.News
	if want, got := "news", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.UUCP
	if want, got := "uucp", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Cron
	if want, got := "cron", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.AuthPriv
	if want, got := "authpriv", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.FTP
	if want, got := "ftp", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Local0
	if want, got := "local0", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Local1
	if want, got := "local1", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Local2
	if want, got := "local2", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Local3
	if want, got := "local3", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Local4
	if want, got := "local4", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Local5
	if want, got := "local5", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Local6
	if want, got := "local6", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = captainslog.Local7
	if want, got := "local7", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestSeverityToString(t *testing.T) {
	var s captainslog.Severity

	s = captainslog.Emerg
	if want, got := "emerg", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = captainslog.Alert
	if want, got := "alert", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = captainslog.Crit
	if want, got := "crit", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = captainslog.Err
	if want, got := "err", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = captainslog.Warning
	if want, got := "warning", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = captainslog.Notice
	if want, got := "notice", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = captainslog.Info
	if want, got := "info", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = captainslog.Debug
	if want, got := "debug", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
