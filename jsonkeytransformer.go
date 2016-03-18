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
// not allow periods in key names, which ECMA-404 does)
type JSONKeyTransformer struct {
	oldVal   string
	newVal   string
	replacer *strings.Replacer
}

// NewJSONKeyTransformer begins construction of a JSONKeyTransformer.
func NewJSONKeyTransformer() *JSONKeyTransformer {
	return &JSONKeyTransformer{}
}

// OldString sets the string that will be replaced in JSON keys
func (t *JSONKeyTransformer) OldString(oldstring string) *JSONKeyTransformer {
	t.oldVal = oldstring
	return t
}

// NewString sets the string that OldString will be converted to
func (t *JSONKeyTransformer) NewString(newstring string) *JSONKeyTransformer {
	t.newVal = newstring
	return t
}

// Do finishes construction of the JSONKeyTransformer and returns an
// error if any arguments are missing
func (t *JSONKeyTransformer) Do() (*JSONKeyTransformer, error) {
	if t.oldVal == "" || t.newVal == "" {
		return t, fmt.Errorf("bad arguments")
	}
	t.replacer = strings.NewReplacer(t.oldVal, t.newVal)
	return t, nil
}

// recurseTransformMap is a helper method to visit multi-level JSON used by Transform
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

// recurseTransformArr is a helper method to visit multi-level JSON used by recurseTransformMap
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

// Transform accepts a SyslogMsg, and if it is a CEE syslog message, "fixes"
// the JSON keys to be compatible with Elasticsearch 2.x
func (t *JSONKeyTransformer) Transform(msg SyslogMsg) (SyslogMsg, error) {
	if !msg.IsCee {
		return msg, ErrTransform
	}

	transformedStructured := make(map[string]interface{})
	t.recurseTransformMap(msg.JSONValues, transformedStructured)
	newContent, _ := json.Marshal(transformedStructured)
	msg.Content = string(newContent)
	msg.JSONValues = transformedStructured
	return msg, nil
}
