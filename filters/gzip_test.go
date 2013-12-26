package filters

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {

	cases := []struct {
		want       []byte
		acceptEnc  string
		contentEnc string
	}{
		{
			[]byte("hello world!"),
			"gzip",
			"gzip",
		},
		{
			[]byte("hello world!"),
			"",
			"",
		},
		{
			[]byte("hello world!"),
			"gzop",
			"",
		},
	}
	for _, tt := range cases {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(tt.want)
		})
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "", nil)
		r.Header.Set("Accept-Encoding", tt.acceptEnc)
		assert.NoError(t, err)

		NewGzip().FilterHTTP(w, r, h)

		assert.Exactly(t, w.Header().Get("Content-Encoding"), tt.contentEnc)
		if tt.acceptEnc == "gzip" {

			gzr, err := gzip.NewReader(w.Body)
			assert.NoError(t, err)

			have, err := ioutil.ReadAll(gzr)
			assert.NoError(t, err)
			assert.Exactly(t, have, tt.want)
		} else {
			assert.Exactly(t, w.Body.Bytes(), tt.want)
		}
	}

}
