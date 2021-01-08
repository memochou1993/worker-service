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

	count := summoned(0)
	for _, v := range factory.attendance {
		count += summoned(v)
	}

	if count != summoned(times) {
		t.Fatal()
	}

	if factory.summoned != summoned(times) {
		t.Fatal()
	}
}
