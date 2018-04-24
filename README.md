# sseserver [![Build Status](https://travis-ci.org/politie/sseserver.svg?branch=master)](https://travis-ci.org/politie/sseserver)

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

Contributors
----------

* [Arno](https://github.com/arnobroekhof)
* [Richard](https://github.com/rkettelerij)
