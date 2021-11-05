package main

import "net/http"

func main() {
	ping := http.Server{
		Addr: ":8081",
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("PING"))
		}),
	}
	pong := http.Server{
		Addr: ":8082",
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("PONG"))
		}),
	}

	go ping.ListenAndServe()
	pong.ListenAndServe()
}
