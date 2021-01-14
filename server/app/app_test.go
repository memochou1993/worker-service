package app

import (
	"sync"
	"testing"

	"github.com/memochou1993/worker-service/server/app/options"
)

func TestEnqueue(t *testing.T) {
	number := 50
	service := NewService(options.Service().SetMaxWorkers(number))

	for i := 0; i < number/2; i++ {
		<-service.Major
		<-service.Minor
	}

	wg := sync.WaitGroup{}
	wg.Add(number)
	for i := 0; i < number; i++ {
		go func(i int) {
			defer wg.Done()
			service.Enqueue(NewWorker(Number(i + 1)))
		}(i)
	}
	wg.Wait()

	if len(service.Major) != number/2 {
		t.Fail()
	}
	if len(service.Minor) != number/2 {
		t.Fail()
	}
}

func TestDequeue(t *testing.T) {
	number := 50
	service := NewService(options.Service().SetMaxWorkers(number))

	wg := sync.WaitGroup{}
	wg.Add(number)
	for i := 0; i < number; i++ {
		go func() {
			defer wg.Done()
			if w := service.Dequeue(); w == nil {
				t.Fail()
			}
		}()
	}
	wg.Wait()

	if len(service.Major) > 0 {
		t.Fail()
	}
	if len(service.Minor) > 0 {
		t.Fail()
	}
}
