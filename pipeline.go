package captainslog

import (
	"bufio"
	"fmt"
	"io"
)

// A Pipeline accepts bytes from an io.Reader, parses them into
// syslog messages, runs them through a slice of transformers
// in order, and writes the result to an io.Writer. By default,
// messages that result in a parse error are skipped, and
// transformation errors are ignored. If ExitOnParseErr() is
// set, the pipeline will exit and return an error on any
// parse error. If ExitOnTransformErr()is set, the pipeline
// will exit on any transformation error.
type Pipeline struct {
	reader             *bufio.Reader
	transformers       []Transformer
	writer             io.Writer
	exitOnParseErr     bool
	exitOnTransformErr bool
}

// NewPipeline starts the construction of a new Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{
		transformers: make([]Transformer, 0),
	}
}

// From sets the io.Reader the pipeline will read from.
func (p *Pipeline) From(r io.Reader) *Pipeline {
	p.reader = bufio.NewReader(r)
	return p
}

// To sets the io.Writer that the pipeline will write to.
func (p *Pipeline) To(w io.Writer) *Pipeline {
	p.writer = w
	return p
}

// Transform adds a Transformer to the list of transformers
// that each message will run through.
func (p *Pipeline) Transform(t Transformer) *Pipeline {
	p.transformers = append(p.transformers, t)
	return p
}

// ExitOnParseError tells the Pipeline to exit and
// return an error if any message produces a parse error.
func (p *Pipeline) ExitOnParseError() *Pipeline {
	p.exitOnParseErr = true
	return p
}

// ExitOnTransformError tells the pipeliine to exit
// and return an error if any message transformation
// results in an error.
func (p *Pipeline) ExitOnTransformError() *Pipeline {
	p.exitOnTransformErr = true
	return p
}

// Do starts the Pipeline. This call will block until
// the underlying reader returns an io.EOF.
func (p *Pipeline) Do() error {
	if p.reader == nil {
		return fmt.Errorf("no reader")
	}

	if p.writer == nil {
		return fmt.Errorf("no writer")
	}

	for {
		line, err := p.reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		msg, err := NewSyslogMsgFromBytes(line)
		if err != nil {
			if p.exitOnParseErr {
				return err
			}
			continue
		}

		for _, t := range p.transformers {
			msg, err = t.Transform(msg)
			if err != nil {
				if p.exitOnTransformErr {
					return err
				}
				continue
			}
		}

		_, err = p.writer.Write(msg.Bytes())
		if err != nil {
			return err
		}
	}
}
