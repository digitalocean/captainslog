# captainslog [![Build Status](https://travis-ci.org/digitalocean/captainslog.svg?branch=master)](https://travis-ci.org/digitalocean/captainslog) [![Doc Status](https://godoc.org/github.com/digitalocean/captainslog?status.png)](https://godoc.org/github.com/digitalocean/captainslog)

Construct, emit, and parse Syslog messages.
# Usage
## Create a captainslog.SyslogMsg from RF3164 bytes:
```go
b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: engage\n")
msg, err := captainslog.NewSyslogMsgFromBytes(b)
if err != nil {
	panic(err)
}
```
## Create a captainslog.SyslogMsg by setting its fields:
```go
msg := captainslog.NewSyslogMsg()
msg.SetFacility(captainslog.Local7)
msg.SetSeverity(captainslog.Err)

msgTime, err := time.Parse("2006 Jan 02 15:04:05", "2017 Aug 15 16:18:34")
if err != nil {
	t.Error(err)
}

msg.SetTime(msgTime)
msg.SetProgram("myprogram")
msg.SetPid("12")
msg.SetHost("host.example.com")
```
captainslog.NewSyslogMsg accepts the following functional options (note: these may also be passed to SyslogMsg.Bytes() and SyslogMsg.String):

**captainslog.OptionUseLocalFormat** tells SyslogMsg.String() and SyslogMsg.Byte() to format the message to be compatible with writing to /dev/log rather than over the wire.

**captainslog.OptionUseRemoteFormat** tells SyslogMsg.String() and SyslogMsg.Byte() to use wire format for the message instead of local format.
## Serialize a captainslog.SyslogMsg to RFC3164 bytes:
```go
b := msg.Bytes()
```
## Create a captainslog.Parser and parse a message:
```go
p := captainslog.NewParser(<options>)
msg, err := p.ParseBytes([]byte(line)
```
Both captainslog.NewSyslogMsgFromBytes and captainslog.NewParser accept the following functional arguments:

**captainslog.OptionNoHostname** sets the parser to not expect the hostname as part of the syslog message, and instead ask the host for its hostname.

**captainslog.OptionDontParseJSON** sets the parser to not parse JSON in the content field of the message. A subsequent call to SyslogMsg.String() or SyslogMsg.Bytes() will then use SyslogMsg.Content for the content field, unless SyslogMsg.JSONValues have been added since the message was originally parsed. If SyslogMsg.JSONValues have been added, the call to SyslogMsg.String() or SyslogMsg.Bytes() will then parse the JSON, and merge the results with the keys in SyslogMsg.JSONVaues.

**captainslog.OptionLocation** is a helper function to configure the parser to parse time in the given timezone, If the parsed time contains a valid timezone identifier this takes precedence. Default timezone is UTC.
## Contibution Guidelines
We use the [Collective Code Construction Contract](http://rfc.zeromq.org/spec:22) for the development of captainslog. For details, see [CONTRIBUTING.md](https://github.com/digitalocean/captainslog/blob/master/CONTRIBUTING.md).
## License
Copyright 2016 DigitalOcean

Captainslog is released under the [Mozilla Public License, version 2.0](https://github.com/digitalocean/captainslog/blob/master/LICENSE)
