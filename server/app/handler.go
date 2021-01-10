package app

import (
	"context"
	gw "github.com/memochou1993/worker-service/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"
)

var (
	ws = NewService().Recruit(30)
)

// Server 服務
type Server struct {
	gw.UnimplementedServiceServer
}

// GetWorker 取出工人
func (s *Server) GetWorker(ctx context.Context, r *gw.GetWorkerRequest) (*gw.GetWorkerResponse, error) {
	w := ws.Dequeue()
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
	ws.Enqueue(NewWorker(Number(r.Number)))
	return &gw.PutWorkerResponse{}, nil
}

// ListWorkers 列出工人
func (s *Server) ListWorkers(ctx context.Context, r *gw.ListWorkersRequest) (*gw.ListWorkersResponse, error) {
	var records []*gw.Record
	mutex.Lock()
	for number, summoned := range ws.Attendance {
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
	n := Number(r.Number)
	mutex.Lock()
	if _, ok := ws.Attendance[n]; !ok {
		mutex.Unlock()
		return &gw.ShowWorkerResponse{}, status.Error(codes.NotFound, "")
	}
	record := gw.Record{Number: float32(n), Summoned: float32(ws.Attendance[n])}
	mutex.Unlock()
	return &gw.ShowWorkerResponse{Worker: &record}, nil
}
