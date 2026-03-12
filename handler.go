package main

import (
	"cmp"
	"io"
	"net/http"
	"slices"
)

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

	writeHeaders(w, r.Header)

	if r.Body != nil {
		_, _ = w.Write([]byte("\n\n"))
		_, _ = io.Copy(w, r.Body)
	}

	if len(r.Trailer) > 0 {
		_, _ = w.Write([]byte("\n"))
		writeHeaders(w, r.Trailer)
	}
}

func writeHeaders(w io.Writer, headers http.Header) {
	for _, name := range mapKeysSorted(headers, cmp.Compare) {
		_, _ = w.Write([]byte("\n"))
		_, _ = w.Write([]byte(name))
		_, _ = w.Write([]byte(": "))

		for i, v := range headers[name] {
			if i > 0 {
				_, _ = w.Write([]byte(","))
			}

			_, _ = w.Write([]byte(v))
		}
	}
}

func mapKeysSorted[K comparable, V any](m map[K]V, cmp func(k1, k2 K) int) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	slices.SortFunc(keys, cmp)

	return keys
}
