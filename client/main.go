package main

import (
	"github.com/gorilla/mux"
	"github.com/memochou1993/worker-service/client/handler"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/api/worker", handler.GetWorker).Methods(http.MethodGet)
	r.HandleFunc("/api/worker", handler.PutWorker).Methods(http.MethodPut)
	r.HandleFunc("/api/workers", handler.ListWorkers).Methods(http.MethodGet)
	r.HandleFunc("/api/worker/{n}", handler.ShowWorker).Methods(http.MethodGet)
	r.HandleFunc("/api/workers/summon", handler.SummonWorker).Methods(http.MethodGet)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("client/public/assets/"))))

	log.Fatalln(http.ListenAndServe(":80", r))
}
