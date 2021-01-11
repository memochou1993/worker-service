package app

import (
	"log"
	"math/rand"
	"sync"

	"github.com/memochou1993/worker-service/server/app/options"
)

var (
	mutex = &sync.Mutex{}
)

// Number 工人號碼
type Number int64

// Summoned 工人被傳喚次數
type Summoned int64

// Worker 工人
type Worker struct {
	Number `json:"Number"`
	Delay  int64 `json:"delay"`
}

// Service 服務
type Service struct {
	Workers    chan *Worker
	Attendance map[Number]Summoned
	Summoned
}

// Enqueue 放入工人
func (s *Service) Enqueue(w *Worker) bool {
	select {
	case s.Workers <- NewWorker(w.Number, options.Worker().SetMaxDelay(10)):
		return true
	default:
		return false
	}
}

// Dequeue 取出工人
func (s *Service) Dequeue() *Worker {
	select {
	case w := <-s.Workers:
		s.log(*w)
		return w
	default:
		return nil
	}
}

// recruit 填充工人
func (s *Service) recruit(n int) {
	for i := 1; i <= n; i++ {
		s.Workers <- NewWorker(Number(i), options.Worker().SetMaxDelay(10))
	}
}

// log 紀錄出勤表
func (s *Service) log(w Worker) {
	mutex.Lock()
	if _, ok := s.Attendance[w.Number]; ok {
		s.Attendance[w.Number]++
	} else {
		s.Attendance[w.Number] = 1
	}
	s.Summoned++
	if s.Summoned%100 == 0 {
		log.Println(s.Attendance)
	}
	mutex.Unlock()
}

// NewService 建立新服務
func NewService(opts ...*options.ServiceOptions) *Service {
	sOpts := options.MergeServiceOptions(opts...)
	s := &Service{
		Workers:    make(chan *Worker, *sOpts.MaxWorkers),
		Attendance: make(map[Number]Summoned),
	}
	s.recruit(*sOpts.MaxWorkers)
	return s
}

// NewWorker 建立新工人
func NewWorker(n Number, opts ...*options.WorkerOptions) *Worker {
	wOpts := options.MergeWorkerOptions(opts...)
	return &Worker{
		Number: n,
		Delay:  int64(rand.Intn(*wOpts.MaxDelay + 1)),
	}
}
