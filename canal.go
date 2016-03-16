package captainslog

type Canal struct {
	input    *InputChanneler
	pipeline []Transformer
	output   *OutputChanneler
}

func NewCanal(input *InputChanneler, output *OutputChanneler, transformers ...Transformer) *Canal {
	c := &Canal{
		input:    input,
		pipeline: transformers,
		output:   output,
	}
	return c
}

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
