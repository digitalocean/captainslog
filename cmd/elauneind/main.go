package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/digitalocean/captainslog"
)

const (
	retryInterval = 5
)

func main() {
	var recvAddress = flag.String("in", ":33333", "TCP endpoint to receive from")
	var sendAddress = flag.String("out", "127.0.0.1:514", "TCP endpoint to send to")

	flag.Parse()

	sendChan := make(chan *captainslog.SyslogMsg)

	go func() {
	Connect:
		conn, err := net.Dial("tcp", *sendAddress)
		if err != nil {
			log.Printf("E: connect error - %s", err)
			time.Sleep(time.Duration(retryInterval) * time.Second)
			goto Connect
		}
		defer conn.Close()

		for msg := range sendChan {
			_, err := conn.Write(msg.Bytes())
			if err != nil {
				log.Printf("E: write error - %s", err)
				time.Sleep(time.Duration(retryInterval) * time.Second)
				goto Connect
			}
		}
	}()

	l, err := net.Listen("tcp", *recvAddress)
	if err != nil {
		log.Printf("E: listen error - %s", err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("E: accept error - %s", err)
			continue
		}
		defer conn.Close()

		go func() {
			reader := bufio.NewReader(conn)
			mutator := &captainslog.JSONForElasticMutator{}

			for {
				b, err := reader.ReadBytes('\n')
				if err != nil {
					if err != io.EOF {
						log.Printf("E: readbytes error - %s", err)
					}
					break
				}

				var original captainslog.SyslogMsg
				err = captainslog.Unmarshal(b, &original)
				if err != nil {
					log.Printf("E: unmarshal error - %s", err)
					log.Print(err)
					continue
				}

				mutated, err := mutator.Mutate(original)
				if err != nil {
					log.Printf("E: mutate error - %s", err)
					continue
				}

				sendChan <- &mutated
			}
		}()
	}
}
