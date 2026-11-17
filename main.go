package main

import (
	"codeproccessor/server"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatalf("error: %w", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	mux := http.NewServeMux()
	server.SetHandlers(mux)

	addr := ":8080"
	server := &http.Server{Addr: addr, Handler: mux}

	errc := make(chan error, 1)
	go func() {
		log.Printf("listening addr: %s\n", addr)
		errc <- server.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	log.Print("shutting down\n")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}
