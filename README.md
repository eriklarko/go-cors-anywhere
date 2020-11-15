# go-cors-anywhere

A proxy adding CORS headers to the proxied requests. Useful if you don't control the API you're using and it doesn't handle CORS the way you want.

This is a minimal implementation not intended for production use


## Without go-cors-anywhere
```
+---------+   GET example.com        +------------+
|         |------------------------->|            |
| Browser |                          | API Server |
|         |<-------------------------|            |
+---------+   Response without CORS  +------------+
              headers
```

## With go-cors-anywhere

```

+---------+   GET go-cors-anywhere.com/http://example.com  +------------------+
|         |------------------------------------------------|                  |
| Browser |                                                | go-cors-anywhere |
|         |<-----------------------------------------------|                  |
+---------+          Response with CORS headers            +------------------+
                                                             |              ^
                                                             |              |
                                                          GET example.com   | Response without
                                                             |              | CORS headers
                                                             |              |
                                                             v              |
                                                           +------------------+
                                                           |                  |
                                                           |    API Server    |
                                                           |                  |
                                                           +------------------+
```

# Run

From the root of this repository, run

```
go run main.go
```

You can of course also build it first

```
go build
./go-cors-anywhere
```

The server starts on port 8080 and is currently not configurable


# Headers added

* `Access-Control-Allow-Origin`: `"*"`
* `Access-Control-Allow-Methods`: what the browser sends the proxy in `Access-Control-Request-Method`
* `Access-Control-Allow-Headers`: what the browser sends the proxy in `Access-Control-Request-Headers`
* `Access-Control-Expose-Headers`: any non-CORS headers the API server replies with


No CORS headers returned from the API server are forwarded to the browser.
