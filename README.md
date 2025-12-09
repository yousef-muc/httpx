# httpx - A lightweight, minimalistic HTTP client for Go
![Hero Image](./artifacts/general/img/stage.png)

A lightweight **minimalistic HTTP client** written in Go.  
`httpx` provides a clean, minimalistic wrapper around Go's native `net/http` package—
with support for global headers, per-request headers, timeouts and automatic JSON/XML request body encoding.

---
## 🚀 Installation


```bash
go get github.com/yousef-muc/httpx@v0.2.1
```

[![Go Reference](https://pkg.go.dev/badge/github.com/yousef-muc/token-captcha.svg)](https://pkg.go.dev/github.com/yousef-muc/token-captcha) 

**Keywords:** go http client, axios for go, golang http wrapper, go fetch api, http utility

---
## 🧠 Overview

`httpx` is a **minimalistic HTTP client** that simplifies making requests in Go.  
It provides an ergonomic API inspired by Axios, while staying fully compatible with
Go's standard `net/http` package.

Key concepts:

- Create a client instance
- Set optional **global headers**
- Pass per-request headers
- Provide a body (any type) → automatically serialized (JSON/XML)
- Perform HTTP methods (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`)

The client remains **lightweight**, avoids unnecessary abstractions, and keeps
full control in the developer’s hands.

---
## ⚙️ Features

- 🔄 Global client headers
- 🧩 Per-request headers with automatic merging
- 📝 JSON & XML body serialization based on `Content-Type`
- 🚀 Simple, clean API inspired by Axios
- 🧪 Works seamlessly with Go's native `http.Response`
- 🌱 Zero dependencies (uses only stdlib)

---
## 🖥️ Example: Usage in Go

Below is a minimal example showing how to perform GET and POST requests.

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "io"

    "github.com/yousef-muc/httpx"
)

func main() {
    client := httpx.New()

    // Set global headers
    defaultHeaders := make(http.Header)
    defaultHeaders.Set("Authorization", "Bearer ABC-123")
    client.SetHeaders(defaultHeaders)

    // --- GET request example ---
    reqHeaders := make(http.Header)
    reqHeaders.Set("Content-Type", "application/json")

    res, err := client.Get("https://dummyjson.com/carts", reqHeaders)
    if err != nil { log.Fatal(err) }
    defer res.Body.Close()

    body, _ := io.ReadAll(res.Body)
    log.Println(string(body))


    // --- POST request example ---
    type User struct {
        Firstname string `json:"firstname"`
        Lastname  string `json:"lastname"`
    }

    user := User{Firstname: "Yousef", Lastname: "Hejazi"}
    reqHeaders.Set("Content-Type", "application/json")

    res, err = client.Post("https://dummyjson.com/carts/add", reqHeaders, user)
    if err != nil { log.Fatal(err) }
    defer res.Body.Close()

    body, _ = io.ReadAll(res.Body)
    log.Println(string(body))
}
```

---
## ⚙️ Header Behavior

`httpx` merges **global headers** (set via `SetHeaders`) with **per-request headers**:

- Per-request headers **override** global headers when keys match
- Only the first value of each header key is used (simple, predictable behavior)

---
## 📦 Automatic Body Encoding

`httpx` serializes request bodies based on their `Content-Type` header:

| Content-Type          | Encoding Method     |
| --------------------- | ------------------- |
| `application/json`    | `json.Marshal`      |
| `application/xml`     | `xml.Marshal`       |
| *(anything else)*     | defaults to JSON    |

Passing `nil` as body → sends no request body.

---
## 🧾 API Reference

### **Client creation**
```go
client := httpx.New()
```

### **Set global headers**
```go
client.SetHeaders(http.Header{"Authorization": []string{"Bearer XYZ"}})
```

### **Send requests**
```go
client.Get(url, headers)
client.Post(url, headers, body)
client.Put(url, headers, body)
client.Patch(url, headers, body)
client.Delete(url, headers)
```

All methods return:
```go
(*http.Response, error)
```

---
## 🔒 Notes

- `httpx` does **not** replace Go's `http.Client`; it wraps it for convenience.
- You can still modify the response, stream data, inspect headers, etc.
- No global state — every `client` instance is isolated.

---
## 🧩 Why httpx?

Go's `net/http` package is powerful but verbose.

`httpx` aims to:

- Reduce boilerplate
- Provide a cleaner API similar to Axios/fetch
- Offer structured, predictable header and body behavior
- Remain extremely lightweight and dependency-free

Perfect for microservices, REST API clients, and CLI tools.

---

## 🧾 License

MIT License

Copyright (c) 2025 yousef-muc

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
