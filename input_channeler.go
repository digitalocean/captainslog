package captainslog

// InputChanneler is an incoming endpoint in
// a chain of Channelers. An InputChanneler uses
// an Inputter to translate RFC3164 syslog bytes
// received over a transport to a *SyslogMsg
// that is sent over a channel.
type InputChanneler struct {
	CmdChan  chan<- ChannelerCmd
	InChan   <-chan *SyslogMsg
	inputter Inputter
}

// NewInputChanneler accepts an Inputter and returns
// a new InputChanneler using the Inputter.
func NewInputChanneler(in Inputter) *InputChanneler {
	cmdChan := make(chan ChannelerCmd)
	inChan := make(chan *SyslogMsg)

	i := &InputChanneler{
		CmdChan:  cmdChan,
		InChan:   inChan,
		inputter: in,
	}

	go i.actor(cmdChan, inChan)
	return i
}

func (i *InputChanneler) actor(cmdChan <-chan ChannelerCmd, inChan chan<- *SyslogMsg) {
	pipe := i.inputter.Listen()
	for {
		select {
		case cmd := <-cmdChan:
			switch cmd {
			case CmdStop:
				goto Stop
			}
		case msg := <-pipe:
			inChan <- msg
		}
	}
Stop:
	close(inChan)
	close(i.CmdChan)
	i.inputter.Close()
}
