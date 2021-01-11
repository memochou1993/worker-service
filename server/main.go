package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gw "github.com/memochou1993/worker-service/gen"
	"github.com/memochou1993/worker-service/server/handler"
	"google.golang.org/grpc"
)

const (
	grpcServerEndpoint = ":8600"
	httpServerEndpoint = ":8890"
)

func main() {
	go grpcServer()
	httpServer()
}

func grpcServer() {
	ln, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		log.Fatalln(err.Error())
	}
	s := grpc.NewServer()
	gw.RegisterServiceServer(s, new(handler.Server))
	log.Printf("Worker service gRPC server started: http://localhost%s", grpcServerEndpoint)
	if err := s.Serve(ln); err != nil {
		log.Fatalln(err.Error())
	}
}

func httpServer() {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := gw.RegisterServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts); err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Worker service HTTP server started: http://localhost%s", httpServerEndpoint)
	log.Fatalln(http.ListenAndServe(httpServerEndpoint, mux))
}
