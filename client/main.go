package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	var remoteAddr string

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	_ = &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				log.Printf("new dial to %s", addr)
				return dialer.Dial(network, addr)
			},
		},
	}

	for i := 1; i < 1024; i++ {
		time.Sleep(time.Second)
		resp, err := http.DefaultClient.Get("http://localhost:8123")
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("#%02d %v %s", i, remoteAddr, resp.Request.RemoteAddr)
	}
}
