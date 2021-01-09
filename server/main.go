package main

import (
	"context"
	pb "github.com/memochou1993/worker-server/gen"
	"github.com/memochou1993/worker-server/server/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

var (
	ws = worker.NewService().Recruit(30)
)

type service struct {
	pb.UnimplementedServiceServer
}

func (s *service) GetWorker(ctx context.Context, r *pb.GetWorkerRequest) (*pb.GetWorkerResponse, error) {
	w := ws.Dequeue()
	if w == nil {
		return &pb.GetWorkerResponse{}, status.Error(codes.NotFound, "")
	}
	return &pb.GetWorkerResponse{Number: int64(w.Number), Delay: w.Delay}, nil
}

func (s *service) PutWorker(ctx context.Context, r *pb.PutWorkerRequest) (*pb.PutWorkerResponse, error) {
	ws.Enqueue(worker.NewWorker(worker.Number(r.Number)))
	return &pb.PutWorkerResponse{}, nil
}

func main() {
	addr := "127.0.0.1:8080"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterServiceServer(s, new(service))
	if err := s.Serve(ln); err != nil {
		log.Fatalln(err.Error())
	}

	// FIXME
	// r := mux.NewRouter()
	// // 索取一個工人
	// r.HandleFunc("/worker", worker.GetWorker).Methods(http.MethodGet)
	// // 退還一個工人
	// r.HandleFunc("/worker", worker.PutWorker).Methods(http.MethodPut)
	// // 列出所有工人被傳喚記錄
	// r.HandleFunc("/workers", worker.ListWorkers).Methods(http.MethodGet)
	// // 查看特定工人被傳喚記錄
	// r.HandleFunc("/workers/{n}", worker.ShowWorker).Methods(http.MethodGet)
	//
	// log.Fatalln(http.ListenAndServe(":8890", r))
}
