package app

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	mutex = &sync.Mutex{}
)

// Number 工人號碼
type Number int64

// Summoned 工人被傳喚次數
type Summoned int64

// Service 服務
type Service struct {
	Workers    chan *Worker
	Attendance map[Number]Summoned
	Summoned
}

// Worker 工人
type Worker struct {
	Number `json:"Number"`
	Delay  int64 `json:"delay"`
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

// Recruit 應徵工人
func (s *Service) Recruit(n int) *Service {
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 1; i <= n; i++ {
		go func(i int) {
			defer wg.Done()
			s.Workers <- NewWorker(Number(i))
		}(i)
	}
	wg.Wait()
	return s
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
func NewService() *Service {
	return &Service{
		Workers:    make(chan *Worker, 30),
		Attendance: make(map[Number]Summoned),
	}
}

// NewWorker 建立新工人
func NewWorker(n Number) *Worker {
	rand.Seed(time.Now().UnixNano())
	return &Worker{
		Number: n,
		Delay:  int64(rand.Intn(11)),
	}
}
