package captainslog

import (
	"encoding/json"
	"strings"
)

// JSONForElasticMutator is a Mutator implementation that finds periods in JSON
// keys in CEE syslog messages and replaces them with a character Elasticsearch
// 2.x will accept. Elasticsearch 2.x does not allow the full JSON specification
// for JSON keys.
type JSONForElasticMutator struct{}

// Mutate accepts a SyslogMsg, and if it is a CEE syslog message, "fixes"
// the JSON keys to be compatible with Elasticsearch 2.x
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

	newContent, _ := json.Marshal(mutatedStructured)
	mutated := msg
	mutated.Content = string(newContent)
	return mutated, err
}
