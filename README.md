# echo-cdn-proxy

[![Build Status](https://github.com/eiri/echo-cdn-proxy/workflows/build/badge.svg)](https://github.com/eiri/echo-cdn-proxy/actions)

## Summary

An `echo` middleware to proxy and cache js requests to a specified CDN

## Usage

```go
package main

import (
    "net/http"
    "time"

    "github.com/eiri/echo-cdn-proxy"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    e := echo.New()

    // Configure the middleware
    cfg := cdnproxy.NewConfig("https://cdn.jsdelivr.net", "/npm")
    e.Use(cfg.Proxy)

    // Start the server
    e.Logger.Fatal(e.Start(":8000"))
}
```

```bash
$ curl http://localhost:8000/npm/vue/dist/vue.min.js
...
```

## Dev

```bash
$ git clone https://github.com/eiri/echo-cdn-proxy
$ cd echo-cdn-proxy
$ make test
go build ./...
go test -v ./...
=== RUN   TestProxy
=== RUN   TestProxy/echo_router_simple
=== RUN   TestProxy/echo_router_with_static
=== RUN   TestProxy/echo_router_with_prefix_clash
--- PASS: TestProxy (0.00s)
    --- PASS: TestProxy/echo_router_simple (0.00s)
    --- PASS: TestProxy/echo_router_with_static (0.00s)
    --- PASS: TestProxy/echo_router_with_prefix_clash (0.00s)
=== RUN   TestProxyNotFound
--- PASS: TestProxyNotFound (0.00s)
PASS
ok      github.com/eiri/echo-cdn-proxy  0.274s
```

## License

[MIT](https://github.com/eiri/echo-cdn-proxy/blob/master/LICENSE)
