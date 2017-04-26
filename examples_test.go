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

	fmt.Printf("Syslog message was from host %q", msg.Host)
	// Output: Syslog message was from host "host.example.org"

}
