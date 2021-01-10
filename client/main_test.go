package main

import (
	"context"
	"github.com/memochou1993/worker-server/client/handler"
	gw "github.com/memochou1993/worker-server/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestGetWorker(t *testing.T) {
	_, err := handler.Client.GetWorker(context.Background(), &gw.GetWorkerRequest{})
	s, ok := status.FromError(err)
	if !ok {
		t.Fatal(err.Error())
	}
	if s.Code() == codes.NotFound {
		t.Fatal(err.Error())
	}
}

func TestPutWorker(t *testing.T) {
	if _, err := handler.Client.PutWorker(context.Background(), &gw.PutWorkerRequest{Number: 1}); err != nil {
		t.Fatal(err.Error())
	}
}
