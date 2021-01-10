package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gw "github.com/memochou1993/worker-service/gen"
	"github.com/memochou1993/worker-service/server/app"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

const (
	grpcServerEndpoint = ":8080"
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
	gw.RegisterServiceServer(s, new(app.Server))
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
	log.Fatalln(http.ListenAndServe(httpServerEndpoint, mux))
}
