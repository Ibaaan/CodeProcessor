package main

import (
	"codeproccessor/repository/ram"
	"codeproccessor/server"
	"codeproccessor/service"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	mux := http.NewServeMux()
	repository := ram.NewRamRepository()
	service := service.NewTaskService(repository, ctx)
	server.SetHandlers(mux, &service)

	addr := ":8000"
	server := &http.Server{Addr: addr, Handler: mux}

	errc := make(chan error, 1)
	go func() {
		slog.Info("listening addr", "addr", addr)
		errc <- server.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}
