package handler

import "net/http"

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
