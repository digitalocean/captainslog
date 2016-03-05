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

func TestFacilityTextToFacility(t *testing.T) {
	facility, err := FacilityTextToFacility("KERN")
	if err != nil {
		t.Error(err)
	}
	if want, got := Kern, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("USER")
	if err != nil {
		t.Error(err)
	}
	if want, got := User, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("MAIL")
	if err != nil {
		t.Error(err)
	}
	if want, got := Mail, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("DAEMON")
	if err != nil {
		t.Error(err)
	}
	if want, got := Daemon, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("AUTH")
	if err != nil {
		t.Error(err)
	}
	if want, got := Auth, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("SYSLOG")
	if err != nil {
		t.Error(err)
	}
	if want, got := Syslog, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("LPR")
	if err != nil {
		t.Error(err)
	}
	if want, got := LPR, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("NEWS")
	if err != nil {
		t.Error(err)
	}
	if want, got := News, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("UUCP")
	if err != nil {
		t.Error(err)
	}
	if want, got := UUCP, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("CRON")
	if err != nil {
		t.Error(err)
	}
	if want, got := Cron, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("AUTHPRIV")
	if err != nil {
		t.Error(err)
	}
	if want, got := AuthPriv, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("FTP")
	if err != nil {
		t.Error(err)
	}
	if want, got := FTP, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("LOCAL0")
	if err != nil {
		t.Error(err)
	}
	if want, got := Local0, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("LOCAL1")
	if err != nil {
		t.Error(err)
	}
	if want, got := Local1, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("LOCAL2")
	if err != nil {
		t.Error(err)
	}
	if want, got := Local2, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("LOCAL3")
	if err != nil {
		t.Error(err)
	}
	if want, got := Local3, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("LOCAL4")
	if err != nil {
		t.Error(err)
	}
	if want, got := Local4, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("LOCAL5")
	if err != nil {
		t.Error(err)
	}
	if want, got := Local5, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("LOCAL6")
	if err != nil {
		t.Error(err)
	}
	if want, got := Local6, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("LOCAL7")
	if err != nil {
		t.Error(err)
	}
	if want, got := Local7, facility; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	facility, err = FacilityTextToFacility("BOGUS")
	if want, got := ErrBadFacility, err; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}
