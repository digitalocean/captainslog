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
## Serialize a captainslog.SyslogMsg to RFC3164 bytes:
```go
b := msg.Bytes()
```
## Contibution Guidelines
We use the [Collective Code Construction Contract](http://rfc.zeromq.org/spec:22) for the development of captainslog. For details, see [CONTRIBUTING.md](https://github.com/digitalocean/captainslog/blob/master/CONTRIBUTING.md).

## License
Copyright 2016 DigitalOcean

Captainslog is released under the [Mozilla Public License, version 2.0](https://github.com/digitalocean/captainslog/blob/master/LICENSE)
