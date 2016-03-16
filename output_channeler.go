package captainslog

import (
	"log"
	"time"
)

// OutputChanneler is an outgoing endpoint in a
// chain of Channelers. An OutputChanneler uses
// an Outputter to translate *SyslogMsg
// received on its OutChan to other encodings
// to be sent on other transports.
type OutputChanneler struct {
	CmdChan   chan<- ChannelerCmd
	OutChan   chan<- *SyslogMsg
	outputter Outputter
}

// NewOutputChanneler accepts an Outputter and returns
// an new OutputChanneler using the Outputter.
func NewOutputChanneler(a Outputter) *OutputChanneler {
	cmdChan := make(chan ChannelerCmd)
	outChan := make(chan *SyslogMsg)

	o := &OutputChanneler{
		CmdChan:   cmdChan,
		OutChan:   outChan,
		outputter: a,
	}

	go o.actor(cmdChan, outChan)

	return o
}

func (o *OutputChanneler) actor(cmdChan <-chan ChannelerCmd, outChan <-chan *SyslogMsg) {
Connect:
	err := o.outputter.Connect()
	if err != nil {
		log.Print("could not connect")
		time.Sleep(time.Duration(o.outputter.RetryInterval()) * time.Second)
		goto Connect
	}

	for {
		select {
		case cmd := <-cmdChan:
			switch cmd {
			case CmdStop:
				goto Stop
			}
		case msg := <-outChan:
			_, err := o.outputter.Output(msg)
			if err != nil {
				log.Print("could not send message")
				goto Connect
			}
		}
	}
Stop:
	o.outputter.Close()
	close(o.CmdChan)
	close(o.OutChan)
}
