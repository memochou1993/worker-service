package main

import (
	"context"
	gw "github.com/memochou1993/worker-server/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestGetWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn := NewClientConn(ctx, target)
	defer conn.Close()

	client = gw.NewServiceClient(conn)

	_, err := client.GetWorker(context.Background(), &gw.GetWorkerRequest{})
	s, ok := status.FromError(err)
	if !ok {
		t.Fatal(err.Error())
	}
	if s.Code() == codes.NotFound {
		t.Fatal(err.Error())
	}
}

func TestPutWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn := NewClientConn(ctx, target)
	defer conn.Close()

	client = gw.NewServiceClient(conn)

	_, err := client.PutWorker(context.Background(), &gw.PutWorkerRequest{Number: 100})
	if err != nil {
		t.Fatal(err.Error())
	}
}
