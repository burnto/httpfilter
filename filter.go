// httpfilter provides a very slim filter (i.e. "middleware") package around Go's net/http.
//
// Filters can be applied to a http.Handler to produce a new http.Handler.
//
// Filters can be chained into a Stack, which is a composite filter.
package httpfilter

import "net/http"

// Filter wraps an http.Handler with its FilterHTTP function
type Filter interface {
	FilterHTTP(http.ResponseWriter, *http.Request, http.Handler)
}

// FilterFunc is a FilterHTTP func that knows how to call itself.
// Offered for convenience.
type FilterFunc func(http.ResponseWriter, *http.Request, http.Handler)

// FilterHTTP calls its FilterFunc
func (f FilterFunc) FilterHTTP(w http.ResponseWriter, r *http.Request, h http.Handler) {
	f(w, r, h)
}

// NewHandler applies a filter to an http.Handler to produce a
// new, wrapped http.Handler
func NewHandler(f Filter, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.FilterHTTP(w, r, h)
	})
}

// Stack is a list of filters, to be applied in order
type Stack []Filter

// FilterHTTP wraps an http.Handler with the stack of Filters
func (s Stack) FilterHTTP(w http.ResponseWriter, r *http.Request, h http.Handler) {
	if len(s) == 0 {
		h.ServeHTTP(w, r)
	} else {
		s[0].FilterHTTP(w, r, NewHandler(s[1:], h))
	}
}
