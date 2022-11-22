package middleware

import (
	"log"
	"net/http"
	"time"
)

type wrapper struct {
	http.ResponseWriter
	written int
	status  int
}

func (w *wrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *wrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

func getLevel(code int) string {
	if code == 0 {
		return "ERROR"
	}
	if code < 300 {
		return "INFO"
	}
	if code < 400 {
		return "INFO"
	}
	if code < 500 {
		return "WARN"
	}
	return "ERROR"
}

type Log struct {
	Target http.Handler
}

func (l Log) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	wr := &wrapper{
		ResponseWriter: w,
	}
	defer func() {
		dt := time.Since(start).Seconds()
		level := getLevel(wr.status)
		log.Printf("[%s] %s %s [%d] %0.3fs\n", level, r.Method, r.URL.Path, wr.status, dt)
	}()

	l.Target.ServeHTTP(wr, r)
}
