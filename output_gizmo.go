package captainslog

import (
	"log"
	"time"
)

// Command represents a command that can
// be sent to a Gizmo
type Command int

const (
	// CmdStop tells a Gizmo to stop
	CmdStop Command = iota
)

// OutputGizmo is an output that forwards
// syslog messages to another destination
type OutputGizmo struct {
	CmdChan chan<- Command
	OutChan chan<- *SyslogMsg
	gadget  OutputGadget
}

// NewOutputGizmo accepts an OutputGadget and returns
// an new OutputGizmo that uses that OutputGadget
func NewOutputGizmo(g OutputGadget) *OutputGizmo {
	cmdChan := make(chan Command)
	outChan := make(chan *SyslogMsg)

	o := &OutputGizmo{
		CmdChan: cmdChan,
		OutChan: outChan,
		gadget:  g,
	}

	go o.actor(cmdChan, outChan)

	return o
}

func (o *OutputGizmo) actor(cmdChan <-chan Command, outChan <-chan *SyslogMsg) {
Connect:
	err := o.gadget.Connect()
	if err != nil {
		log.Print("could not connect")
		time.Sleep(time.Duration(o.gadget.RetryInterval()) * time.Second)
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
			_, err := o.gadget.Output(msg)
			if err != nil {
				log.Print("could not send message")
				goto Connect
			}
		}
	}
Stop:
	o.gadget.Close()
	close(o.CmdChan)
	close(o.OutChan)
}
