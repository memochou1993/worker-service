package main

import (
	"sync"
	"testing"
	"time"

	"github.com/memochou1993/worker-service/server/app"
	"github.com/memochou1993/worker-service/server/options"
)

func TestSummon(t *testing.T) {
	service := app.NewService(options.Service().SetMaxWorkers(50))

	times := 100

	wg := sync.WaitGroup{}
	wg.Add(times)
	for i := 0; i < times; i++ {
		go func() {
			defer wg.Done()
			summon(service)
		}()
	}
	wg.Wait()

	summoned := app.Summoned(0)
	for _, v := range service.Attendance {
		summoned += v
	}
	if summoned != app.Summoned(times) {
		t.Fatal()
	}

	if service.Summoned != app.Summoned(times) {
		t.Fatal()
	}
}

func summon(s *app.Service) {
	if w := s.Dequeue(); w != nil {
		time.Sleep(time.Duration(w.Delay) * time.Microsecond)
		s.Enqueue(w)
		return
	}
	time.Sleep(time.Second)
	summon(s)
}
