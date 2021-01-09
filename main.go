package main

import (
	"github.com/gorilla/mux"
	"github.com/memochou1993/worker-server/app"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	// 索取一個工人
	r.HandleFunc("/worker", app.GetWorker).Methods(http.MethodGet)
	// 退還一個工人
	r.HandleFunc("/worker", app.PutWorker).Methods(http.MethodPut)
	// 列出所有工人被傳喚記錄
	r.HandleFunc("/workers", app.ListWorkers).Methods(http.MethodGet)
	// 查看特定工人被傳喚記錄
	r.HandleFunc("/workers/{n}", app.ShowWorker).Methods(http.MethodGet)

	log.Fatalln(http.ListenAndServe(":8890", r))
}
