# captainslog [![Build Status](https://travis-ci.org/digitalocean/captainslog.svg?branch=master)](https://travis-ci.org/digitalocean/captainslog) [![Doc Status](https://godoc.org/github.com/digitalocean/captainslog?status.png)](https://godoc.org/github.com/digitalocean/captainslog)

## Introduction
Captainslog started as an RFC3164 (Syslog) parser written to solve a specific problem: there was a [breaking change](https://www.elastic.co/guide/en/elasticsearch/reference/current/breaking_20_mapping_changes.html#_field_names_may_not_contain_dots) in the 2.0 release of Elasticsearch that no longer allowed periods in field names. In order to support the creation of a log sanitization service for replacing characters in JSON keys within syslog messages, we created a syslog parser along with an interface for writing "transformers" - plugins for modifying syslog messages in a stream.

We are now continuing development along this path, and plan on working on a set of Inputter, Transformer, and Outputter interface implemenatations that can be used together for processing syslog messages.

The RFC3164 syslog parser has been tested rigorously with go-fuzz and is being used heavily in production.
## Usage
### NewSyslogMsgFromBytes
```go
b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: engage\n")
msg, err := captainslog.NewSyslogMsgFromBytes(b)
if err != nil {
	panic(err)
}
```
### Mutators
```go
tagger := captainslog.NewTagArrayMutator("tags", "captain")
err = tagger.Mutate(&msg)
```

### Transformers
```go
transformer, err := captainslog.NewTagRangeTransformer().
	Select(captainslog.Program, "test:").
	StartMatch(captainslog.Contains, "the final frontier").
	EndMatch(captainslog.Contains, "to boldly go").
	AddTag("tags", "intro").
	WaitDuration(1).
	Do()

msg, err = transformer.Transform(msg)
```

### Pipelines
```go
elasticKeyFixer, err := captainslog.NewJSONKeyTransformer().
	OldString(".").
	NewString("_").
	Do()

err = captainslog.NewPipeline().
	From(reader).
	Transform(elasticKeyFixer).
	To(writer).
	Do()
```

## Contibution Guidelines
We use the [Collective Code Construction Contract](http://rfc.zeromq.org/spec:22) for the development of captainslog. For details, see [CONTRIBUTING.md](https://github.com/digitalocean/captainslog/blob/master/CONTRIBUTING.md).

## License
Copyright 2016 DigitalOcean

Captainslog is released under the [Mozilla Public License, version 2.0](https://github.com/digitalocean/captainslog/blob/master/LICENSE)
