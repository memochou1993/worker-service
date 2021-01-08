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

// number 代表工人號碼
type number int64

// summoned 代表工人傳喚次數
type summoned int64

// Factory 代表工廠
type Factory struct {
	workers    chan *Worker
	attendance map[number]summoned
	summoned
}

// Worker 代表工人
type Worker struct {
	Number number `json:"number"`
	Delay  int64  `json:"delay"`
}

// Record 代表工人傳喚記錄
type Record struct {
	Number   number   `json:"number"`
	Summoned summoned `json:"summoned"`
}

// Payload 代表回應資料
type Payload struct {
	Data interface{} `json:"data"`
}

func init() {
	// 應徵工人
	factory.recruit(30)
}

func main() {
	r := mux.NewRouter()
	// 索取一個工人
	r.HandleFunc("/worker", getWorker).Methods(http.MethodGet)
	// 退還一個工人
	r.HandleFunc("/worker", putWorker).Methods(http.MethodPut)
	// 列出所有工人
	r.HandleFunc("/workers", listWorkers).Methods(http.MethodGet)
	// 查看特定工人
	r.HandleFunc("/workers/{n}", showWorker).Methods(http.MethodGet)

	log.Fatalln(http.ListenAndServe(":8890", r))
}

// 索取一個工人
func getWorker(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if worker := factory.dequeue(); worker != nil {
		response(w, http.StatusOK, worker)
		return
	}

	response(w, http.StatusNotFound, nil)
}

// 退還一個工人
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

// 列出所有工人
func listWorkers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var records []Record
	for n, u := range factory.attendance {
		records = append(records, Record{n, u})
	}

	response(w, http.StatusOK, records)
}

// 查看特定工人
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

// 回應
func response(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(Payload{Data: data}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 客戶端呼叫
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

// 更新出勤表
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

// 印出出勤表
func (f *Factory) alert() {
	// server 於每 100 次抽出時，印出每個號碼被抽取次數
	f.summoned++
	if f.summoned%100 == 0 {
		log.Println(f.attendance)
	}
}

// 放入工人
func (f *Factory) enqueue(w *Worker) bool {
	// client 抽出的 Delay 須每次隨機不同
	select {
	case f.workers <- newWorker(w.Number):
		return true
	default:
		return false
	}
}

// 取出工人
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

// 應徵工人
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

// 建立新工廠
func newFactory() *Factory {
	return &Factory{
		workers:    make(chan *Worker, 30),
		attendance: make(map[number]summoned),
	}
}

// 建立新工人
func newWorker(n number) *Worker {
	return &Worker{
		Number: n,
		Delay:  int64(rand.Intn(10)),
	}
}
