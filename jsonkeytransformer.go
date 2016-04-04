package captainslog

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONKeyTransformer is a Transformer implementation that finds periods in JSON
// keys in CEE syslog messages and replaces them. This can be used in
// conjunction with systems such as Elasticsearch 2.x which do not
// fully support ECMA-404 (for instance, Elasticsearch 2.x does
// not allow periods in key names, which ECMA-404 does).
type JSONKeyTransformer struct {
	replacer *strings.Replacer
}

// NewJSONKeyTransformer creates a new JSONKeyTransformer.
func NewJSONKeyTransformer(oldVal, newVal string) *JSONKeyTransformer {
	return &JSONKeyTransformer{
		replacer: strings.NewReplacer(oldVal, newVal),
	}
}

// recurseTransformMap is a helper method to visit multi-level JSON used by Transform.
func (t *JSONKeyTransformer) recurseTransformMap(in, out map[string]interface{}) {
	for k, v := range in {
		transformedKey := t.replacer.Replace(k)
		switch cv := v.(type) {
		case map[string]interface{}:
			nv := make(map[string]interface{})
			out[transformedKey] = nv
			t.recurseTransformMap(cv, nv)
		case []interface{}:
			nv := make([]interface{}, len(cv))
			out[transformedKey] = nv
			t.recurseTransformArr(cv, nv)
		default:
			out[transformedKey] = v
		}
	}
}

// recurseTransformArr is a helper method to visit multi-level JSON used by recurseTransformMap.
func (t *JSONKeyTransformer) recurseTransformArr(in, out []interface{}) {
	for i, v := range in {
		switch cv := v.(type) {
		case map[string]interface{}:
			nv := make(map[string]interface{})
			out[i] = nv
			t.recurseTransformMap(cv, nv)
		case []interface{}:
			nv := make([]interface{}, len(cv))
			out[i] = nv
			t.recurseTransformArr(cv, nv)
		default:
			out[i] = v
		}
	}
}

// Transform accepts a SyslogMsg, and if it is a CEE syslog message, replaces
// the old string with the new string.
func (t *JSONKeyTransformer) Transform(msg SyslogMsg) (SyslogMsg, error) {
	if !msg.IsCee {
		return msg, fmt.Errorf("Transform expected msg.IsCee == true")
	}

	transformedStructured := make(map[string]interface{})
	t.recurseTransformMap(msg.JSONValues, transformedStructured)
	newContent, _ := json.Marshal(transformedStructured)
	msg.Content = string(newContent)
	msg.JSONValues = transformedStructured
	return msg, nil
}
