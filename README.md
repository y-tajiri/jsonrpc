# jsonrpc

[![Travis branch](https://img.shields.io/travis/osamingo/jsonrpc/master.svg)](https://travis-ci.org/osamingo/jsonrpc)
[![codecov](https://codecov.io/gh/osamingo/jsonrpc/branch/master/graph/badge.svg)](https://codecov.io/gh/osamingo/jsonrpc)
[![Go Report Card](https://goreportcard.com/badge/osamingo/jsonrpc)](https://goreportcard.com/report/osamingo/jsonrpc)
[![codebeat badge](https://codebeat.co/badges/cbd0290d-200b-4693-80dc-296d9447c35b)](https://codebeat.co/projects/github-com-osamingo-jsonrpc)
[![GoDoc](https://godoc.org/github.com/osamingo/jsonrpc?status.svg)](https://godoc.org/github.com/osamingo/jsonrpc)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/osamingo/jsonrpc/master/LICENSE)

## About

- Simple, Poetic, Pithy.
- No `reflect` package.
  - But `reflect` package is used only when invoke the debug handler.
- Support both packages `context` and `golang.org/x/net/context`.
- Support GAE/Go Standard Environment.
- Compliance with [JSON-RPC 2.0](http://www.jsonrpc.org/specification).

## Install

```
$ go get -u github.com/osamingo/jsonrpc
```

## Usage

```go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/osamingo/jsonrpc"
)

type (
	EchoHandler struct {}
	EchoParams struct {
		Name string `json:"name"`
	}
	EchoResult struct {
		Message string `json:"message"`
	}
)

var _ (jsonrpc.Handler) = (*EchoHandler)(nil)

func (h *EchoHandler)ServeJSONRPC(c context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {

	var p EchoParams
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return EchoResult{
		Message: "Hello, " + p.Name,
	}, nil
}

func init() {
	jsonrpc.RegisterMethod("Main.Echo", &EchoHandler{}, EchoParams{}, EchoResult{})
}

func main() {
	http.HandleFunc("/jrpc", jsonrpc.HandlerFunc)
	http.HandleFunc("/jrpc/debug", jsonrpc.DebugHandlerFunc)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
```

### Result

#### Invoke the Echo method

```
POST /jrpc HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 82
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/0.9.6

{
  "jsonrpc": "2.0",
  "method": "Main.Echo",
  "params": {
    "name": "John Doe"
  },
  "id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"
}

HTTP/1.1 200 OK
Content-Length: 68
Content-Type: application/json
Date: Mon, 28 Nov 2016 13:48:13 GMT

{
  "jsonrpc": "2.0",
  "result": {
    "message": "Hello, John Doe"
  },
  "id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"
}
```

#### Access to debug handler

```
GET /jrpc/debug HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8080
User-Agent: HTTPie/0.9.6



HTTP/1.1 200 OK
Content-Length: 408
Content-Type: application/json
Date: Mon, 28 Nov 2016 13:56:24 GMT

[
  {
    "handler": "EchoHandler",
    "name": "Main.Echo",
    "params": {
      "$ref": "#/definitions/EchoParams",
      "$schema": "http://json-schema.org/draft-04/schema#",
      "definitions": {
        "EchoParams": {
          "additionalProperties": false,
          "properties": {
            "name": {
              "type": "string"
            }
          },
          "required": [
            "name"
          ],
          "type": "object"
        }
      }
    },
    "result": {
      "$ref": "#/definitions/EchoResult",
      "$schema": "http://json-schema.org/draft-04/schema#",
      "definitions": {
        "EchoResult": {
          "additionalProperties": false,
          "properties": {
            "message": {
              "type": "string"
            }
          },
          "required": [
            "message"
          ],
          "type": "object"
        }
      }
    }
  }
]
```

## License

Released under the [MIT License](https://github.com/osamingo/jsonrpc/blob/master/LICENSE).
