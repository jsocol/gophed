package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jsocol/gophed/handlers"
	"github.com/jsocol/gophed/middleware"
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

	srv := http.Server{
		Handler:           &middleware.Log{Target: mux},
		ReadHeaderTimeout: 1 * time.Second,
	}

	serverClosed := make(chan struct{}, 1)
	go func() {
		log.Printf("Listening on %s...\n", addr)
		if err := srv.Serve(listener); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
		close(serverClosed)
	}()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	<-serverClosed
}
