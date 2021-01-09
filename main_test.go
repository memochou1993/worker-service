package main

import (
	"github.com/memochou1993/worker-server/app"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	service = app.NewService()
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	service.Recruit(30)
}

func TestFetch(t *testing.T) {
	times := 100

	// 重複輪迴抽 100 次
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
	// client 抽出的 Entity 須確實在 server 端消失, 並於放回後重新於 server 產生
	if w := s.Dequeue(); w != nil {
		time.Sleep(time.Duration(w.Delay) * time.Microsecond)
		// log.Println(fmt.Sprintf("Number: %d, Delay: %d", w.Number, w.Delay))
		s.Enqueue(w)
		return
	}
	// client 抽不到號碼須等待
	time.Sleep(time.Microsecond)
	fetch(s)
}
