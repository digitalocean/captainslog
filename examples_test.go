package captainslog_test

import (
	"fmt"

	"github.com/digitalocean/captainslog"
)

func ExampleNewSyslogMsgFromBytes() {
	msg, err := captainslog.NewSyslogMsgFromBytes([]byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: engage\n"))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Syslog message was from host '%s'", msg.Host)
	// Output: Syslog message was from host 'host.example.org'

}

func ExampleJSONKeyTransformer() {
	msg, err := captainslog.NewSyslogMsgFromBytes([]byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"benjamin\", \"last.name\":\"sisko\"}\n"))

	transformer := captainslog.NewJSONKeyTransformer(".", "_")

	msg, err = transformer.Transform(msg)
	if err != nil {
		panic(err)
	}

	fmt.Print(msg.String())
	// Output: <191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{"first_name":"benjamin","last_name":"sisko"}
}

func ExampleMutateRangeTransformer() {
	msg, err := captainslog.NewSyslogMsgFromBytes([]byte("<4>2016-03-08T14:59:36.293816+00:00 host.example.com trek: space, the final frontier\n"))
	if err != nil {
		panic(err)
	}

	transformer := captainslog.NewMutateRangeTransformer(
		captainslog.NewTagMatcher("trek:"),
		captainslog.NewContentContainsMatcher("space"),
		captainslog.NewContentContainsMatcher("space"),
		captainslog.NewTagArrayMutator("tags", "intro"),
		1,
	)

	msg, err = transformer.Transform(msg)
	if err != nil {
		panic(err)
	}

	fmt.Print(msg.String())
	// Output: <4>2016-03-08T14:59:36.293816+00:00 host.example.com trek: @cee:{"msg":"space, the final frontier","tags":["intro"]}
}

func ExampleTagArrayMutator() {
	msg, err := captainslog.NewSyslogMsgFromBytes([]byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first_name\":\"kathryn\", \"last_name\":\"janeway\"}\n"))
	if err != nil {
		panic(err)
	}

	captainTagger := captainslog.NewTagArrayMutator("tags", "captain")
	err = captainTagger.Mutate(&msg)
	if err != nil {
		panic(err)
	}

	voyagerTagger := captainslog.NewTagArrayMutator("tags", "voyager")
	err = voyagerTagger.Mutate(&msg)
	if err != nil {
		panic(err)
	}

	fmt.Print(msg.String())
	// Output: <191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{"first_name":"kathryn","last_name":"janeway","tags":["captain","voyager"]}
}
