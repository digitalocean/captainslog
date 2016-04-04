# captainslog [![Build Status](https://travis-ci.org/digitalocean/captainslog.svg?branch=master)](https://travis-ci.org/digitalocean/captainslog) [![Doc Status](https://godoc.org/github.com/digitalocean/captainslog?status.png)](https://godoc.org/github.com/digitalocean/captainslog)

## Introduction
Captainslog started as an RFC3164 (Syslog) parser written to solve a specific problem: there was a [breaking change](https://www.elastic.co/guide/en/elasticsearch/reference/current/breaking_20_mapping_changes.html#_field_names_may_not_contain_dots) in the 2.0 release of Elasticsearch that no longer allowed periods in field names. In order to support the creation of a log sanitization service for replacing characters in JSON keys within syslog messages, we created a syslog parser along with an interface for writing "transformers" - plugins for modifying syslog messages in a stream.

We are now continuing development along this path, and plan on working on a set of Inputter, Transformer, and Outputter interface implemenatations that can be used together for processing syslog messages.

The RFC3164 syslog parser has been tested rigorously with go-fuzz and is being used heavily in production.
## Usage
### NewSyslogMsgFromBytes
NewSyslogMsgFromBytes accepts a []byte containing an RFC3164 message and returns a SyslogMsg. If the original RFC3164 message is a CEE enhanced message, the JSON will be parsed into the JSONValues map[string]inferface{}.

```go
b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: engage\n")
msg, err := captainslog.NewSyslogMsgFromBytes(b)
if err != nil {
	panic(err)
}
```
### Mutators
#### TagArrayMutator
TagArrayMutator is a Mutator that modifies a syslog message by adding a tag value to an array of tags. If the message was not already a CEE message, it will be converted to one and the original message content will be assigned to the "msg" attribute.
```go
```
### Transformers
#### JSONKeyTransformer
JSONKeyTransformer is a Transformer implementation that finds periods in JSON keys in CEE syslog messages and replaces them. This can be used in conjunction with systems such as Elasticsearch 2.x which do not fully support ECMA-404 (for instance, Elasticsearch 2.x does not allow periods in key names, which ECMA-404 does).
```go
msg, err := captainslog.NewSyslogMsgFromBytes([]byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"benjamin\", \"last.name\":\"sisko\"}\n"))

transformer := captainslog.NewJSONKeyTransformer(".", "_")

msg, err = transformer.Transform(msg)
if err != nil {
	panic(err)
}

fmt.Print(msg.String())
```
```
<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{"first_name":"benjamin","last_name":"sisko"}
```

#### MutateRangeTransformer
MutateRangeTransformer is a Transformer implementation that mutates log lines that meet a selection criteria and are logged between a start and end match. Matches are performed by implementations of the Matcher interface.
```go
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
```
```
<4>2016-03-08T14:59:36.293816+00:00 host.example.com trek: @cee:{"msg":"space, the final frontier","tags":["intro"]}
```

## Contibution Guidelines
We use the [Collective Code Construction Contract](http://rfc.zeromq.org/spec:22) for the development of captainslog. For details, see [CONTRIBUTING.md](https://github.com/digitalocean/captainslog/blob/master/CONTRIBUTING.md).

## License
Copyright 2016 DigitalOcean

Captainslog is released under the [Mozilla Public License, version 2.0](https://github.com/digitalocean/captainslog/blob/master/LICENSE)
