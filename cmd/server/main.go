package main

import (
	"log"
	"net/http"
	"ruralfolk/api"
	"ruralfolk/config"
	"ruralfolk/service"
	"ruralfolk/store"
)

func main() {
	c := config.FromEnv()
	if err := c.Validate(); err != nil {
		log.Fatal(err)
	}
	s, err := store.Open(c.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	handler := api.WithHeaders(api.WithBodyLimit(api.WithStatic(api.Routes(service.New(s)), c.StaticDir), c.MaxBodyBytes))
	log.Printf("rural folk exhibition listening on %s", c.Address)
	log.Fatal(http.ListenAndServe(c.Address, handler))
}
