package captainslog

import "testing"

func TestTimeSinceTransformerTransform(t *testing.T) {
	cases := []testCase{
		testCase{
			input:  []byte("<182>2016-04-07T19:31:49.643625+00:00 test.example.com rsyslogd-pstats: @cee:{\"inblock\":0,\"majflt\":0,\"maxrss\":5020,\"minflt\":624,\"name\":\"resource-usage\",\"nivcsw\":14,\"nvcsw\":1358,\"origin\":\"impstats\",\"oublock\":136,\"stime\":59937,\"utime\":174816}\n"),
			result: true,
		},
	}

	transformer := NewTimeSinceTransformer(
		"since", 86400,
		NewContentContainsMatcher("inblock"),
		NewTagMatcher("rsyslogd-pstats:"),
	)

	for i, v := range cases {
		original := NewSyslogMsg()
		err := Unmarshal(v.input, &original)
		if err != nil {
			t.Error(err)
		}

		mutated, err := transformer.Transform(original)
		if err != nil {
			t.Error(err)
		}

		_, hasTagsKey := mutated.JSONValues["since"]

		if want, got := v.result, hasTagsKey; want != got {
			t.Errorf("case %d: want '%v', got '%v'", i, want, got)
		}
	}
}
