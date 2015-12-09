# captainslog [![Build Status](https://travis-ci.org/digitalocean/captainslog.svg?branch=master)](https://travis-ci.org/digitalocean/captainslog) [![Doc Status](https://godoc.org/github.com/digitalocean/captainslog?status.png)](https://godoc.org/github.com/digitalocean/captainslog)

## Installation

```
go get github.com/digitalocean/captainslog
```

## Usage
```go
package main

import (
	"fmt"

	"github.com/digitalocean/captainslog"
)

func main() {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	var msg captainslog.SyslogMsg
	err := captainslog.Unmarshal(b, &msg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Syslog message was from host '%s'\n", msg.Host)
}
```

