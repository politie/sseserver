# sseserver [![Build Status](https://travis-ci.org/politie/sseserver.svg?branch=master)](https://travis-ci.org/politie/sseserver) [![Go Report Card](https://goreportcard.com/badge/github.com/politie/sseserver)](https://goreportcard.com/report/github.com/politie/sseserver)

Server to stream a JSON file over [Server-Sent-Events (SSE)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/) to a client/browser.

## Build

Install Go as described: https://golang.org/doc/install and then:

```
make dev-dependencies
make release 
```

## Test

```
make test
```

## Run

```
./sseserver-<os>-<arch> -inputFile <file-to-serve>
```

Open `examples/test_sse.html` to receive the SSE events.

### Resend file to all connected clients

```
kill -HUP <PID of sseserver>
```

### Use with Consul Template

The `sseserver` works nicely when combined with [Consul Template](https://github.com/hashicorp/consul-template) to stream information from Consul to clients/browsers.

```
./consul-template \
-template=/endpoints.ctmpl:/endpoints.json \
-exec='./sseserver -enable-syslog -context /endpoints -input-file /endpoints.json' \
-exec-reload-signal=SIGHUP
```

Contributors
----------

* [Arno](https://github.com/arnobroekhof)
* [Richard](https://github.com/rkettelerij)
