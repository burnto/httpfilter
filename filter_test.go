package httpfilter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {

	pleaseCount := 0

	please := FilterFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		pleaseCount += 1
		r.Method = r.Method + "_PLEASE"
		h.ServeHTTP(w, r)
	})

	thankyou := FilterFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		r.Method = r.Method + "_THANKYOU"
		w.Header().Set("Extra", "thanks")
		h.ServeHTTP(w, r)
	})

	cases := []struct {
		name            string
		stack           Stack
		method          string
		wantMethod      string
		wantExtraHeader string
		wantPleaseCount int
	}{
		{
			"thankyou should be applied before please",
			Stack{thankyou, please},
			"GET",
			"GET_THANKYOU_PLEASE",
			"thanks",
			1,
		},
		{
			"please should be applied before thankyou",
			Stack{please, thankyou},
			"GET",
			"GET_PLEASE_THANKYOU",
			"thanks",
			1,
		},
		{
			"please should be repeated",
			Stack{please, please},
			"GET",
			"GET_PLEASE_PLEASE",
			"",
			2,
		},
		{
			"single filter",
			Stack{please},
			"GET",
			"GET_PLEASE",
			"",
			1,
		},
		{
			"no filters should just result in the bare handler",
			Stack{},
			"GET",
			"GET",
			"",
			0,
		},
	}

	for _, tt := range cases {
		pleaseCount = 0
		t.Log(tt.name)
		var handlerCalled bool
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			assert.Exactly(t, r.Method, tt.wantMethod)
			assert.Exactly(t, w.Header().Get("Extra"), tt.wantExtraHeader)
			w.WriteHeader(http.StatusTeapot)
		})
		r, err := http.NewRequest("GET", "", nil)
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		tt.stack.FilterHTTP(w, r, handler)
		assert.Exactly(t, pleaseCount, tt.wantPleaseCount)
		assert.Exactly(t, w.Code, http.StatusTeapot)
	}

}
