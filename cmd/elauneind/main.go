package main

import (
	"flag"

	"github.com/digitalocean/captainslog"
)

func main() {
	var recvAddress = flag.String("in", ":33333", "TCP endpoint to receive from")
	var sendAddress = flag.String("out", "127.0.0.1:514", "TCP endpoint to send to")

	flag.Parse()

	retryInterval := 5
	outputter := captainslog.NewTCPOutputter(*sendAddress, retryInterval)
	outputChanneler := captainslog.NewOutputChanneler(outputter)

	inputter, err := captainslog.NewTCPInputter(*recvAddress)
	if err != nil {
		panic(err)
	}

	inputChanneler := captainslog.NewInputChanneler(inputter)

	sanitizer, err := captainslog.NewJSONKeyTransformer().
		OldString(".").
		NewString("_").
		Do()

	if err != nil {
		panic(err)
	}

	tagger, err := captainslog.NewTagRangeTransformer().
		Select(captainslog.Program, "kernel:").
		StartMatch(captainslog.Contains, "[ cut here ]").
		EndMatch(captainslog.Contains, "[ end trace").
		AddTag("tags", "trace").
		WaitDuration(60).
		Do()

	if err != nil {
		panic(err)
	}

	canal := captainslog.NewCanal(inputChanneler, outputChanneler, tagger, sanitizer)
	canal.Ship()
}
