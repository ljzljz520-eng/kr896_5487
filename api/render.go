package api

import (
	"encoding/json"
	"net/http"
	"ruralfolk/domain"
	"strings"
)

type ErrorBody struct {
	Code, Message string
	Fields        map[string]string
}

func writeError(w http.ResponseWriter, status int, code, message string, fields map[string]string) {
	writeJSON(w, ErrorBody{Code: code, Message: message, Fields: fields}, status)
}
func decodeExhibit(r *http.Request) (domain.Exhibit, error) {
	var e domain.Exhibit
	err := json.NewDecoder(r.Body).Decode(&e)
	return e, err
}
func decodeBooking(r *http.Request) (domain.Booking, error) {
	var b domain.Booking
	err := json.NewDecoder(r.Body).Decode(&b)
	return b, err
}
func decodeGuestbook(r *http.Request) (domain.GuestbookEntry, error) {
	var g domain.GuestbookEntry
	err := json.NewDecoder(r.Body).Decode(&g)
	return g, err
}
func methodAllowed(r *http.Request, methods ...string) bool {
	for _, m := range methods {
		if r.Method == m {
			return true
		}
	}
	return false
}
func setLocation(w http.ResponseWriter, path string) {
	if path != "" {
		w.Header().Set("Location", path)
	}
}
func noCache(w http.ResponseWriter) { w.Header().Set("Cache-Control", "no-store") }
func acceptJSON(r *http.Request) bool {
	return r.Header.Get("Accept") == "" || r.Header.Get("Accept") == "application/json" || r.Header.Get("Accept") == "*/*"
}
func contentTypeJSON(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "" || r.Header.Get("Content-Type") == "application/json"
}
func statusText(status int) string {
	switch status {
	case 200:
		return "ok"
	case 201:
		return "created"
	case 400:
		return "bad_request"
	case 403:
		return "forbidden"
	case 404:
		return "not_found"
	case 405:
		return "method_not_allowed"
	case 500:
		return "internal_error"
	}
	return "unknown"
}
func pageLimit(raw string, defaultValue, max int) int {
	n := 0
	for _, c := range raw {
		if c < '0' || c > '9' {
			return defaultValue
		}
		n = n*10 + int(c-'0')
		if n > max {
			return max
		}
	}
	if n == 0 {
		return defaultValue
	}
	return n
}
func writeCollection(w http.ResponseWriter, items any) {
	noCache(w)
	writeJSON(w, items, http.StatusOK)
}

func requestID(r *http.Request) string {
	if v := r.Header.Get("X-Request-ID"); v != "" {
		return v
	}
	return "fixture-request"
}
func isMutation(r *http.Request) bool {
	return r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete
}
func routeName(r *http.Request) string {
	if r.URL.Path == "" {
		return "/"
	}
	return r.URL.Path
}
func queryValue(r *http.Request, key, fallback string) string {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	return v
}
func pathID(r *http.Request) string {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
func writeCreated(w http.ResponseWriter, v any, location string) {
	setLocation(w, location)
	writeJSON(w, v, http.StatusCreated)
}
