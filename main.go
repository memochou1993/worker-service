package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	mutex   = &sync.Mutex{}
	factory = NewFactory()
)

type Factory struct {
	Workers    chan *Worker
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
	// 模擬 client
	for i := 0; i < 100; i++ {
		go fetch()
	}

	time.Sleep(10 * time.Second)
}

func fetch() {
	// 6. client 抽出的 Entity 需確實在 server 端消失, 並於放回後重新於 server 產生
	if w := factory.dequeue(); w != nil {
		time.Sleep(time.Duration(w.Delay) * time.Millisecond)
		log.Println(fmt.Sprintf("Number: %d, Delay: %d", w.Number, w.Delay))
		factory.enqueue(w)
		return
	}
	// 9. client 抽不到號碼需等待
	fetch()
}

func (f *Factory) record(w Worker) {
	// 更新出勤表
	mutex.Lock()
	if _, ok := f.attendance[w.Number]; ok {
		f.attendance[w.Number]++
	} else {
		f.attendance[w.Number] = 1
	}
	mutex.Unlock()

	// 4. server 於每 100 次抽出時，印出每個號碼被抽取次數
	f.count++
	if f.count%100 == 0 {
		log.Println(f.attendance)
	}
}

func (f *Factory) recruit(n int) {
	for i := 1; i <= n; i++ {
		go func(i int) {
			f.Workers <- NewWorker(i)
		}(i)
	}
}

func (f *Factory) dequeue() *Worker {
	// 7. 號碼被 client 抽出期間，不可再被抽出
	select {
	case w := <-f.Workers:
		f.record(*w)
		return w
	default:
		return nil
	}
}

func (f *Factory) enqueue(w *Worker) bool {
	select {
	// 8. client 抽出的 Delay 需每次隨機不同
	case f.Workers <- NewWorker(w.Number):
		return true
	default:
		return false
	}
}

func NewFactory() *Factory {
	return &Factory{
		Workers:    make(chan *Worker, 30),
		attendance: make(map[int]int),
	}
}

func NewWorker(n int) *Worker {
	return &Worker{
		Number: n,
		Delay:  int64(rand.Intn(10)),
	}
}
