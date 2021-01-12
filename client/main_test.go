package main

import (
	"context"
	"testing"
	"time"

	"github.com/memochou1993/worker-service/client/handler"
	gw "github.com/memochou1993/worker-service/gen"
)

const (
	target = ":8500"
)

var (
	client gw.ServiceClient
)

func init() {
	client = gw.NewServiceClient(handler.NewClientConn(context.Background(), target))
}

func TestGetWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := client.GetWorker(ctx, &gw.GetWorkerRequest{}); err != nil {
		t.Fatal(err.Error())
	}
}

func TestPutWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := client.PutWorker(ctx, &gw.PutWorkerRequest{Number: 100}); err != nil {
		t.Fatal(err.Error())
	}
}
