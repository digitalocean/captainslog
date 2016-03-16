package captainslog

import (
	"encoding/json"
	"strings"
)

// JSONKeyTransformer is a Transformer implementation that finds periods in JSON
// keys in CEE syslog messages and replaces them. This can be used in
// conjunction with systems such as Elasticsearch 2.x which do not
// fully support ECMA-404 (for instance, Elasticsearch 2.x does
// not allow periods in key names, which ECMA-404 does)
type JSONKeyTransformer struct {
	replacer *strings.Replacer
}

// NewJSONKeyTransformer applies a strings.Replacer to all
// keys in a JSON document in a CEE syslog message.
func NewJSONKeyTransformer(replacer *strings.Replacer) *JSONKeyTransformer {
	return &JSONKeyTransformer{
		replacer: replacer,
	}
}

// recurseTransformMap is a helper method to visit multi-level JSON used by Transform
func (m *JSONKeyTransformer) recurseTransformMap(in, out map[string]interface{}) {
	for k, v := range in {
		transformedKey := m.replacer.Replace(k)
		switch cv := v.(type) {
		case map[string]interface{}:
			nv := make(map[string]interface{})
			out[transformedKey] = nv
			m.recurseTransformMap(cv, nv)
		case []interface{}:
			nv := make([]interface{}, len(cv))
			out[transformedKey] = nv
			m.recurseTransformArr(cv, nv)
		default:
			out[transformedKey] = v
		}
	}
}

// recurseTransformArr is a helper method to visit multi-level JSON used by recurseTransformMap
func (m *JSONKeyTransformer) recurseTransformArr(in, out []interface{}) {
	for i, v := range in {
		switch cv := v.(type) {
		case map[string]interface{}:
			nv := make(map[string]interface{})
			out[i] = nv
			m.recurseTransformMap(cv, nv)
		case []interface{}:
			nv := make([]interface{}, len(cv))
			out[i] = nv
			m.recurseTransformArr(cv, nv)
		default:
			out[i] = v
		}
	}
}

// Transform accepts a SyslogMsg, and if it is a CEE syslog message, "fixes"
// the JSON keys to be compatible with Elasticsearch 2.x
func (m *JSONKeyTransformer) Transform(msg SyslogMsg) (SyslogMsg, error) {
	if !msg.IsCee {
		return msg, ErrTransform
	}

	transformedStructured := make(map[string]interface{})
	m.recurseTransformMap(msg.JSONValues, transformedStructured)
	newContent, _ := json.Marshal(transformedStructured)
	msg.Content = string(newContent)
	msg.JSONValues = transformedStructured
	return msg, nil
}
