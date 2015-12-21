package captainslog

import (
	"log"
	"time"
)

// ChannelerCmd represents a command that can
// be sent to a Channeler
type ChannelerCmd int

const (
	// CmdStop tells a Gizmo to stop
	CmdStop ChannelerCmd = iota
)

// OutputChanneler is an outgoing endpoint in a
// chain of Channelers. An OutputChanneler uses
// an OutputAdapter to translate *SyslogMsg
// received on its OutChan to other encodings
// to be sent on other transports.
type OutputChanneler struct {
	CmdChan chan<- ChannelerCmd
	OutChan chan<- *SyslogMsg
	adapter OutputAdapter
}

// NewOutputChanneler accepts an OutputAdapter and returns
// an new OutputChanneler using the adapter.
func NewOutputChanneler(a OutputAdapter) *OutputChanneler {
	cmdChan := make(chan ChannelerCmd)
	outChan := make(chan *SyslogMsg)

	o := &OutputChanneler{
		CmdChan: cmdChan,
		OutChan: outChan,
		adapter: a,
	}

	go o.actor(cmdChan, outChan)

	return o
}

func (o *OutputChanneler) actor(cmdChan <-chan ChannelerCmd, outChan <-chan *SyslogMsg) {
Connect:
	err := o.adapter.Connect()
	if err != nil {
		log.Print("could not connect")
		time.Sleep(time.Duration(o.adapter.RetryInterval()) * time.Second)
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
			_, err := o.adapter.Output(msg)
			if err != nil {
				log.Print("could not send message")
				goto Connect
			}
		}
	}
Stop:
	o.adapter.Close()
	close(o.CmdChan)
	close(o.OutChan)
}
