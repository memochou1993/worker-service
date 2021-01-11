package handler

import (
	"context"
	"math/rand"
	"sort"
	"sync"
	"time"

	gw "github.com/memochou1993/worker-service/gen"
	"github.com/memochou1993/worker-service/server/app"
	"github.com/memochou1993/worker-service/server/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	service *app.Service
	mutex   = &sync.Mutex{}
)

// Server 服務
type Server struct {
	gw.UnimplementedServiceServer
}

func init() {
	rand.Seed(time.Now().UnixNano())
	service = app.NewService(options.Service().SetMaxWorkers(30))
}

// GetWorker 取出工人
func (s *Server) GetWorker(ctx context.Context, r *gw.GetWorkerRequest) (*gw.GetWorkerResponse, error) {
	w := service.Dequeue()
	if w == nil {
		return &gw.GetWorkerResponse{}, status.Error(codes.NotFound, "")
	}
	return &gw.GetWorkerResponse{Worker: &gw.Worker{Number: float32(w.Number), Delay: float32(w.Delay)}}, nil
}

// PutWorker 放回工人
func (s *Server) PutWorker(ctx context.Context, r *gw.PutWorkerRequest) (*gw.PutWorkerResponse, error) {
	if r.Number < 1 {
		return &gw.PutWorkerResponse{}, status.Error(codes.InvalidArgument, "")
	}
	service.Enqueue(app.NewWorker(app.Number(r.Number), options.Worker().SetMaxDelay(10)))
	return &gw.PutWorkerResponse{}, nil
}

// ListWorkers 列出工人
func (s *Server) ListWorkers(ctx context.Context, r *gw.ListWorkersRequest) (*gw.ListWorkersResponse, error) {
	var records []*gw.Record
	mutex.Lock()
	for number, summoned := range service.Attendance {
		records = append(records, &gw.Record{Number: float32(number), Summoned: float32(summoned)})
	}
	mutex.Unlock()
	sort.Slice(records, func(i, j int) bool {
		return records[i].Number < records[j].Number
	})
	return &gw.ListWorkersResponse{Workers: records}, nil
}

// ShowWorker 查看工人
func (s *Server) ShowWorker(ctx context.Context, r *gw.ShowWorkerRequest) (*gw.ShowWorkerResponse, error) {
	n := app.Number(r.Number)
	mutex.Lock()
	if _, ok := service.Attendance[n]; !ok {
		mutex.Unlock()
		return &gw.ShowWorkerResponse{}, status.Error(codes.NotFound, "")
	}
	record := gw.Record{Number: float32(n), Summoned: float32(service.Attendance[n])}
	mutex.Unlock()
	return &gw.ShowWorkerResponse{Worker: &record}, nil
}
