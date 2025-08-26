package possum

import (
	"bufio"
	"fmt"
	"net"
	"net/http"

	"github.com/mikespook/possum/log"
)

var ErrHijackerNotImplement = fmt.Errorf("underlying ResponseWriter does not implement Hijacker")

// Log is a middleware that wraps an http.HandlerFunc with request/response logging functionality.
func Log(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logWriter := &logResponseWriter{
			status: http.StatusOK,
			err:    nil,
			writer: w,
		}
		log.Trace().Any("header", r.Header).
			Str("host", r.Host).
			Str("method", r.Method).
			Str("proto", r.Proto).
			Str("remote_addr", r.RemoteAddr).
			Str("url", r.URL.String()).
			Str("user_agent", r.UserAgent()).
			Msg("request")
		next(logWriter, r)
		log.Trace().Any("header", w.Header()).Int("status", logWriter.status).Msg("response")
		log.Info().
			Str("host", r.Host).
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", logWriter.status).
			Err(logWriter.err).
			Send()
	}
}

type logResponseWriter struct {
	writer http.ResponseWriter
	err    error
	status int
}

func (lrw *logResponseWriter) Write(p []byte) (int, error) {
	return lrw.writer.Write(p)
}

func (lrw *logResponseWriter) Header() http.Header {
	return lrw.writer.Header()
}

func (lrw *logResponseWriter) WriteHeader(statusCode int) {
	lrw.status = statusCode
	lrw.writer.WriteHeader(statusCode)
}

// Hijack implements http.Hijacker.Hijack
func (lrw *logResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	// Check if the underlying writer implements http.Hijacker
	hijacker, ok := lrw.writer.(http.Hijacker)
	if !ok {
		// If not, return an error indicating it's not supported
		return nil, nil, ErrHijackerNotImplement
	}
	// If it does, delegate the Hijack call to the underlying writer
	return hijacker.Hijack()
}

func (lrw *logResponseWriter) Flush() {
	if flusher, ok := lrw.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}
