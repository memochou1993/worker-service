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

// Number represents an ID of the worker.
type Number int64

// Summoned represents a number of times that a worker summoned by client.
type Summoned int64

// Worker represents a worker.
type Worker struct {
	Number `json:"Number"`
	Delay  int64 `json:"delay"`
}

// Service represents a service.
type Service struct {
	Major      chan *Worker
	Minor      chan *Worker
	Attendance map[Number]Summoned
	Summoned
}

// Enqueue enqueues a Worker instance.
func (s *Service) Enqueue(w *Worker) bool {
	worker := NewWorker(w.Number, options.Worker().SetMaxDelay(10))
	select {
	case s.Major <- worker:
		return true
	case s.Minor <- worker:
		return true
	default:
		return false
	}
}

// Dequeue dequeues a Worker instance.
func (s *Service) Dequeue() *Worker {
	select {
	case w := <-s.Major:
		s.log(*w)
		return w
	case w := <-s.Minor:
		s.log(*w)
		return w
	default:
		return nil
	}
}

func (s *Service) recruit(n int) {
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 1; i <= n; i++ {
		go func(i int) {
			defer wg.Done()
			s.Enqueue(NewWorker(Number(i), options.Worker().SetMaxDelay(10)))
		}(i)
	}
	wg.Wait()
}

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

// NewService creates a new Service instance.
func NewService(opts ...*options.ServiceOptions) *Service {
	sOpts := options.MergeServiceOptions(opts...)
	s := &Service{
		Major:      make(chan *Worker, *sOpts.MaxWorkers/2),
		Minor:      make(chan *Worker, *sOpts.MaxWorkers/2),
		Attendance: make(map[Number]Summoned),
	}
	s.recruit(*sOpts.MaxWorkers)
	return s
}

// NewWorker creates a new Worker instance.
func NewWorker(n Number, opts ...*options.WorkerOptions) *Worker {
	wOpts := options.MergeWorkerOptions(opts...)
	return &Worker{
		Number: n,
		Delay:  int64(rand.Intn(*wOpts.MaxDelay + 1)),
	}
}
