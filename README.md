[![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/crazygreenpenguin/beats-output-http/master/LICENSE)
[![Build Status](https://travis-ci.org/crazygreenpenguin/beats-output-http.svg?branch=master)](https://travis-ci.com/crazygreenpenguin/beats-output-http)
[![Go Report Card](https://goreportcard.com/badge/github.com/crazygreenpenguin/beats-output-http)](https://goreportcard.com/report/github.com/crazygreenpenguin/beats-output-http)
# beats-output-http
HTTP output producer for the Elastic Beats framework
beats-output-http. Output for the Elastic Beats platform that simply
POSTs events to an HTTP endpoint.

Compatibilities
=====
This output require v7 beat API for older version try using https://github.com/raboof/beats-output-http

Compatibility testing conducted with beats 7.8 and output version v0.0.6

Attention
=====

Not using pre-release version! It's only for tests.

Usage
=====

To add support for this output plugin to a beat, you
have to import this plugin into your main beats package,
like this:

```
package main

import (
	"os"
	_ "github.com/crazygreenpenguin/beats-output-http"
	"github.com/elastic/beats/v7/filebeat/cmd"
)


func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```
for filebeat it's be filebeat/main.go fo example

Plugin configuration
=====

Then configure the http output plugin in filebeat.yaml:

```yaml
output.http:
  only_fields: true
  # only fields set force json codec for body and send
  # only Fields content, no @metadata (from version 0.0.7 @timestamp saving)
  # default=false
  add_fields:
    field1: 123
    field2: "fedggd"
  #add fields in out message, only if only_fields = true from v0.0.7
  url: 'http://some.example.com:80/foo'
  # URL for sending POST request
  max_retries: -1
  # How many retry fail send: -1=infinite, 0=no retry, default=-1
  compression: false
  # Use HTTP client compression? default=false
  keep_alive: true
  # HTTP KeepAlive using default=true
  max_idle_conns: 1
  # Controls the maximum number of idle (keep-alive) must be 1 and greater
  # Default=1
  idle_conn_timeout: 0
  # Is the maximum amount of time in seconds an idle (keep-alive)
  # connection will remain idle before closing itself.
  # Zero means no limit
  # Default=0
  response_header_timeout: 100
  # Specifies the amount of time in milliseconds to wait for a server's response
  # headers after fully writing the request (including its body, if any).
  # This time does not include the time to read the response body.
  # default=100ms
  username: 'test'
  # Basic Auth username if empty or none Authorization header not set
  password: 'password'
  # Basic Auth password if auth set. If auth set^ can't be empty!
```
