# httpx - A lightweight, minimalistic HTTP client for Go
![Hero Image](./artifacts/general/img/stage.png)

A lightweight **minimalistic HTTP client** written in Go.  
`httpx` provides a clean, minimalistic wrapper around Go's native `net/http` package—
with support for global headers, per-request headers, timeouts and automatic JSON/XML request body encoding. 
It provides:

- An **axios-style API**
- Automatic body encoding (JSON, XML, forms, multipart)
- Query parameters via `WithParams`
- Global & per-request headers
- Configurable timeouts
- Strong response helpers (`JSON`, `Text`, `Bytes`, `XML`)

This library stays **zero-dependency**, fully predictable, and idiomatic.

---
## 🚀 Installation


```bash
go get github.com/yousef-muc/httpx@v0.3.0
```

[![Go Reference](https://pkg.go.dev/badge/github.com/yousef-muc/token-captcha.svg)](https://pkg.go.dev/github.com/yousef-muc/token-captcha) 

**Keywords:** go http client, axios for go, golang http wrapper, go fetch api, http utility

---
## 🧠 Overview

`httpx` is a **minimalistic HTTP client** that simplifies making requests in Go.  
It provides an ergonomic API inspired by axios, while staying fully compatible with
Go's standard `net/http` package.
`httpx` aims to remove boilerplate from everyday HTTP calls while keeping the full power of `net/http`.

Key concepts:

- Simple request API: `Get`, `Post`, `Put`, `Patch`, `Delete`
- Options pattern for per-request configuration
- Automatic request body encoding
- Clean error reporting (`HttpError`)
- Helper functions for reading responses

The client remains **lightweight**, avoids unnecessary abstractions, and keeps
full control in the developer’s hands.

---
## ⚙️ Features

- 🔄 Global & per-request headers
- ❓ Optional query parameters
- 📝 Automatic JSON / XML / Form / Multipart encoding
- ⏱️ Request + connection timeouts
- 🧪 Response helpers (typed decoding)
- 🌱 Zero dependencies

---

# 🧭 Request Options (core concept)

Every request optionally accepts `Option` modifiers.

### Supported options:

- `WithHeaders(http.Header)`
- `WithParams(map[string]string)`
- `WithBody(any)`

Example:

```go
client.Get(
    "https://api.com/users",
    httpx.WithParams(map[string]string{"limit": "10"}),
)
```

---

# 📘 Usage Examples

The following section is split into:

1. **Simple (minimal) examples** – the most common usage
2. **Advanced examples** – using headers, parameters, and bodies

---

# 1️⃣ Simple Client Creation

### **Default client (no config)**

```go
client := httpx.New(nil)
```

### **Client with configuration**

```go
client := httpx.New(&httpx.Config{
    RequestTimeout:     5 * time.Second,
    ConnectionTimeout:  1 * time.Second,
    MaxIdleConnections: 10,
    Headers: http.Header{
        "Authorization": []string{"Bearer TOKEN"},
    },
})
```

---

# 2️⃣ Simple Request Examples

These examples demonstrate the absolute minimum required to use each method.

---

## 📗 Simple GET

```go
res, err := client.Get("https://api.com/products")
if err != nil { panic(err) }

fmt.Println(client.Text(res))
```

---

## 📘 Simple POST (JSON Body)

```go
body := map[string]any{"name": "John"}

res, err := client.Post(
    "https://api.com/users/add",
    httpx.WithBody(body),
)
if err != nil { panic(err) }

fmt.Println(client.Text(res))
```

---

## 🟦 Simple PUT

```go
res, err := client.Put(
    "https://api.com/users/1",
    httpx.WithBody(map[string]any{"active": true}),
)
fmt.Println(client.Text(res))
```

---

## 🟧 Simple PATCH

```go
res, err := client.Patch(
    "https://api.com/users/1",
    httpx.WithBody(map[string]any{"lastname": "Smith"}),
)
fmt.Println(client.Text(res))
```

---

## 🟥 Simple DELETE

```go
res, err := client.Delete("https://api.com/users/1")
fmt.Println(client.Text(res))
```

---

# 3️⃣ Advanced Examples

These examples show the full power of request options.

---

## 📗 GET with params + headers

```go
res, err := client.Get(
    "https://api.com/products",
    httpx.WithParams(map[string]string{"limit": "5"}),
    httpx.WithHeaders(http.Header{"Accept": []string{"application/json"}}),
)
if err != nil { panic(err) }
fmt.Println(client.Text(res))
```

---

## 📘 POST JSON

```go
type User struct {
    Firstname string `json:"firstname"`
    Lastname  string `json:"lastname"`
}

res, err := client.Post(
    "https://api.com/users/add",
    httpx.WithBody(User{"John", "Smith"}),
    httpx.WithHeaders(http.Header{"Content-Type": []string{"application/json"}}),
)
created, _ := httpx.JSON[User](res)
fmt.Println(created)
```

---

## 📙 POST Form

```go
res, err := client.Post(
    "https://api.com/auth/login",
    httpx.WithBody(map[string]string{
        "username": "demo",
        "password": "secret",
    }),
    httpx.WithHeaders(http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}),
)
fmt.Println(client.Text(res))
```

---

## 📕 POST Multipart Upload

```go
file := []byte("binary data here…")

res, err := client.Post(
    "https://api.com/upload",
    httpx.WithBody(map[string]any{
        "avatar": file,
        "username": "John",
    }),
    httpx.WithHeaders(http.Header{"Content-Type": []string{"multipart/form-data"}}),
)
fmt.Println(client.Text(res))
```

---

# 📦 Response Helpers

### JSON (generic)

```go
user, err := httpx.JSON[User](res)
```

### JSON (struct pointer)

```go
var user User
err := client.ReadJSON(res, &user)
```

### Text

```go
text, _ := client.Text(res)
```

### Bytes

```go
data, _ := client.Bytes(res)
```

### XML

```go
feed, _ := httpx.XML[Feed](res)
```

---

# ⚠️ Error Handling (Axios-like)

Non-2xx responses return a structured `HttpError`:

```go
res, err := client.Get("https://api.com/404")
if err != nil {
    if httpErr, ok := err.(*httpx.HttpError); ok {
        fmt.Println(httpErr.StatusCode)
        fmt.Println(string(httpErr.Body))
        fmt.Println(httpErr.Method)
        fmt.Println(httpErr.URL)
    }
}
```

---

# 🧩 Why httpx?

`httpx` removes friction from everyday HTTP operations:

- No verbose boilerplate
- No repeated marshaling logic
- Easy header & parameter handling
- Clean, consistent API
- Predictable & explicit behavior

Perfect for REST clients, CLI tools, and microservices.

---

## 🧾 License

MIT License

Copyright (c) 2025 yousef-muc

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
