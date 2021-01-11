package main

import (
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
	r.HandleFunc("/api/summon/workers/sync", handler.SummonWorkersSync).Methods(http.MethodGet)
	r.HandleFunc("/api/summon/workers/async", handler.SummonWorkersAsync).Methods(http.MethodGet)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("client/public/assets/"))))

	log.Printf("\033[1;33mWorker service HTTP client started: http://localhost%s\033[0m", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}
