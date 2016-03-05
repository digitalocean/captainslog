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

// recurseMutateMap is a helper method to visit multi-level JSON used by Mutate
func (m *JSONKeyMutator) recurseMutateMap(in, out map[string]interface{}) {
	for k, v := range in {
		mutatedKey := m.replacer.Replace(k)
		switch cv := v.(type) {
		case map[string]interface{}:
			nv := make(map[string]interface{})
			out[mutatedKey] = nv
			m.recurseMutateMap(cv, nv)
		case []interface{}:
			nv := make([]interface{}, len(cv))
			out[mutatedKey] = nv
			m.recurseMutateArr(cv, nv)
		default:
			out[mutatedKey] = v
		}
	}
}

// recurseMutateArr is a helper method to visit multi-level JSON used by recurseMutateMap
func (m *JSONKeyMutator) recurseMutateArr(in, out []interface{}) {
	for i, v := range in {
		switch cv := v.(type) {
		case map[string]interface{}:
			nv := make(map[string]interface{})
			out[i] = nv
			m.recurseMutateMap(cv, nv)
		case []interface{}:
			nv := make([]interface{}, len(cv))
			out[i] = nv
			m.recurseMutateArr(cv, nv)
		default:
			out[i] = v
		}
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
	m.recurseMutateMap(contentStructured, mutatedStructured)

	newContent, _ := json.Marshal(mutatedStructured)
	msg.Content = string(newContent)
	return msg, nil
}
