package main

import (
	"flag"
	"strings"

	"github.com/digitalocean/captainslog"
)

func main() {
	var recvAddress = flag.String("in", ":33333", "TCP endpoint to receive from")
	var sendAddress = flag.String("out", "127.0.0.1:514", "TCP endpoint to send to")

	flag.Parse()

	retryInterval := 5
	outputter := captainslog.NewTCPOutputter(*sendAddress, retryInterval)
	outputChanneler := captainslog.NewOutputChanneler(outputter)

	inputter, _ := captainslog.NewTCPInputter(*recvAddress)
	inputChanneler := captainslog.NewInputChanneler(inputter)

	replacer := strings.NewReplacer(".", "_")
	sanitizer := captainslog.NewJSONKeyTransformer(replacer)

	tagger := captainslog.NewTagRangeTransformer(
		captainslog.NewTagMatcher("kernel:"),
		captainslog.NewContentContainsMatcher("[ cut here ]"),
		captainslog.NewContentContainsMatcher("[ end trace"),
		captainslog.NewTagArrayMutator("tags", "trace"),
		60, 30)

	canal := captainslog.NewCanal(inputChanneler, outputChanneler, tagger, sanitizer)
	canal.Ship()
}
