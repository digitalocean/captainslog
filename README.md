# captainslog [![Build Status](https://travis-ci.org/digitalocean/captainslog.svg?branch=master)](https://travis-ci.org/digitalocean/captainslog) [![Doc Status](https://godoc.org/github.com/digitalocean/captainslog?status.png)](https://godoc.org/github.com/digitalocean/captainslog)

## Introduction
Captainslog started as an RFC3164 (Syslog) parser written to solve a specific problem: there was a [breaking change](https://www.elastic.co/guide/en/elasticsearch/reference/current/breaking_20_mapping_changes.html#_field_names_may_not_contain_dots) in the 2.0 release of Elasticsearch that no longer allowed periods in field names. In order to support the creation of a log sanitization service for replacing characters in JSON keys within syslog messages, we created a syslog parser along with an interface for writing "mutators" - plugins for modifying syslog messages in a stream.

We are now continuing development along this path, and plan on working on a set of Inputter, Transformer, and Outputter interface implemenatations that can be used together for processing syslog messages.

## Usage
```go
package main

import (
	"fmt"
	"strings"

	"github.com/digitalocean/captainslog"
)

func main() {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"my.message\": \"hello world\"}\n")

	var originalMsg captainslog.SyslogMsg
	err := captainslog.Unmarshal(b, &originalMsg)
	if err != nil {
		panic(err)
	}

	replacer := strings.NewReplacer(".", "_")
	mutator := captainslog.NewJSONKeyTransformer(replacer)

	mutatedMsg, err := mutator.Transform(originalMsg)
	if err != nil {
		panic(err)
	}

	fmt.Print(mutatedMsg.String())
}
```

## Contibution Guidelines
We use the [Collective Code Construction Contract](http://rfc.zeromq.org/spec:22) for the development of captainlog. For details, see [CONTRIBUTING.md](https://github.com/digitalocean/captainslog/blob/master/CONTRIBUTING.md).

## License
Copyright 2016 DigitalOcean

Captainslog is released under the [Mozilla Public License, version 2.0](https://github.com/digitalocean/captainslog/blob/master/LICENSE)
