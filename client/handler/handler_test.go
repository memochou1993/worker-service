package handler

import (
	"context"
	"testing"
	"time"

	gw "github.com/memochou1993/worker-service/gen"
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client = gw.NewServiceClient(NewClientConn(ctx, target))
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
