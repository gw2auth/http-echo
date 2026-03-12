package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	var listenAddr string
	var enableKeepAlive bool
	flag.StringVar(&listenAddr, "listenaddr", ":8080", "listen address; port only: ':8080' or with interface to bind on '127.0.0.1:8080'")
	flag.BoolVar(&enableKeepAlive, "keepalive", true, "enable keepalive; true or false")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	server := http.Server{
		Addr:    listenAddr,
		Handler: echoHandler{},
	}

	server.SetKeepAlivesEnabled(enableKeepAlive)

	if err := listenAndServe(ctx, &server); err != nil {
		log.Fatal(err)
	}
}

func listenAndServe(ctx context.Context, server *http.Server) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(errors.New("defer"))

	shutdownErrCh := make(chan error)
	go func() {
		defer close(shutdownErrCh)

		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			shutdownErrCh <- err
		}
	}()

	var err error
	if listenErr := server.ListenAndServe(); !errors.Is(listenErr, http.ErrServerClosed) {
		err = errors.Join(err, listenErr)
	}

	cancel(errors.New("ListenAndServe exit"))

	if shutdownErr := <-shutdownErrCh; shutdownErr != nil {
		err = errors.Join(err, shutdownErr)
	}

	return err
}
