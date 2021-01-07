package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

var (
	factory *Factory
)

type Factory struct {
	Workers chan *Worker
	count   int
}

type Worker struct {
	Number int
	Delay  int64
}

func init() {
	factory = NewFactory()
	factory.recruit(30)
}

func main() {
	for i := 0; i < 100; i++ {
		go fetch()
	}

	time.Sleep(1 * time.Second)
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

func (f *Factory) record() {
	f.count++
	batch := 100
	if f.count%batch == 0 {
		log.Println(100)
	}
}

func (f *Factory) recruit(n int) {
	for i := 0; i < n; i++ {
		go func(i int) {
			f.Workers <- NewWorker(i)
		}(i)
	}
}

func (f *Factory) dequeue() *Worker {
	// 7. 號碼被 client 抽出期間，不可再被抽出
	select {
	case e := <-f.Workers:
		f.record()
		return e
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
		Workers: make(chan *Worker, 30),
	}
}

func NewWorker(n int) *Worker {
	return &Worker{
		Number: n,
		Delay:  int64(rand.Intn(10)),
	}
}
