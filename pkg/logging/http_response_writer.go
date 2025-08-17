package logging

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPResponseWriter wraps http.ResponseWriter to capture status code and body
type HTTPResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (w *HTTPResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *HTTPResponseWriter) Write(data []byte) (int, error) {
	if w.body != nil {
		w.body.Write(data)
	}
	return w.ResponseWriter.Write(data)
}

// Add missing methods for gin.ResponseWriter compatibility
func (w *HTTPResponseWriter) CloseNotify() <-chan bool {
	if cn, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}

func (w *HTTPResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *HTTPResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacking not supported")
}

func (w *HTTPResponseWriter) Size() int {
	if sizer, ok := w.ResponseWriter.(gin.ResponseWriter); ok {
		return sizer.Size()
	}
	return 0
}

func (w *HTTPResponseWriter) Status() int {
	return w.statusCode
}

func (w *HTTPResponseWriter) WriteHeaderNow() {
	if w.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}
}

func (w *HTTPResponseWriter) Written() bool {
	return w.statusCode != 0
}

func (w *HTTPResponseWriter) WriteString(s string) (int, error) {
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.Write([]byte(s))
}

func (w *HTTPResponseWriter) Pusher() http.Pusher {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
