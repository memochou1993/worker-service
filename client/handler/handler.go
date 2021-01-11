package handler

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

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
	Client gw.ServiceClient
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	Client = gw.NewServiceClient(newClientConn(ctx, target))
}

// Index 渲染首頁
func Index(w http.ResponseWriter, r *http.Request) {
	render(w, "index")
}

// GetWorker 取出工人
func GetWorker(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	resp, err := Client.GetWorker(context.Background(), &gw.GetWorkerRequest{})
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

// PutWorker 放回工人
func PutWorker(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	var req gw.PutWorkerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response(w, http.StatusBadRequest, nil)
		return
	}
	if _, err := Client.PutWorker(context.Background(), &req); err != nil {
		response(w, http.StatusInternalServerError, nil)
		return
	}

	response(w, http.StatusNoContent, nil)
}

// ListWorkers 列出工人
func ListWorkers(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	resp, err := Client.ListWorkers(context.Background(), &gw.ListWorkersRequest{})
	if err != nil {
		response(w, http.StatusInternalServerError, nil)
		return
	}

	response(w, http.StatusOK, resp)
}

// ShowWorker 查看工人
func ShowWorker(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	n, err := strconv.Atoi(mux.Vars(r)["n"])
	if err != nil {
		response(w, http.StatusNotFound, nil)
		return
	}

	resp, err := Client.ShowWorker(context.Background(), &gw.ShowWorkerRequest{Number: float32(n)})
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

// SummonWorkersSync 同步傳喚工人
func SummonWorkersSync(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	times := 100
	for i := 0; i < times; i++ {
		summon(context.Background())
	}

	ListWorkers(w, r)
}

// SummonWorkersAsync 非同步傳喚工人
func SummonWorkersAsync(w http.ResponseWriter, r *http.Request) {
	defer closeBody(r)

	times := 100
	wg := sync.WaitGroup{}
	wg.Add(times)
	for i := 0; i < times; i++ {
		go func() {
			defer wg.Done()
			summon(context.Background())
		}()
	}
	wg.Wait()

	ListWorkers(w, r)
}

// summon 傳喚工人
func summon(ctx context.Context) {
	// 取出工人
	resp, err := Client.GetWorker(ctx, &gw.GetWorkerRequest{})

	// 檢查錯誤
	s, ok := status.FromError(err)
	if !ok {
		log.Println(err.Error())
		return
	}

	// 重試
	if s.Code() != codes.OK {
		time.Sleep(time.Second)
		log.Println("Retrying...")
		summon(ctx)
		return
	}

	// 延遲
	time.Sleep(time.Duration(resp.Worker.Delay) * time.Microsecond)
	log.Printf("Number: %d, Delay: %d", int64(resp.Worker.Number), int64(resp.Worker.Delay))

	// 放回工人
	if _, err = Client.PutWorker(ctx, &gw.PutWorkerRequest{Number: resp.Worker.Number}); err != nil {
		time.Sleep(time.Second)
		log.Println("Retrying...")
		summon(ctx)
		return
	}
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
	var tmpl = template.Must(template.ParseFiles("Client/public/" + name + ".html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Fatalln(err.Error())
	}
}

func newClientConn(ctx context.Context, addr string) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err.Error())
	}

	return conn
}
