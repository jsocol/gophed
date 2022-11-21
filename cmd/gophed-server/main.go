package main

import (
	"log"
	"net"
	"net/http"

	"github.com/jsocol/gophed/handlers"
	"github.com/jsocol/gophed/repository"
)

func main() {
	repo := repository.New()
	wh := handlers.WebfingerHandler{
		Users: repo,
	}

	mux := http.NewServeMux()

	mux.Handle("/.well-known/webfinger", &wh)

	addr := "127.0.0.1:8000"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %s\n", addr)

	if err := http.Serve(listener, mux); err != nil {
		log.Fatal(err)
	}
}
