package webutil

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
)

func LogHTTPRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.FromContext(r.Context())
		start := time.Now()
		r = r.WithContext(logging.WithLogger(r.Context(), logger))

		w2 := NewMyResponseWriter(w)

		next.ServeHTTP(w2, r)

		duration := time.Since(start).String()

		httpReqResLog := r.RemoteAddr + " " + "[" + start.Format(time.RFC3339) + "]" + " " + r.Method + " " + r.URL.Path + " " +
			r.Proto + " " + strconv.Itoa(w2.Status()) + " " + duration + " " + r.UserAgent()

		logger.Info(httpReqResLog)
	})
}

type MyResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func NewMyResponseWriter(w http.ResponseWriter) *MyResponseWriter {
	return &MyResponseWriter{ResponseWriter: w}
}

func (w *MyResponseWriter) Status() int {
	return w.status
}

func (w *MyResponseWriter) Write(p []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func (w *MyResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true
}
