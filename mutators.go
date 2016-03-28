package captainslog

import "fmt"

// TagArrayMutator is a Mutator that modifies a syslog message
// by adding a tag value to an array of tags.
type TagArrayMutator struct {
	tagKey   string
	tagValue string
}

// NewTagArrayMutator constructs a new TagArrayMutator from a supplied
// key and value. The tag will be added to the array at the key if
// it exists - if the key does not exist in the SyslogMsg's JSONValues
// map, it will be  created.
func NewTagArrayMutator(tagKey, tagValue string) Mutator {
	return &TagArrayMutator{
		tagKey:   tagKey,
		tagValue: tagValue,
	}
}

// Mutate modifies the SyslogMsg passed by reference.
func (t *TagArrayMutator) Mutate(msg *SyslogMsg) error {
	var err error

	if _, ok := msg.JSONValues[t.tagKey]; !ok {
		msg.JSONValues[t.tagKey] = make([]interface{}, 0)
	}

	switch val := msg.JSONValues[t.tagKey].(type) {
	case []interface{}:
		msg.JSONValues[t.tagKey] = append(val, t.tagValue)
		return err
	default:
		err = fmt.Errorf("tags key in message was not an array")
		return err
	}
}
