package main

import (
	"context"
	"log"
	"net/http"

	"simple-repo/pkg/server"
	. "simple-repo/pkg/signals"
)

const port = 2345

func main() {
	srv, err := server.New(port)
	if err != nil {
		log.Fatalf("server error: %s\n", err.Error())
	}

	graceful := SetupSignals()

	go func() {
		log.Printf("server started on port %d\n", port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server stopped: %s\n", err.Error())
		}
		log.Println("server stopped")
	}()

	<-graceful
	log.Println("starting graceful shutdown")
	ctx := context.Background()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server shutdown error: %s\n", err.Error())
	}
	log.Println("server shutdown")
}
