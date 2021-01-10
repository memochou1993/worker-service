package main

import (
	"github.com/memochou1993/worker-server/server/app"
	"sync"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	service := app.NewService().Recruit(30)

	times := 100

	wg := sync.WaitGroup{}
	wg.Add(times)
	for i := 0; i < times; i++ {
		go func() {
			defer wg.Done()
			fetch(service)
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

func fetch(s *app.Service) {
	if w := s.Dequeue(); w != nil {
		time.Sleep(time.Duration(w.Delay) * time.Microsecond)
		s.Enqueue(w)
		return
	}
	time.Sleep(time.Second)
	fetch(s)
}
