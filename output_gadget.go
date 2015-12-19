package captainslog

// OutputGadget is an interface for Output Gadgets.
// Output Gadgets plug in to an OutputGizmo to provide
// its functionality.
type OutputGadget interface {
	Output(s *SyslogMsg) (int, error)
	Connect() error
	RetryInterval() int
	Close()
}
