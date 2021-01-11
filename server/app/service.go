package app

import (
	"log"
	"math/rand"
	"sync"
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

// ServiceOptions 服務選項
type ServiceOptions struct {
	MaxWorkers int
}

// SetMaxWorkers 設置工人最大數量
func (s *ServiceOptions) SetMaxWorkers(max int) *ServiceOptions {
	s.MaxWorkers = max
	return s
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
	case s.Workers <- NewWorker(w.Number):
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
		s.Workers <- NewWorker(Number(i))
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

// NewServiceOptions 建立新服務選項
func NewServiceOptions() *ServiceOptions {
	return &ServiceOptions{
		MaxWorkers: 10,
	}
}

// NewService 建立新服務
func NewService(opts *ServiceOptions) *Service {
	s := &Service{
		Workers:    make(chan *Worker, opts.MaxWorkers),
		Attendance: make(map[Number]Summoned),
	}
	s.recruit(opts.MaxWorkers)
	return s
}

// NewWorker 建立新工人
func NewWorker(n Number) *Worker {
	return &Worker{
		Number: n,
		Delay:  int64(rand.Intn(11)),
	}
}
