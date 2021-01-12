package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gobuffalo/packr/v2"

	"github.com/gorilla/mux"
	gw "github.com/memochou1993/worker-service/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	target = ":8600"
)

var (
	client gw.ServiceClient
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client = gw.NewServiceClient(NewClientConn(ctx, target))
}

// Index renders a home page.
func Index(w http.ResponseWriter, r *http.Request) {
	render(w, "index")
}

// GetWorker dequeues a worker.
func GetWorker(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	resp, err := client.GetWorker(context.Background(), &gw.GetWorkerRequest{})
	s, ok := status.FromError(err)
	if !ok {
		response(w, http.StatusInternalServerError, nil)
		return
	}
	if s.Code() == codes.NotFound {
		response(w, http.StatusNotFound, nil)
		return
	}
	if s.Code() != codes.OK {
		response(w, http.StatusInternalServerError, nil)
		return
	}
	log.Printf("Number: %d, Delay: %d", int64(resp.Worker.Number), int64(resp.Worker.Delay))

	response(w, http.StatusOK, resp)
}

// PutWorker enqueues a worker.
func PutWorker(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	var req gw.PutWorkerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response(w, http.StatusBadRequest, nil)
		return
	}
	if _, err := client.PutWorker(context.Background(), &req); err != nil {
		response(w, http.StatusInternalServerError, nil)
		return
	}

	response(w, http.StatusNoContent, nil)
}

// ListWorkers lists workers.
func ListWorkers(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	resp, err := client.ListWorkers(context.Background(), &gw.ListWorkersRequest{})
	if err != nil {
		response(w, http.StatusInternalServerError, nil)
		return
	}

	response(w, http.StatusOK, resp)
}

// ShowWorker shows a worker.
func ShowWorker(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	n, err := strconv.Atoi(mux.Vars(r)["n"])
	if err != nil {
		response(w, http.StatusNotFound, nil)
		return
	}

	resp, err := client.ShowWorker(context.Background(), &gw.ShowWorkerRequest{Number: float32(n)})
	s, ok := status.FromError(err)
	if !ok {
		response(w, http.StatusInternalServerError, nil)
		return
	}
	if s.Code() == codes.NotFound {
		response(w, http.StatusNotFound, nil)
		return
	}
	if s.Code() != codes.OK {
		response(w, http.StatusInternalServerError, nil)
		return
	}

	response(w, http.StatusOK, resp)
}

// SummonWorkers dequeues and enqueues a worker.
func SummonWorkers(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	a, err := strconv.Atoi(mux.Vars(r)["a"])
	if err != nil {
		response(w, http.StatusNotFound, nil)
		return
	}
	s, err := strconv.Atoi(mux.Vars(r)["s"])
	if err != nil {
		response(w, http.StatusNotFound, nil)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(a)
	for i := 0; i < a; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < s; i++ {
				summon(context.Background())
			}
		}()
	}
	wg.Wait()

	ListWorkers(w, r)
}

func summon(ctx context.Context) {
	resp, err := client.GetWorker(ctx, &gw.GetWorkerRequest{})
	if err != nil {
		time.Sleep(time.Second)
		log.Println("Retrying...")
		summon(ctx)
		return
	}

	time.Sleep(time.Duration(resp.Worker.Delay) * time.Microsecond)
	log.Printf("Number: %d, Delay: %d", int64(resp.Worker.Number), int64(resp.Worker.Delay))

	if _, err = client.PutWorker(ctx, &gw.PutWorkerRequest{Number: resp.Worker.Number}); err != nil {
		time.Sleep(time.Second)
		log.Println("Retrying...")
		summon(ctx)
		return
	}
}

func closeBody(r *http.Request) {
	if err := r.Body.Close(); err != nil {
		log.Fatal(err.Error())
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
	box := packr.New("public", "../public")
	html, err := box.FindString(fmt.Sprintf("%s.html", name))
	if err != nil {
		log.Fatal(err.Error())
	}
	tmpl, err := template.New(name).Parse(html)
	if err := tmpl.Execute(w, nil); err != nil {
		log.Fatal(err.Error())
	}
}

// NewClientConn creates a new ClientConn instance.
func NewClientConn(ctx context.Context, addr string) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err.Error())
	}

	return conn
}
