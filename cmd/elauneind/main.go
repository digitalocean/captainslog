package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/digitalocean/captainslog"
)

const (
	retryInterval = 5
)

func main() {
	var recvAddress = flag.String("in", ":33333", "TCP endpoint to receive from")
	var sendAddress = flag.String("out", "127.0.0.1:514", "TCP endpoint to send to")
	var logFacility = flag.String("facility", "LOCAL6", "Syslog facility to log to")

	flag.Parse()

	facility, err := captainslog.FacilityTextToFacility(*logFacility)
	if err != nil {
		log.Printf("facility '%s' is not a valid facility\n", *logFacility)
		os.Exit(1)
	}

	ceelog, err := captainslog.NewMostlyFeaturelessLogger(facility)
	if err != nil {
		log.Printf("could not start system logger: %s\n", err)
		os.Exit(1)
	}

	ceelog.InfoWithFields(captainslog.Fields{
		"component": "elauneind",
		"action":    "start",
		"msg":       "starting service",
	})

	a := captainslog.NewTCPOutputAdapter(*sendAddress, retryInterval)
	o := captainslog.NewOutputChanneler(a)

	l, err := net.Listen("tcp", *recvAddress)
	if err != nil {
		ceelog.ErrorWithFields(captainslog.Fields{
			"component": "listener",
			"action":    "listen",
			"msg":       err})

		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			ceelog.ErrorWithFields(captainslog.Fields{
				"ocmponent": "listener",
				"action":    "accept",
				"msg":       err})
			continue
		}
		defer conn.Close()

		go func() {
			reader := bufio.NewReader(conn)
			replacer := strings.NewReplacer(".", "_")
			mutator := captainslog.NewJSONKeyMutator(replacer)

			for {
				b, err := reader.ReadBytes('\n')
				if err != nil {
					if err != io.EOF {
						ceelog.ErrorWithFields(captainslog.Fields{
							"component": "reader",
							"action":    "readbytes",
							"msg":       err})
					}
					break
				}

				var original captainslog.SyslogMsg
				err = captainslog.Unmarshal(b, &original)
				if err != nil {
					ceelog.ErrorWithFields(captainslog.Fields{
						"component": "captainslog",
						"action":    "unmarshal",
						"msg":       err})

					continue
				}

				mutated, err := mutator.Mutate(original)
				if err != nil {
					ceelog.ErrorWithFields(captainslog.Fields{
						"component": "mutator",
						"action":    "Mutate",
						"msg":       err})

					continue
				}

				o.OutChan <- &mutated
			}
		}()
	}
}
