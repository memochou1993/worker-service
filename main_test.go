package main

import (
	"sync"
	"testing"
)

func TestFetch(t *testing.T) {
	times := 100

	wg := sync.WaitGroup{}
	wg.Add(times)
	for i := 0; i < times; i++ {
		go func() {
			defer wg.Done()
			fetch()
		}()
	}
	wg.Wait()

	count := 0
	for _, v := range factory.attendance {
		count += v
	}

	if count != times {
		t.Fatal()
	}

	if factory.count != times {
		t.Fatal()
	}
}
