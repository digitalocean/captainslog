package captainslog

import (
	"fmt"
	"testing"
)

func TestJSONForElasticMutatorMutate(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\"}\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
	if err != nil {
		t.Error(err)
	}

	mutator := &JSONForElasticMutator{}
	err = mutator.Mutate(&msg)
	if err != nil {
		t.Error(err)
	}

	if want, got := "@cee:{\"first_name\":\"captain\"}\n", msg.Content; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func ExampleJSONForElasticMutatorMutate() {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\"}\n")

	var msg SyslogMsg
	err := Unmarshal(b, &msg)
	if err != nil {
		panic(err)
	}

	mutator := &JSONForElasticMutator{}
	err = mutator.Mutate(&msg)
	if err != nil {
		panic(err)
	}

	fmt.Printf(msg.Content)
	// Output: @cee:{"first_name":"captain"}
}

func BenchmarkJSONForElasticMutatorMutate(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\"}\n")
	mutator := &JSONForElasticMutator{}

	for i := 0; i < b.N; i++ {
		var msg SyslogMsg
		err := Unmarshal(m, &msg)
		if err != nil {
			panic(err)
		}

		err = mutator.Mutate(&msg)
		if err != nil {
			panic(err)
		}
	}
}
