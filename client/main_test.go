package main

import (
	"context"
	"github.com/memochou1993/worker-service/client/handler"
	gw "github.com/memochou1993/worker-service/gen"
	"testing"
	"time"
)

func TestGetWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := handler.Client.GetWorker(ctx, &gw.GetWorkerRequest{}); err != nil {
		t.Fatal(err.Error())
	}
}

func TestPutWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := handler.Client.PutWorker(ctx, &gw.PutWorkerRequest{Number: 100}); err != nil {
		t.Fatal(err.Error())
	}
}
