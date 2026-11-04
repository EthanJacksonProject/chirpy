# HandlerFunc vs Handle in Go's net/http

## `http.Handler` (interface)

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

This is an **interface**. Anything with a `ServeHTTP(w, r)` method satisfies it. For example:

```go
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.HandlerFunc {
```

The parameter `next` is typed as `http.Handler` - meaning "anything that knows how to serve an HTTP request." This is what makes middleware composable: you can wrap any handler, regardless of how it was created, as long as it satisfies this interface.

## `http.HandlerFunc` (type, not interface)

```go
type HandlerFunc func(ResponseWriter, *Request)
```

This is a **named function type** with one special trick: it has a `ServeHTTP` method defined on it:

```go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

This means any plain function matching the signature `func(http.ResponseWriter, *http.Request)` can be converted to `http.HandlerFunc`, and once converted, it automatically satisfies the `http.Handler` interface. It's an adapter - it lets ordinary functions be used wherever a `Handler` is expected.

Example:

```go
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.HandlerFunc {
	newHandler := func(r http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(r, req)
	}

	return newHandler
}
```

Here, the function's return type is `http.HandlerFunc`, so Go implicitly converts `newHandler` (a plain closure) into that named type when returning it. That's why this compiles even though `newHandler` was never explicitly cast.

## `mux.Handle` vs `mux.HandleFunc`

This is where the distinction becomes practical:

```go
mux.Handle(pattern string, handler Handler)
mux.HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

- **`mux.Handle`** expects something that already satisfies the `Handler` interface (has a `ServeHTTP` method). You use this when you have a `Handler` value already - like the result of `http.FileServer(...)`, or your own middleware-wrapped handler.
- **`mux.HandleFunc`** expects a plain function matching the signature. Under the hood, it just wraps your function in `http.HandlerFunc(...)` for you and calls `Handle`.

Example usage:

```go
mux.Handle("/app/", fsHandler) // fsHandler is already an http.Handler (or HandlerFunc)
mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate) // plain function, auto-wrapped
```

## Summary table

| Concept            | What it is                 | Purpose                                                     |
| ------------------ | -------------------------- | ----------------------------------------------------------- |
| `http.Handler`     | Interface with `ServeHTTP` | The contract anything "servable" must meet                  |
| `http.HandlerFunc` | Function type + adapter    | Lets plain functions satisfy `Handler`                      |
| `mux.Handle`       | Method                     | Registers a route with something that's already a `Handler` |
| `mux.HandleFunc`   | Method                     | Registers a route with a plain function, auto-converting it |

A middleware function like `middlewareMetricsInc` takes a `Handler` in, and returns a `HandlerFunc` out - which is itself also a valid `Handler`, so it can be passed right back into another layer of middleware or registered directly with `mux.Handle`.
