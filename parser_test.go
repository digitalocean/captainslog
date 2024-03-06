package captainslog_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/digitalocean/captainslog"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		name            string
		input           string
		options         []func(*captainslog.Parser)
		sanitizeProgram bool
		err             error
		facility        captainslog.Facility
		severity        captainslog.Severity
		year            int
		month           int
		day             int
		hour            int
		minute          int
		second          int
		millis          int
		offset          int
		host            string
		program         string
		tag             string
		pid             string
		cee             bool
		json            bool
		content         string
		jsonKeys        []string
	}{
		{
			name:     "parse plain text with pid",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[12]: hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test[12]:",
			pid:      "12",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse plain text with period in tag",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test.rb: hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test.rb",
			tag:      "test.rb:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse plain text with bracket in hostname",
			input:    "<36>2006-01-02T15:04:05.999999-07:00 pdu.example.org [Sentry3_53d65d] AUTH: User \"ADMN\" logged out -- connection source \"CONSOLE\" [Console]\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Auth,
			severity: captainslog.Warning,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "pdu.example.org",
			program:  "Sentry3_53d65d",
			tag:      "[Sentry3_53d65d]",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " AUTH: User \"ADMN\" logged out -- connection source \"CONSOLE\" [Console]",
			jsonKeys: []string{},
		},
		{
			name:     "parse plain text with bracket in hostname and pid",
			input:    "<36>2006-01-02T15:04:05.999999-07:00 pdu.example.org [Sentry3_53d65d][88] AUTH: User \"ADMN\" logged out -- connection source \"CONSOLE\" [Console]\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Auth,
			severity: captainslog.Warning,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "pdu.example.org",
			program:  "Sentry3_53d65d",
			tag:      "[Sentry3_53d65d][88]",
			pid:      "88",
			cee:      false,
			json:     false,
			content:  " AUTH: User \"ADMN\" logged out -- connection source \"CONSOLE\" [Console]",
			jsonKeys: []string{},
		},
		{
			name:     "parse plain test without pid",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse plain test without pid",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test-with-hyphen: hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test-with-hyphen",
			tag:      "test-with-hyphen:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},

		{
			name:     "parse plain test without pid no space after tag",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  "hello world",
			jsonKeys: []string{},
		},
		{
			name:     "missing host with OptionNoHostname",
			input:    "<86>Jul 24 11:53:47 sudo: pam_unix(sudo:session): session opened for user root by (uid=0)\n",
			options:  []func(*captainslog.Parser){captainslog.OptionNoHostname},
			err:      nil,
			facility: captainslog.AuthPriv,
			severity: captainslog.Info,
			year:     time.Now().Year(),
			month:    7,
			day:      24,
			hour:     11,
			minute:   53,
			second:   47,
			millis:   0,
			offset:   0,
			host:     "localhost",
			program:  "sudo",
			tag:      "sudo:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " pam_unix(sudo:session): session opened for user root by (uid=0)",
			jsonKeys: []string{},
		},
		{
			name:     "parse @cee json",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      true,
			json:     true,
			content:  "{\"a\":\"b\"}",
			jsonKeys: []string{"a"},
		},
		{
			name:     "parse log with no micros",
			input:    "<191>2006-01-02T15:04:05.999-07:00 host.example.org test[12]: hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999000,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test[12]:",
			pid:      "12",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse log with no millis",
			input:    "<171>2015-12-18T18:08:17+00:00 host.example.org test[12]: hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local5,
			severity: captainslog.Err,
			year:     2015,
			month:    12,
			day:      18,
			hour:     18,
			minute:   8,
			second:   17,
			millis:   0,
			offset:   0,
			host:     "host.example.org",
			program:  "test",
			tag:      "test[12]:",
			pid:      "12",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse cee with space",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee: {\"a\":\"b\"}\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      true,
			json:     true,
			content:  " {\"a\":\"b\"}",
			jsonKeys: []string{"a"},
		},
		{
			name:     "parse cee with no space",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      true,
			json:     true,
			content:  "{\"a\":\"b\"}",
			jsonKeys: []string{"a"},
		},
		{
			name:     "parse json without cee",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: {\"a\":\"b\"}\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      false,
			json:     true,
			content:  " {\"a\":\"b\"}",
			jsonKeys: []string{"a"},
		},
		{
			name:    "parse program with slashes",
			input:   "<191>2006-01-02T15:04:05.999999-07:00 host.example.org /usr/bin/test: some content\n",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadTag,
		},
		{
			name:            "parse program with slashes with OptionSanitizeProgram",
			input:           "<191>2006-01-02T15:04:05.999999-07:00 host.example.org /usr/bin/test: some content\n",
			options:         []func(*captainslog.Parser){captainslog.OptionSanitizeProgram},
			sanitizeProgram: true,
			err:             nil,
			facility:        captainslog.Local7,
			severity:        captainslog.Debug,
			year:            2006,
			month:           1,
			day:             2,
			hour:            15,
			minute:          4,
			second:          5,
			millis:          999999,
			offset:          -25200,
			host:            "host.example.org",
			program:         "test",
			tag:             "test:",
			pid:             "",
			cee:             false,
			json:            false,
			content:         " some content",
			jsonKeys:        []string{},
		},
		{
			name:     "parse json without cee with OptionDontParseJSON",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: {\"a\":\"b\"}\n",
			options:  []func(*captainslog.Parser){captainslog.OptionDontParseJSON},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " {\"a\":\"b\"}",
			jsonKeys: []string{},
		},
		{
			name:     "parse cee early termination",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " @cee",
			jsonKeys: []string{},
		},
		{
			name:     "parse with tag with no colon",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse with truncated subseconds",
			input:    "<191>2006-01-02T15:04:05.99999-07:00 host.example.org test hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999990,
			offset:   -25200,
			host:     "host.example.org",
			program:  "test",
			tag:      "test",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse with unix time",
			input:    fmt.Sprintf("<191>%s host.example.org test: hello world\n", generateDate("Mon Jan _2 15:04:05")),
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     time.Now().Year(),
			month:    1,
			day:      1,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   0,
			offset:   0,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse with ansi time",
			input:    "<191>Mon Jan  2 15:04:05 2006 host.example.org test: hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   0,
			offset:   0,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse with unix date time",
			input:    "<191>Mon Jan  2 15:04:05 MST 2006 host.example.org test: hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   0,
			offset:   0,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:    "parse cee with no json after",
			input:   "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:\n",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadContent,
		},
		{
			name:    "parse bad content",
			input:   "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test:\n",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadContent,
		},
		{
			name:    "parse no priority",
			input:   "2006-01-02T15:04:05.999999-07:00 host.example.org test:\n",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadPriority,
		},
		{
			name:    "parse no priority end",
			input:   "<1912006-01-02T15:04:05.999999-07:00 host.example.org test:\n",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadPriority,
		},
		{
			name:    "parse bad priority end",
			input:   "<9999>2006-01-02T15:04:05.999999-07:00 host.example.org test:\n",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadPriority,
		},
		{
			name:    "parse priority truncated",
			input:   "<99\n",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadPriority,
		},
		{
			name:    "parse date truncated",
			input:   "<191>2006-01-02T15:0",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadTime,
		},
		{
			name:    "parse host truncated",
			input:   "<191>2006-01-02T15:04:05.999999-07:00 ",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadHost,
		},
		{
			name:    "parse no host",
			input:   "<191>2006-01-02T15:04:05.999999-07:00 host.examp",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadHost,
		},
		{
			name:    "parse tag truncated",
			input:   "<191>2006-01-02T15:04:05.999999-07:00 host.example.org tes",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadTag,
		},
		{
			name:    "parse no tag",
			input:   "<191>2006-01-02T15:04:05.999999-07:00 host.example.org ",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadTag,
		},
		{
			name:    "parse rsyslog could not determine hostname but has ip",
			input:   "<28>2017-08-01T13:49:57.537348+00:00 192.168.1.123 <localhost> test[146407]: hello world",
			options: []func(*captainslog.Parser){},
			err:     captainslog.ErrBadTag,
		},
		{
			name:     "parse with time zone option",
			input:    "<191>Mon Jan  2 15:04:05 2006 host.example.org test: hello world\n",
			options:  []func(*captainslog.Parser){captainslog.OptionLocation(time.FixedZone("UTC-8", -8*60*60))},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   0,
			offset:   -8 * 60 * 60,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse with time zone option overridden by specified time zone",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n",
			options:  []func(*captainslog.Parser){captainslog.OptionLocation(time.FixedZone("UTC-8", -8*60*60))},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -7 * 60 * 60,
			host:     "host.example.org",
			program:  "test",
			tag:      "test:",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},
		{
			name:     "parse CEE using optional GJSON parser",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[241]: @cee:{\"egid\":0,\"eid\":0,\"env\":\"production\",\"host\":\"myhost.example.org\",\"level\":\"info\",\"msg\":\"request complete\",\"pid\":7,\"pname\":\"/bin/myprogram\",\"req_method\":\"GET\",\"req_path\":\"/bin/myprogram/\",\"req_remote_ip\":\"172.0.0.1\",\"req_useragent\":\"my-client;go=go1.11.1\",\"resp_bytes_per_sec\":40562776957,\"resp_code\":200,\"resp_duration\":\"1.345193477s\",\"resp_duration_ms\":1345.193477,\"resp_latency\":\"1.345069502s\",\"resp_latency_ms\":1345.069502,\"resp_mebibytes_per_sec\":38683.678586006165,\"resp_size\":54922,\"system\":\"server\",\"time\":\"2019-03-04T19:21:26.895323594Z\",\"version\":\"5b82fdcddaf0286e7fec3a5f8dbf7a67a325fd6b\"}\n",
			options:  []func(*captainslog.Parser){captainslog.OptionLocation(time.FixedZone("UTC-8", -8*60*60)), captainslog.OptionUseGJSONParser},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -7 * 60 * 60,
			host:     "host.example.org",
			program:  "test",
			tag:      "test[241]:",
			pid:      "241",
			cee:      true,
			json:     true,
			content:  "{\"egid\":0,\"eid\":0,\"env\":\"production\",\"host\":\"myhost.example.org\",\"level\":\"info\",\"msg\":\"request complete\",\"pid\":7,\"pname\":\"/bin/myprogram\",\"req_method\":\"GET\",\"req_path\":\"/bin/myprogram/\",\"req_remote_ip\":\"172.0.0.1\",\"req_useragent\":\"my-client;go=go1.11.1\",\"resp_bytes_per_sec\":40562776957,\"resp_code\":200,\"resp_duration\":\"1.345193477s\",\"resp_duration_ms\":1345.193477,\"resp_latency\":\"1.345069502s\",\"resp_latency_ms\":1345.069502,\"resp_mebibytes_per_sec\":38683.678586006165,\"resp_size\":54922,\"system\":\"server\",\"time\":\"2019-03-04T19:21:26.895323594Z\",\"version\":\"5b82fdcddaf0286e7fec3a5f8dbf7a67a325fd6b\"}",
			jsonKeys: []string{"egid", "eid", "env", "host", "level", "msg", "pid", "pname", "req_method", "req_path", "req_remote_ip", "req_useragent", "resp_bytes_per_sec", "resp_code", "resp_duration", "resp_duration_ms", "resp_latency", "resp_latency_ms", "resp_mebibytes_per_sec", "resp_size", "system", "time", "version"},
		},
		{
			name:     "parse with syslog tag over 32 chars",
			input:    "<191>2006-01-02T15:04:05.999999-07:00 host.example.org abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz hello world\n",
			options:  []func(*captainslog.Parser){},
			err:      nil,
			facility: captainslog.Local7,
			severity: captainslog.Debug,
			year:     2006,
			month:    1,
			day:      2,
			hour:     15,
			minute:   4,
			second:   5,
			millis:   999999,
			offset:   -25200,
			host:     "host.example.org",
			program:  "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz",
			tag:  "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz",
			pid:      "",
			cee:      false,
			json:     false,
			content:  " hello world",
			jsonKeys: []string{},
		},


	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := captainslog.NewParser(tc.options...)

			msg, err := p.ParseBytes([]byte(tc.input))

			if want, got := tc.err, err; want != got {
				t.Errorf("error: want %v, got %v", want, got)
			}

			if tc.err == nil {
				if want, got := tc.facility, msg.Pri.Facility; want != got {
					t.Errorf("facility: want %q, got %q", want, got)
				}

				if want, got := tc.severity, msg.Pri.Severity; want != got {
					t.Errorf("facility: want %q, got %q", want, got)
				}

				if want, got := tc.year, msg.Time.Year(); want != got {
					t.Errorf("year: want %d, got %d", want, got)
				}

				if want, got := time.Month(tc.month), msg.Time.Month(); want != got {
					t.Errorf("month: want %d, got %d", want, got)
				}

				if want, got := tc.day, msg.Time.Day(); want != got {
					t.Errorf("day: want %d, got %d", want, got)
				}

				if want, got := tc.hour, msg.Time.Hour(); want != got {
					t.Errorf("hour: want %d, got %d", want, got)
				}

				if want, got := tc.minute, msg.Time.Minute(); want != got {
					t.Errorf("minute: want %d, got %d", want, got)
				}

				if want, got := tc.second, msg.Time.Second(); want != got {
					t.Errorf("second: want %d, got %d", want, got)
				}

				if want, got := tc.millis, msg.Time.Nanosecond()/1000; want != got {
					t.Errorf("millis: want %d, got %d", want, got)
				}

				_, offsetSeconds := msg.Time.Zone()
				if want, got := tc.offset, offsetSeconds; want != got {
					t.Errorf("offset: want %d, got %d", want, got)
				}

				var useLocal bool
				if tc.host == "localhost" {
					host, err := os.Hostname()
					if err != nil {
						t.Error(err)
					}
					tc.host = host
					useLocal = true
				}

				if want, got := tc.host, msg.Host; want != got {
					t.Errorf("host: want %q, got %q", want, got)
				}

				if want, got := tc.program, msg.Tag.Program; want != got {
					t.Errorf("program: want %q, got %q", want, got)
				}

				if want, got := tc.tag, msg.Tag.String(); want != got {
					t.Errorf("tag: want %q, got %q", want, got)
				}

				if want, got := tc.pid, msg.Tag.Pid; want != got {
					t.Errorf("pid: want %q, got %q", want, got)
				}

				if want, got := tc.cee, msg.IsCee; want != got {
					t.Errorf("cee: want %v, got %v", want, got)
				}

				if want, got := tc.json, msg.IsJSON; want != got {
					t.Errorf("json: want %v, got %v for msg %s", want, got, msg.String())
				}

				if want, got := tc.content, msg.Content; want != got {
					t.Errorf("content: want %q, got %q", want, got)
				}

				if want, got := len(tc.jsonKeys), len(msg.JSONValues); want != got {
					t.Errorf("keys: want %d, got %d", want, got)
				}

				for _, v := range tc.jsonKeys {
					if _, ok := msg.JSONValues[v]; !ok {
						t.Errorf("could not find expected key %q in msg.JSONValues", v)
					}
				}

				// NOTE: for now we do not do a byte level reconstruction test if the original
				// message had JSON, since re-encoding the JSON to reconstruct the message
				// can remove spaces that were in the origin message. Also, if we choose to
				// sanitize the program name, it will lose the context of a prefixed path, so
				// comparing with the original bytes would fail.
				if !useLocal && !msg.IsJSON && !tc.sanitizeProgram {
					if want, got := 0, bytes.Compare([]byte(tc.input), msg.Bytes()); want != got {
						t.Errorf("want %q, got  %q", tc.input, msg.Bytes())
					}
				}
			}
		})
	}
}

func testFuzzFindings(fuzzData string, t *testing.T) {
	b := []byte(fuzzData)
	_, err := captainslog.NewSyslogMsgFromBytes(b)

	if want, got := false, err == nil; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestFuzzFindings(t *testing.T) {
	inputs := []string{
		"<0>Mon Jan 00 00:00:000 0 ",
		"<0>Mon Jan 00 00:00:000 :",
		"<0>Jan 02 00:00:00",
	}

	for _, fuzzData := range inputs {
		testFuzzFindings(fuzzData, t)
	}
}

func TestBadPids(t *testing.T) {
	inputs := []string{
		"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[1: hello world\n",
		"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[1 hello world\n",
		"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[\n",
		"<191>2006-01-02T15:04:05.999999-07:00 host.example.org test[]\n",
	}

	for _, fuzzData := range inputs {
		testFuzzFindings(fuzzData, t)
	}
}

// This test demonstrates that long lines (this one is just shy of 32k) parse successfully; max line length is not enforced by captainslog today.
func TestLongOssecLine(t *testing.T) {
	m := []byte("<168>2022-02-15T20:39:21.508396+00:00 myhost.mydc.internal.digitalocean.com ossec-alerts @cee: {\"message\":\" @cee:{\\\"rule\\\":{\\\"level\\\":13,\\\"comment\\\":\\\"Non standard syslog message (size too large).\\\",\\\"sidid\\\":1003,\\\"firedtimes\\\":652081,\\\"groups\\\":[\\\"syslog\\\",\\\"errors\\\"]},\\\"id\\\":\\\"1644957561.960287735\\\",\\\"TimeStamp\\\":1644957561000,\\\"location\\\":\\\"/var/log/syslog\\\",\\\"full_log\\\":\\\"\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"}\n")

	msg, err := captainslog.NewSyslogMsgFromBytes(m)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		assert.NoError(t, err)
	}
	assert.Equal(t, "myhost.mydc.internal.digitalocean.com", msg.Host)
	assert.Equal(t, 15960, len(msg.Content))
}

func BenchmarkParserParse(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Content != " hello world" {
			panic("unexpected msg.Content")
		}
	}
}

