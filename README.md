httpfilter
==========

A very slim filter (i.e. "middleware") package around Go's net/http.

```go
// Create some filters
logFilter := new(filters.Log)
gzipFilter := new(filters.GZIP)

// Stack them up
stack := httpfilter.Stack{logFilter, gzipFilter}

// Apply the stack to your handler
h := httpfilter.NewHandler(stack, myHandler)
http.Handle("/resource", h)
```
