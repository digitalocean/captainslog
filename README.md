# captainslog [![Build Status](https://travis-ci.org/digitalocean/captainslog.svg?branch=master)](https://travis-ci.org/digitalocean/captainslog) [![Doc Status](https://godoc.org/github.com/digitalocean/captainslog?status.png)](https://godoc.org/github.com/digitalocean/captainslog)

Efficient and accurate syslog parser written in Golang. Tested rigorously with go-fuzz and being used to process over a billion syslog messages a day in production.

## Usage
### NewSyslogMsgFromBytes
NewSyslogMsgFromBytes accepts a []byte containing an RFC3164 message and returns a SyslogMsg. If the original RFC3164 message is a CEE message, the JSON will be parsed into the JSONValues map[string]inferface{}.

```go
b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: engage\n")
msg, err := captainslog.NewSyslogMsgFromBytes(b)
if err != nil {
	panic(err)
}
```

## Contibution Guidelines
We use the [Collective Code Construction Contract](http://rfc.zeromq.org/spec:22) for the development of captainslog. For details, see [CONTRIBUTING.md](https://github.com/digitalocean/captainslog/blob/master/CONTRIBUTING.md).

## License
Copyright 2016 DigitalOcean

Captainslog is released under the [Mozilla Public License, version 2.0](https://github.com/digitalocean/captainslog/blob/master/LICENSE)
