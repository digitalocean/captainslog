package captainslog

import (
	"encoding/json"
	"strings"
)

// JSONKeyMutator is a Mutator implementation that finds periods in JSON
// keys in CEE syslog messages and replaces them. This can be used in
// conjunction with systems such as Elasticsearch 2.x which do not
// fully support ECMA-404 (for instance, Elasticsearch 2.x does
// not allow periods in key names, which ECMA-404 does)
type JSONKeyMutator struct {
	replacer *strings.Replacer
}

// NewJSONKeyMutator applies a strings.Replacer to all
// keys in a JSON document in a CEE syslog message.
func NewJSONKeyMutator(replacer *strings.Replacer) *JSONKeyMutator {
	return &JSONKeyMutator{
		replacer: replacer,
	}
}

// Mutate accepts a SyslogMsg, and if it is a CEE syslog message, "fixes"
// the JSON keys to be compatible with Elasticsearch 2.x
func (m *JSONKeyMutator) Mutate(msg SyslogMsg) (SyslogMsg, error) {
	if !msg.IsCee {
		return msg, ErrMutate
	}

	var contentStructured map[string]interface{}

	err := json.Unmarshal([]byte(msg.Content), &contentStructured)
	if err != nil {
		return msg, err
	}

	mutatedStructured := make(map[string]interface{})
	for k, v := range contentStructured {
		k = m.replacer.Replace(k)
		mutatedStructured[k] = v
	}

	newContent, _ := json.Marshal(mutatedStructured)
	msg.Content = string(newContent)
	return msg, nil
}
