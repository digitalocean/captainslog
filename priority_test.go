package captainslog

import "testing"

func TestNewPriority(t *testing.T) {
	_, err := NewPriority(Local0, Err)
	if err != nil {
		t.Error(err)
	}
}

func TestNewPriorityBadFacility(t *testing.T) {
	_, err := NewPriority(Facility(30), Err)
	if want, got := ErrBadFacility, err; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestNewPriorityBadSeverity(t *testing.T) {
	_, err := NewPriority(Local0, Severity(50))
	if want, got := ErrBadSeverity, err; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestFacilityToString(t *testing.T) {
	var f Facility

	f = Kern
	if want, got := "kern", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = User
	if want, got := "user", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Mail
	if want, got := "mail", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Daemon
	if want, got := "daemon", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Auth
	if want, got := "auth", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Syslog
	if want, got := "syslog", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = LPR
	if want, got := "lpr", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = News
	if want, got := "news", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = UUCP
	if want, got := "uucp", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Cron
	if want, got := "cron", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = AuthPriv
	if want, got := "authpriv", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = FTP
	if want, got := "ftp", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Local0
	if want, got := "local0", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Local1
	if want, got := "local1", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Local2
	if want, got := "local2", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Local3
	if want, got := "local3", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Local4
	if want, got := "local4", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Local5
	if want, got := "local5", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Local6
	if want, got := "local6", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	f = Local7
	if want, got := "local7", f.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestSeverityToString(t *testing.T) {
	var s Severity

	s = Emerg
	if want, got := "emerg", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = Alert
	if want, got := "alert", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = Crit
	if want, got := "crit", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = Err
	if want, got := "err", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = Warning
	if want, got := "warning", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = Notice
	if want, got := "notice", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = Info
	if want, got := "info", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	s = Debug
	if want, got := "debug", s.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
