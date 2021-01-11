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

func TestListWorkers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := handler.Client.ListWorkers(ctx, &gw.ListWorkersRequest{}); err != nil {
		t.Fatal(err.Error())
	}
}

func TestShowWorkers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := handler.Client.PutWorker(ctx, &gw.PutWorkerRequest{Number: 100}); err != nil {
		t.Fatal(err.Error())
	}

	if _, err := handler.Client.ShowWorker(ctx, &gw.ShowWorkerRequest{Number: 100}); err != nil {
		t.Fatal(err.Error())
	}
}
