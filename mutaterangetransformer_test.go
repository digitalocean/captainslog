package captainslog

import (
	"testing"
	"time"
)

type testCase struct {
	input  []byte
	result bool
}

func TestRangeTransformerTTL(t *testing.T) {
	input := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: [15803005.789011] ------------[ cut here ]------------\n")

	transformer, err := NewTagRangeTransformer().
		Select(Program, "kernel:").
		StartMatch(Contains, "[ cut here ]").
		EndMatch(Contains, "[ end trace").
		AddTag("tags", "trace").
		WaitDuration(1).
		Do()

	if err != nil {
		t.Error(err)
	}

	original := NewSyslogMsg()
	err = Unmarshal(input, &original)
	if err != nil {
		t.Error(err)
	}

	transformer.mutex.Lock()
	if want, got := 0, len(transformer.trackingDB); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
	transformer.mutex.Unlock()

	_, err = transformer.Transform(original)
	if err != nil {
		t.Error(err)
	}

	transformer.mutex.Lock()
	if want, got := 1, len(transformer.trackingDB); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
	transformer.mutex.Unlock()

	time.Sleep(2 * time.Second)

	transformer.mutex.Lock()
	if want, got := 0, len(transformer.trackingDB); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
	transformer.mutex.Unlock()
}

func TestTagRangeTransformerTransform(t *testing.T) {
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

	transformer, err := NewTagRangeTransformer().
		Select(Program, "kernel:").
		StartMatch(Contains, "[ cut here ]").
		EndMatch(Contains, "[ end trace").
		AddTag("tags", "trace").
		WaitDuration(60).
		Do()

	if err != nil {
		t.Error(err)
	}

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

		_, hasTagsKey := mutated.JSONValues["tags"]

		if want, got := v.result, hasTagsKey; want != got {
			t.Errorf("case %d: want '%v', got '%v'", i, want, got)
		}
	}
}

func BenchmarkTagRangeTransformerTransform(b *testing.B) {
	m := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: [15803005.789011] ------------[ cut here ]------------\n")

	transformer, err := NewTagRangeTransformer().
		Select(Program, "kernel:").
		StartMatch(Contains, "[ cut here ]").
		EndMatch(Contains, "[ end trace").
		AddTag("tags", "trace").
		WaitDuration(60).
		Do()

	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		original := NewSyslogMsg()
		err := Unmarshal(m, &original)
		if err != nil {
			panic(err)
		}

		_, err = transformer.Transform(original)
		if err != nil {
			panic(err)
		}
	}
}