func BenchmarkParserParseCEE(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Host != "host.example.org" {
			panic("unexpected msg.Host")
		}
	}
}

func BenchmarkParserParseLongCEE(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"egid\":0,\"eid\":0,\"env\":\"production\",\"host\":\"myhost.example.org\",\"level\":\"info\",\"msg\":\"request complete\",\"pid\":7,\"pname\":\"/bin/myprogram\",\"req_method\":\"GET\",\"req_path\":\"/bin/myprogram/\",\"req_remote_ip\":\"172.0.0.1\",\"req_useragent\":\"my-client;go=go1.11.1\",\"resp_bytes_per_sec\":40562776957,\"resp_code\":200,\"resp_duration\":\"1.345193477s\",\"resp_duration_ms\":1345.193477,\"resp_latency\":\"1.345069502s\",\"resp_latency_ms\":1345.069502,\"resp_mebibytes_per_sec\":38683.678586006165,\"resp_size\":54922,\"system\":\"server\",\"time\":\"2019-03-04T19:21:26.895323594Z\",\"version\":\"5b82fdcddaf0286e7fec3a5f8dbf7a67a325fd6b\"}\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Host != "host.example.org" {
			panic("unexpected msg.Host")
		}
		if msg.JSONValues["host"] != "myhost.example.org" {
			panic("unexpected JSON value host")
		}
	}
}

