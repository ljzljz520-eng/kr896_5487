package main

import (
	"net/http"
	"ruralfolk/api"
	"ruralfolk/service"
	"ruralfolk/store"
	"testing"
)

func TestServerReady(t *testing.T) {
	s, e := store.Open(":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if !serverReady(api.Routes(service.New(s))) {
		t.Fatal("server not ready")
	}
	_ = http.MethodGet
}
