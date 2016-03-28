package captainslog

import (
	"bytes"
	"io"
	"testing"
)

type pipelineTestCase struct {
	testInput  string
	testOutput string
}

func getTestReader(testCases []pipelineTestCase) io.Reader {
	var testInputs []byte
	for _, testCase := range testCases {
		testCase.testInput = testCase.testInput + "\n"
		testInputs = append(testInputs, []byte(testCase.testInput)...)
	}

	testReader := bytes.NewBuffer(testInputs)
	return testReader
}

func getExpectedTestResults(testCases []pipelineTestCase) []string {
	var wants []string
	for _, testCase := range testCases {
		wants = append(wants, testCase.testOutput)
	}
	return wants
}

func TestPipeline(t *testing.T) {
	testCases := []pipelineTestCase{
		{
			"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\",\"one.two.three\":\"four.five.six\"}",
			"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first_name\":\"captain\",\"one_two_three\":\"four.five.six\"}",
		},
		{
			"<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: [15803005.789011] ------------[ cut here ]------------",
			"<4>2016-03-08T14:59:36.293816+00:00 host.example.com kernel: @cee:{\"msg\":\"[15803005.789011] ------------[ cut here ]------------\",\"tags\":[\"trace\"]}",
		},
		{
			"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world",
			"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world",
		},
	}

	testReader := getTestReader(testCases)
	testWriter := bytes.NewBuffer(make([]byte, 0))

	sanitizer, err := NewJSONKeyTransformer().
		OldString(".").
		NewString("_").
		Do()

	if err != nil {
		t.Error(err)
	}

	tagger, err := NewTagRangeTransformer().
		Select(Program, "kernel:").
		StartMatch(Contains, "[ cut here ]").
		EndMatch(Contains, "[ end trace").
		AddTag("tags", "trace").
		WaitDuration(1).
		Do()

	if err != nil {
		t.Error(err)
	}

	err = NewPipeline().
		From(testReader).
		Transform(sanitizer).
		Transform(tagger).
		To(testWriter).
		Do()

	if want, got := io.EOF, err; want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	wants := getExpectedTestResults(testCases)

	for _, want := range wants {
		want = want + "\n"
		got, _ := testWriter.ReadString('\n')
		if want != got {
			t.Errorf("want %q, got %q", want, got)
		}

	}
}
