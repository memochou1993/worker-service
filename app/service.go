package app

import (
	"log"
	"math/rand"
	"sync"
)

var (
	mutex   = &sync.Mutex{}
	service = NewService()
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

// Record 工人被傳喚記錄
type Record struct {
	Number   `json:"Number"`
	Summoned `json:"Summoned"`
}

// Payload 回應資料
type Payload struct {
	Data interface{} `json:"data"`
}

func init() {
	// 應徵工人
	service.Recruit(30)
}

// 紀錄出勤表
func (s *Service) record(w Worker) {
	// 號碼被 client 抽出後, server 需紀錄號碼被抽出次數
	mutex.Lock()
	if _, ok := s.Attendance[w.Number]; ok {
		s.Attendance[w.Number]++
	} else {
		s.Attendance[w.Number] = 1
	}
	mutex.Unlock()
}

// 印出出勤表
func (s *Service) alert() {
	// server 於每 100 次抽出時，印出每個號碼被抽取次數
	s.Summoned++
	if s.Summoned%100 == 0 {
		log.Println(s.Attendance)
	}
}

// 放入工人
func (s *Service) Enqueue(w *Worker) bool {
	// client 抽出的 Delay 須每次隨機不同
	select {
	case s.Workers <- NewWorker(w.Number):
		return true
	default:
		return false
	}
}

// 取出工人
func (s *Service) Dequeue() *Worker {
	// 號碼被 client 抽出期間，不可再被抽出
	select {
	case w := <-s.Workers:
		s.record(*w)
		s.alert()
		return w
	default:
		return nil
	}
}

// 應徵工人
func (s *Service) Recruit(n int) {
	// server 須於被要求時，隨機決定被抽出的尚存號碼實體，不可預先排序
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 1; i <= n; i++ {
		go func(i int) {
			defer wg.Done()
			s.Workers <- NewWorker(Number(i))
		}(i)
	}
	wg.Wait()
}

// 建立新服務
func NewService() *Service {
	return &Service{
		Workers:    make(chan *Worker, 30),
		Attendance: make(map[Number]Summoned),
	}
}

// 建立新工人
func NewWorker(n Number) *Worker {
	return &Worker{
		Number: n,
		Delay:  int64(rand.Intn(10)),
	}
}
