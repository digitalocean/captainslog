package captainslog

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JSONForElasticMutator struct{}

func (m *JSONForElasticMutator) Mutate(msg *SyslogMsg) error {
	if !msg.Cee {
		return ErrMutate
	}

	var contentStructured map[string]interface{}

	var err error

	err = json.Unmarshal([]byte(msg.Content), &contentStructured)
	if err != nil {
		return err
	}

	mutatedStructured := make(map[string]interface{})
	for k, v := range contentStructured {
		k = strings.Replace(k, ".", "_", -1)
		mutatedStructured[k] = v
	}

	newContent, err := json.Marshal(mutatedStructured)
	if err != nil {
		return err
	}

	msg.Content = fmt.Sprintf("%s%s\n", "@cee:", string(newContent))
	return err
}