func BenchmarkParserParseLongCEEWithGJSON(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"egid\":0,\"eid\":0,\"env\":\"production\",\"host\":\"myhost.example.org\",\"level\":\"info\",\"msg\":\"request complete\",\"pid\":7,\"pname\":\"/bin/myprogram\",\"req_method\":\"GET\",\"req_path\":\"/bin/myprogram/\",\"req_remote_ip\":\"172.0.0.1\",\"req_useragent\":\"my-client;go=go1.11.1\",\"resp_bytes_per_sec\":40562776957,\"resp_code\":200,\"resp_duration\":\"1.345193477s\",\"resp_duration_ms\":1345.193477,\"resp_latency\":\"1.345069502s\",\"resp_latency_ms\":1345.069502,\"resp_mebibytes_per_sec\":38683.678586006165,\"resp_size\":54922,\"system\":\"server\",\"time\":\"2019-03-04T19:21:26.895323594Z\",\"version\":\"5b82fdcddaf0286e7fec3a5f8dbf7a67a325fd6b\"}\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m, captainslog.OptionUseGJSONParser)
		if err != nil {
			panic(err)
		}
		if msg.Host != "host.example.org" {
			panic("unexpected msg.Host")
		}
		if msg.JSONValues["host"] != "myhost.example.org" {
			panic("unexpected JSON value host")
		}
	}
}

