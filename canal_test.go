package captainslog

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

func TestCanal(t *testing.T) {
	cases := []struct{ inMsg, outMsg string }{
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

	retryInterval := 5
	inAddr := "127.0.0.1:45454"
	outAddr := "127.0.0.1:45455"

	outputter := NewTCPOutputter(outAddr, retryInterval)
	outputChanneler := NewOutputChanneler(outputter)

	inputter, _ := NewTCPInputter(inAddr)
	inputChanneler := NewInputChanneler(inputter)

	replacer := strings.NewReplacer(".", "_")
	sanitizer := NewJSONKeyTransformer(replacer)

	tagger := NewTagRangeTransformer(
		NewTagMatcher("kernel:"),
		NewContentContainsMatcher("[ cut here ]"),
		NewContentContainsMatcher("[ end trace"),
		NewTagArrayMutator("tags", "trace"),
		60, 30)

	canal := NewCanal(inputChanneler, outputChanneler, tagger, sanitizer)

	go func() { canal.Ship() }()

	go func() {
		conn, err := net.Dial("tcp", inAddr)
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()

		for _, testCase := range cases {
			_, err = conn.Write([]byte(testCase.inMsg + "\n"))
			if err != nil {
				t.Error(err)
			}
		}
	}()

	l, err := net.Listen("tcp", outAddr)
	if err != nil {
		t.Error(err)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for _, testCase := range cases {
		receivedMsg, _, err := reader.ReadLine()
		if err != nil {
			t.Error(err)
		}

		if want, got := testCase.outMsg, string(receivedMsg); want != got {
			t.Errorf("want %q, got %q", want, got)
		}
	}
}
