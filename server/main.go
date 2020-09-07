package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/mux", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	s := &http.Server{
		Addr:         ":8123",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			log.Printf("new connection from %s", c.RemoteAddr())
			go func() {
				time.Sleep(10 * time.Second)
				c.Close()
			}()
			return ctx
		},
	}
	log.Fatal(s.ListenAndServe())
}
