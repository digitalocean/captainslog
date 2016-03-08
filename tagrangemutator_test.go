package captainslog

import "testing"

func TestTagRangeMutatorIsMutator(t *testing.T) {
}

type testCase struct {
	input  []byte
	result bool
}

func TestTagRangeMutatorMutate(t *testing.T) {
	cases := []testCase{
		testCase{
			input:  []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: [15803005.789010] this line is not part of the trace\n"),
			result: false,
		},
		testCase{
			input:  []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: [15803005.789011] ------------[ cut here ]------------\n"),
			result: true,
		},
		testCase{
			input:  []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: [15803005.789013] this  line should be part of the trace\n"),
			result: true,
		},
		testCase{
			input:  []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: [15803005.789015] this line should also be part of the trace\n"),
			result: true,
		},
		testCase{
			input:  []byte("<4>2016-03-08T14:59:36.293918+00:00 host.example.com kernel: [15803005.789433] ---[ end trace 999999999 ]---\n"),
			result: true,
		},
		testCase{
			input:  []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: [15803005.789010] this line is not part of the trace\n"),
			result: false,
		},
	}

	mutator := NewTagRangeMutator(
		TagMatcher("kernel:"),
		ContentContainsMatcher("[ cut here ]"),
		ContentContainsMatcher("[ end trace"))

	for i, v := range cases {
		original := NewSyslogMsg()
		err := Unmarshal(v.input, &original)
		if err != nil {
			t.Error(err)
		}

		mutated, err := mutator.Mutate(original)
		if err != nil {
			t.Error(err)
		}

		_, hasTagsKey := mutated.JSONValues["tags"]

		if want, got := v.result, hasTagsKey; want != got {
			t.Errorf("case %d: want '%v', got '%v'", i, want, got)
		}
	}
}
