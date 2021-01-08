package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var (
	mutex   = &sync.Mutex{}
	factory = newFactory()
)

type Factory struct {
	workers    chan *Worker
	attendance map[number]used
	count      int
}

type number int64

type used int64

type Worker struct {
	Number number `json:"number"`
	Delay  int64  `json:"delay"`
}

type Record struct {
	Number number `json:"number"`
	Used   used   `json:"used"`
}

type Payload struct {
	Data interface{} `json:"data"`
}

func init() {
	factory.recruit(30)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/worker", getWorker).Methods(http.MethodGet)
	r.HandleFunc("/worker", putWorker).Methods(http.MethodPut)
	r.HandleFunc("/workers", listWorkers).Methods(http.MethodGet)
	r.HandleFunc("/workers/{n}", showWorker).Methods(http.MethodGet)

	log.Fatalln(http.ListenAndServe(":8890", r))
}

func getWorker(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if worker := factory.dequeue(); worker != nil {
		response(w, http.StatusOK, worker)
		return
	}

	response(w, http.StatusNotFound, nil)
}

func putWorker(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var worker Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		response(w, http.StatusInternalServerError, nil)
		return
	}
	factory.enqueue(newWorker(worker.Number))

	response(w, http.StatusNoContent, nil)
}

func listWorkers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var records []Record
	for n, u := range factory.attendance {
		records = append(records, Record{n, u})
	}

	response(w, http.StatusOK, records)
}

func showWorker(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	n, err := strconv.Atoi(mux.Vars(r)["n"])
	if err != nil {
		response(w, http.StatusNotFound, nil)
		return
	}
	if _, ok := factory.attendance[number(n)]; !ok {
		response(w, http.StatusNotFound, nil)
		return
	}

	record := Record{number(n), factory.attendance[number(n)]}

	response(w, http.StatusOK, record)
}

func response(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(Payload{Data: data}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func fetch() {
	// client 抽出的 Entity 須確實在 server 端消失, 並於放回後重新於 server 產生
	if w := factory.dequeue(); w != nil {
		time.Sleep(time.Duration(w.Delay) * time.Microsecond)
		// log.Println(fmt.Sprintf("Number: %d, Delay: %d", w.Number, w.Delay))
		factory.enqueue(w)
		return
	}
	// client 抽不到號碼須等待
	fetch()
}

func (f *Factory) record(w Worker) {
	// 號碼被 client 抽出後, server 需紀錄號碼被抽出次數
	mutex.Lock()
	if _, ok := f.attendance[w.Number]; ok {
		f.attendance[w.Number]++
	} else {
		f.attendance[w.Number] = 1
	}
	mutex.Unlock()
}

func (f *Factory) alert() {
	// server 於每 100 次抽出時，印出每個號碼被抽取次數
	f.count++
	if f.count%100 == 0 {
		log.Println(f.attendance)
	}
}

func (f *Factory) dequeue() *Worker {
	// 號碼被 client 抽出期間，不可再被抽出
	select {
	case w := <-f.workers:
		f.record(*w)
		f.alert()
		return w
	default:
		return nil
	}
}

func (f *Factory) enqueue(w *Worker) bool {
	// client 抽出的 Delay 須每次隨機不同
	select {
	case f.workers <- newWorker(w.Number):
		return true
	default:
		return false
	}
}

func (f *Factory) recruit(n int) {
	// server 須於被要求時，隨機決定被抽出的尚存號碼實體，不可預先排序
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 1; i <= n; i++ {
		go func(i int) {
			defer wg.Done()
			f.workers <- newWorker(number(i))
		}(i)
	}
	wg.Wait()
}

func newFactory() *Factory {
	return &Factory{
		workers:    make(chan *Worker, 30),
		attendance: make(map[number]used),
	}
}

func newWorker(n number) *Worker {
	return &Worker{
		Number: n,
		Delay:  int64(rand.Intn(10)),
	}
}
