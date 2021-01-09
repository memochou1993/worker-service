package app

import (
	"context"
	gw "github.com/memochou1993/worker-server/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ws = NewService().Recruit(30)
)

type Server struct {
	gw.UnimplementedServiceServer
}

func (s *Server) GetWorker(ctx context.Context, r *gw.GetWorkerRequest) (*gw.GetWorkerResponse, error) {
	w := ws.Dequeue()
	if w == nil {
		return &gw.GetWorkerResponse{}, status.Error(codes.NotFound, "")
	}
	return &gw.GetWorkerResponse{Worker: &gw.Worker{Number: int64(w.Number), Delay: w.Delay}}, nil
}

func (s *Server) PutWorker(ctx context.Context, r *gw.PutWorkerRequest) (*gw.PutWorkerResponse, error) {
	ws.Enqueue(NewWorker(Number(r.Number)))
	return &gw.PutWorkerResponse{}, nil
}

func (s *Server) ListWorker(ctx context.Context, r *gw.ListWorkerRequest) (*gw.ListWorkerResponse, error) {
	var records []*gw.Record
	for number, summoned := range ws.Attendance {
		records = append(records, &gw.Record{Number: int64(number), Summoned: int64(summoned)})
	}
	return &gw.ListWorkerResponse{Records: records}, nil
}

func (s *Server) ShowWorker(ctx context.Context, r *gw.ShowWorkerRequest) (*gw.ShowWorkerResponse, error) {
	n := r.Number
	if _, ok := ws.Attendance[Number(n)]; !ok {
		return &gw.ShowWorkerResponse{}, status.Error(codes.NotFound, "")
	}
	record := gw.Record{Number: n, Summoned: int64(ws.Attendance[Number(n)])}
	return &gw.ShowWorkerResponse{Record: &record}, nil
}
