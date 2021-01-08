package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	mutex   = &sync.Mutex{}
	factory = newFactory()
)

type Factory struct {
	workers    chan *Worker
	attendance map[int]int
	count      int
}

type Worker struct {
	Number int
	Delay  int64
}

func init() {
	factory.recruit(30)
}

func main() {
	//
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
			f.workers <- newWorker(i)
		}(i)
	}
	wg.Wait()
}

func newFactory() *Factory {
	return &Factory{
		workers:    make(chan *Worker, 30),
		attendance: make(map[int]int),
	}
}

func newWorker(n int) *Worker {
	return &Worker{
		Number: n,
		Delay:  int64(rand.Intn(10)),
	}
}
