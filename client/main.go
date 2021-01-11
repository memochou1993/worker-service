package main

import (
	"github.com/gobuffalo/packr/v2"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/memochou1993/worker-service/client/handler"
)

const (
	addr = ":9000"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/api/worker", handler.GetWorker).Methods(http.MethodGet)
	r.HandleFunc("/api/worker", handler.PutWorker).Methods(http.MethodPut)
	r.HandleFunc("/api/workers", handler.ListWorkers).Methods(http.MethodGet)
	r.HandleFunc("/api/workers/{n}", handler.ShowWorker).Methods(http.MethodGet)
	r.HandleFunc("/api/workers/summon/async/{a}/sync/{s}", handler.SummonWorkers).Methods(http.MethodGet)

	box := packr.New("assets", "./public/assets")
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(box)))

	log.Printf("Worker service HTTP client started: http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
