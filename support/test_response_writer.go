package support

import (
	"bytes"
	"net/http"
)

type TestResponseWriter struct {
	Body        bytes.Buffer
	StatusCode  int
	WroteHeader bool
	header      http.Header
}

func (w *TestResponseWriter) Header() http.Header {
	if nil == w.header {
		w.header = make(map[string][]string)
	}
	return w.header
}

func (w *TestResponseWriter) Write(p []byte) (int, error) {
	return w.Body.Write(p)
}

func (w *TestResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.WroteHeader = true
}
