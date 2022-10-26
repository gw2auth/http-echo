package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
)

func mapKeysSorted[K comparable, V any](m map[K]V, less func(k1, k2 K) bool) <-chan K {
	sortedKeys := make([]K, 0, len(m))

	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Slice(sortedKeys, func(i, j int) bool { return less(sortedKeys[i], sortedKeys[j]) })

	ch := make(chan K)

	go func() {
		for _, k := range sortedKeys {
			ch <- k
		}

		close(ch)
	}()

	return ch
}

type echoHandler struct{}

func (hf echoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	_, _ = w.Write([]byte(r.Method))
	_, _ = w.Write([]byte(" "))
	_, _ = w.Write([]byte(r.RequestURI))
	_, _ = w.Write([]byte(" "))
	_, _ = w.Write([]byte(r.Proto))
	_, _ = w.Write([]byte("\n"))
	_, _ = w.Write([]byte("Host: "))
	_, _ = w.Write([]byte(r.Host))

	for name := range mapKeysSorted(r.Header, func(k1, k2 string) bool { return k1 < k2 }) {
		_, _ = w.Write([]byte("\n"))
		_, _ = w.Write([]byte(name))
		_, _ = w.Write([]byte(": "))

		for i, v := range r.Header[name] {
			if i > 0 {
				_, _ = w.Write([]byte(","))
			}

			_, _ = w.Write([]byte(v))
		}
	}

	if r.Body != nil {
		_, _ = w.Write([]byte("\n\n"))
		_, _ = io.Copy(w, r.Body)
	}
}

func main() {
	var listenAddr string
	var enableKeepAlive bool
	flag.StringVar(&listenAddr, "listenaddr", ":8080", "listen address; port only: ':8080' or with interface to bind on '127.0.0.1:8080'")
	flag.BoolVar(&enableKeepAlive, "keepalive", true, "enable keepalive; true or false")
	flag.Parse()

	server := http.Server{
		Addr:    listenAddr,
		Handler: echoHandler{},
	}

	server.SetKeepAlivesEnabled(enableKeepAlive)

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}

		close(closed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-closed
}
