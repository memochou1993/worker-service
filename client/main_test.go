package main

import (
	"context"
	"github.com/memochou1993/worker-service/client/handler"
	gw "github.com/memochou1993/worker-service/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestGetWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := handler.Client.GetWorker(ctx, &gw.GetWorkerRequest{})
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

	if _, err := handler.Client.PutWorker(ctx, &gw.PutWorkerRequest{Number: 1}); err != nil {
		t.Fatal(err.Error())
	}
}
