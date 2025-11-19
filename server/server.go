package server

import "net/http"

func Run() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./server")))
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
