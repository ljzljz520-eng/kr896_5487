package main

import "net/http"

func serverReady(h http.Handler) bool {
	r := httptestRequest(h, "GET", "/health")
	return r == http.StatusOK
}
func httptestRequest(h http.Handler, method, path string) int {
	req, _ := http.NewRequest(method, path, nil)
	w := newRecorder()
	h.ServeHTTP(w, req)
	return w.code
}

type recorder struct{ code int }

func newRecorder() *recorder            { return &recorder{} }
func (r *recorder) Header() http.Header { return http.Header{} }
func (r *recorder) WriteHeader(c int)   { r.code = c }
func (r *recorder) Write(b []byte) (int, error) {
	if r.code == 0 {
		r.code = 200
	}
	return len(b), nil
}
