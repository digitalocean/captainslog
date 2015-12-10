package captainslog

import (
	"encoding/json"
	"strings"
)

type JSONForElasticMutator struct{}

func (m *JSONForElasticMutator) Mutate(msg SyslogMsg) (SyslogMsg, error) {
	if !msg.IsCee {
		return msg, ErrMutate
	}

	var contentStructured map[string]interface{}
	var err error

	err = json.Unmarshal([]byte(msg.Content), &contentStructured)
	if err != nil {
		return msg, err
	}

	mutatedStructured := make(map[string]interface{})
	for k, v := range contentStructured {
		k = strings.Replace(k, ".", "_", -1)
		mutatedStructured[k] = v
	}

	newContent, err := json.Marshal(mutatedStructured)
	if err != nil {
		return msg, err
	}

	mutated := msg
	mutated.Content = string(newContent)
	return mutated, err
}
