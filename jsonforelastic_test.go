package captainslog

import (
	"errors"
	"fmt"
	"testing"
)

func checkMutateInterface(m Mutator) {
	return
}

func TestJSONForElastcMutatorIsMutator(t *testing.T) {
	mutator := &JSONForElasticMutator{}
	checkMutateInterface(mutator)
}

func TestJSONForElasticMutatorMutate(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\",\"one.two.three\":\"four.five.six\"}\n")

	var original SyslogMsg
	err := Unmarshal(b, &original)
	if err != nil {
		t.Error(err)
	}

	mutator := &JSONForElasticMutator{}
	mutated, err := mutator.Mutate(original)
	if err != nil {
		t.Error(err)
	}

	if want, got := "{\"first_name\":\"captain\",\"one_two_three\":\"four.five.six\"}", mutated.Content; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestJSONForElasticMutatorMutateNotCee(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: not a json message\n")

	var original SyslogMsg
	err := Unmarshal(b, &original)
	if err != nil {
		t.Error(err)
	}

	mutator := &JSONForElasticMutator{}
	_, err = mutator.Mutate(original)

	if want, got := ErrMutate, err; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestJSONForElasticMutatorMutateBadJSON(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"capt\n")

	var original SyslogMsg
	err := Unmarshal(b, &original)
	if err != nil {
		t.Error(err)
	}

	mutator := &JSONForElasticMutator{}
	_, err = mutator.Mutate(original)

	wantedErr := errors.New("unexpected end of JSON input").Error()

	if want, got := wantedErr, err.Error(); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func ExampleJSONForElasticMutatorMutate() {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\"}\n")

	var original SyslogMsg
	err := Unmarshal(b, &original)
	if err != nil {
		panic(err)
	}

	mutator := &JSONForElasticMutator{}
	mutated, err := mutator.Mutate(original)
	if err != nil {
		panic(err)
	}

	fmt.Printf(mutated.Content)
	// Output: {"first_name":"captain"}
}

func BenchmarkJSONForElasticMutatorMutate(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\"}\n")
	mutator := &JSONForElasticMutator{}

	for i := 0; i < b.N; i++ {
		var original SyslogMsg
		err := Unmarshal(m, &original)
		if err != nil {
			panic(err)
		}

		_, err = mutator.Mutate(original)
		if err != nil {
			panic(err)
		}
	}
}
