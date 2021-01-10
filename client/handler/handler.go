package handler

import (
	"context"
	"encoding/json"
	gw "github.com/memochou1993/worker-server/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	target = ":8080"
)

var (
	client gw.ServiceClient
)

func init() {
	client = gw.NewServiceClient(NewClientConn(context.Background(), target))

	for i := 0; i < 100; i++ {
		summon(context.Background())
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	render(w, "index")
}

func GetWorker(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	worker, err := client.GetWorker(context.Background(), &gw.GetWorkerRequest{})
	s, ok := status.FromError(err)
	if !ok {
		response(w, http.StatusInternalServerError, nil)
		return
	}
	if s.Code() == codes.NotFound {
		response(w, http.StatusNotFound, nil)
		return
	}

	response(w, http.StatusOK, worker)
}

func PutWorker(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	var worker gw.PutWorkerRequest
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := client.PutWorker(context.Background(), &worker)
	if err != nil {
		log.Fatalln(err.Error())
	}

	response(w, http.StatusNoContent, nil)
}

func closeBody(r *http.Request) {
	if err := r.Body.Close(); err != nil {
		log.Fatalln(err.Error())
	}
}

func response(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func render(w http.ResponseWriter, name string) {
	var tmpl = template.Must(template.ParseFiles("client/public/" + name + ".html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Fatalln(err.Error())
	}
}

func summon(ctx context.Context) {
	// 取出工人
	resp, err := client.GetWorker(ctx, &gw.GetWorkerRequest{})
	s, ok := status.FromError(err)
	if !ok {
		log.Println(err.Error())
		return
	}

	// 等待
	if s.Code() == codes.NotFound {
		time.Sleep(time.Second)
		log.Println("retrying...")
		summon(ctx)
		return
	}

	// 延遲
	time.Sleep(time.Duration(resp.Worker.Delay) * time.Microsecond)
	log.Printf("Number: %d, Delay: %d", resp.Worker.Number, resp.Worker.Delay)

	// 放回工人
	_, err = client.PutWorker(ctx, &gw.PutWorkerRequest{Number: resp.Worker.Number})
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func NewClientConn(ctx context.Context, addr string) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err.Error())
	}

	return conn
}
