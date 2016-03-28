package captainslog_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/digitalocean/captainslog"
)

func ExampleNewSyslogMsgFromBytes() {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: engage\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Syslog message was from host '%s'", msg.Host)
	// Output: Syslog message was from host 'host.example.org'

}

func ExampleTransformer() {
	b := []byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com test: space, the final frontier\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		panic(err)
	}

	transformer, err := captainslog.NewTagRangeTransformer().
		Select(captainslog.Program, "test:").
		StartMatch(captainslog.Contains, "the final frontier").
		EndMatch(captainslog.Contains, "to boldly go").
		AddTag("tags", "intro").
		WaitDuration(1).
		Do()

	msg, err = transformer.Transform(msg)
	if err != nil {
		panic(err)
	}
	fmt.Print(msg.String())
	// Output: <4>2016-03-08T14:59:36.293816+00:00 host.example.com test: @cee:{"msg":"space, the final frontier","tags":["intro"]}
}

func ExampleMutator() {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first_name\":\"kathryn\", \"last_name\":\"janeway\"}\n")
	msg, err := captainslog.NewSyslogMsgFromBytes(b)
	if err != nil {
		panic(err)
	}

	tagger := captainslog.NewTagArrayMutator("tags", "captain")
	err = tagger.Mutate(&msg)
	if err != nil {
		panic(err)
	}

	fmt.Print(msg.String())
	// Output: <191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{"first_name":"kathryn","last_name":"janeway","tags":["captain"]}
}

func ExamplePipeline() {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\",\"last.name\":\"kirk\"}\n")

	reader := bytes.NewBuffer(b)
	writer := bytes.NewBuffer(make([]byte, 0))

	elasticKeyFixer, err := captainslog.NewJSONKeyTransformer().
		OldString(".").
		NewString("_").
		Do()

	err = captainslog.NewPipeline().
		From(reader).
		Transform(elasticKeyFixer).
		To(writer).
		Do()

	if err != io.EOF {
		panic(err)
	}

	fmt.Printf("%s", writer.String())
	// Output: <191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{"first_name":"captain","last_name":"kirk"}
}
