package captainslog

// Canal ties together an input, a pipeline of transformers, and an output.
type Canal struct {
	input    *InputChanneler
	pipeline []Transformer
	output   *OutputChanneler
}

// NewCanal accepts an InputChanneler, an OutputChanneler,  and a variadic list
// of Transformers and returns a Canal.
func NewCanal(input *InputChanneler, output *OutputChanneler, transformers ...Transformer) *Canal {
	c := &Canal{
		input:    input,
		pipeline: transformers,
		output:   output,
	}
	return c
}

// Ship starts the Canal.
func (c *Canal) Ship() {
	for {
		var err error
		msg := <-c.input.InChan
		for _, transformer := range c.pipeline {
			*msg, err = transformer.Transform(*msg)
			if err != nil {
				continue
			}
		}
		c.output.OutChan <- msg
	}
}
