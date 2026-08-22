package api

import (
	"io"
	"net/http/httptest"
	"ruralfolk/service"
	"ruralfolk/store"
	"testing"
)

func TestAPIHealth(t *testing.T) {
	s, e := store.Open(":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	Routes(service.New(s)).ServeHTTP(r, req)
	if r.Code != 200 {
		t.Fatalf("health %d", r.Code)
	}
}

func TestStaticEntry(t *testing.T) {
	s, e := store.Open(":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	WithStatic(Routes(service.New(s)), "../web").ServeHTTP(r, req)
	body, e := io.ReadAll(r.Result().Body)
	if e != nil || r.Code != 200 || len(body) == 0 {
		t.Fatalf("static entry unavailable: status=%d bytes=%d err=%v", r.Code, len(body), e)
	}
}
