package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/memochou1993/worker-server/client/handler"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/api/worker", handler.GetWorker).Methods(http.MethodGet)
	r.HandleFunc("/api/worker", handler.PutWorker).Methods(http.MethodPut)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("client/public/assets/"))))

	for i := 0; i < 100; i++ {
		handler.SummonWorker(context.Background())
	}

	log.Fatalln(http.ListenAndServe(":80", r))
}
