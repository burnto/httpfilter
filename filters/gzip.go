package filters

import (
	"compress/gzip"
	"io"
	"net/http"
	"regexp"

	"github.com/burnto/httpfilter"
)

var gzipRegexp *regexp.Regexp

func init() {
	gzipRegexp = regexp.MustCompile(`\bgzip\b`)
}

func gzipFilter(w http.ResponseWriter, r *http.Request, h http.Handler) {
	accept := r.Header.Get("Accept-Encoding")
	if gzipRegexp.MatchString(accept) {
		w.Header().Set("Content-Encoding", "gzip")

		gzw := gzip.NewWriter(w)
		defer gzw.Close()

		gzrw := &gzipResponseWriter{gzw, w}
		h.ServeHTTP(gzrw, r)
		gzw.Flush()
	} else {
		h.ServeHTTP(w, r)
	}
}

func NewGzip() httpfilter.Filter {
	return httpfilter.FilterFunc(gzipFilter)
}

// gzipResponseWriter embeds an http.ResponseWriter, but overrides
// its Write method to use a different Writer.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