func BenchmarkParserParseCEEWithOptionDontParseJSON(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"a\":\"b\"}\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m, captainslog.OptionDontParseJSON)
		if err != nil {
			panic(err)
		}
		if msg.Host != "host.example.org" {
			panic("unexpected msg.Host")
		}
	}
}

func BenchmarkParserParseLeastLikelyTime(b *testing.B) {
	m := []byte("<38>Mon Jan  2 15:04:05 host.example.org test: hello world\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Content != " hello world" {
			panic("unexpected msg.Content")
		}
	}
}

func BenchmarkParserParseAndString(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Content != " hello world" {
			panic("unexpected msg.Content")
		}
	}
}

func BenchmarkParserParseAndBytes(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: hello world\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		msg, err := captainslog.NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}
		if msg.Content != " hello world" {
			panic("unexpected msg.Content")
		}
	}
}

func BenchmarkJSONParseNoJSON(b *testing.B) {
	m := []byte("hello world")
	val := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(m, &val)
		if err == nil {
			panic("wat")
		}
	}
}

func BenchmarkJSONCheckFirstChar(b *testing.B) {
	m := []byte("hello world")
	for i := 0; i < b.N; i++ {
		if bytes.Equal(m[0:0], []byte("{")) {
			panic("wee")
		}
	}
}

func BenchmarkParserParseInvalidDate(b *testing.B) {
	m := []byte("<191>2006-02-30T15:04:05.999999-07:00 host.example.org test: hello world\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		_, err := captainslog.NewSyslogMsgFromBytes(m)
		if err == nil {
			panic(err)
		}
	}
}

func BenchmarkParserParseInvalidSyslog(b *testing.B) {
	m := []byte("Hello I am not a syslog message\n")

	for i := 0; i < b.N; i++ {
		b.SetBytes(int64(len(m)))
		_, err := captainslog.NewSyslogMsgFromBytes(m)
		if err == nil {
			panic(err)
		}
	}
}

func generateDate(f string) string {
	t := time.Date(time.Now().Year(), time.January, 1, 15, 4, 5, 0, time.UTC) // first day of whatever the current year is
	return t.Format(f)
}
