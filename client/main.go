package main

import (
	"context"
	gw "github.com/memochou1993/worker-server/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

const (
	target = ":8080"
)

var (
	client gw.ServiceClient
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn := NewClientConn(ctx, target)
	defer conn.Close()

	client = gw.NewServiceClient(conn)

	for i := 0; i < 100; i++ {
		summon(context.Background())
	}
}

func summon(ctx context.Context) {
	// 取出工人
	w, err := client.GetWorker(ctx, &gw.GetWorkerRequest{})
	s, ok := status.FromError(err)
	if !ok {
		log.Println(err.Error())
		return
	}

	// 等待
	if s.Code() == codes.NotFound {
		time.Sleep(time.Microsecond)
		log.Println("waiting...")
		summon(ctx)
		return
	}

	// 延遲
	time.Sleep(time.Duration(w.Worker.Delay) * time.Microsecond)
	log.Printf("Number: %d, Delay: %d", w.Worker.Number, w.Worker.Delay)

	// 放回工人
	client.PutWorker(ctx, &gw.PutWorkerRequest{Number: w.Worker.Number})
}

func NewClientConn(ctx context.Context, addr string) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err.Error())
	}

	return conn
}
