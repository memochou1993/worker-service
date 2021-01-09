package worker

// // GetWorker 索取一個工人
// func GetWorker(w http.ResponseWriter, r *http.Request) {
// 	defer closeBody(r)
//
// 	if worker := service.Dequeue(); worker != nil {
// 		response(w, http.StatusOK, worker)
// 		return
// 	}
//
// 	response(w, http.StatusNotFound, nil)
// }
//
// // PutWorker 退還一個工人
// func PutWorker(w http.ResponseWriter, r *http.Request) {
// 	defer closeBody(r)
//
// 	var worker Worker
// 	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
// 		response(w, http.StatusInternalServerError, nil)
// 		return
// 	}
// 	service.Enqueue(NewWorker(worker.Number))
//
// 	response(w, http.StatusNoContent, nil)
// }
//
// // ListWorkers 列出所有工人被傳喚記錄
// func ListWorkers(w http.ResponseWriter, r *http.Request) {
// 	defer closeBody(r)
//
// 	var records []Record
// 	for n, u := range service.Attendance {
// 		records = append(records, Record{n, u})
// 	}
//
// 	response(w, http.StatusOK, records)
// }
//
// // ShowWorker 查看特定工人被傳喚記錄
// func ShowWorker(w http.ResponseWriter, r *http.Request) {
// 	defer closeBody(r)
//
// 	n, err := strconv.Atoi(mux.Vars(r)["n"])
// 	if err != nil {
// 		response(w, http.StatusNotFound, nil)
// 		return
// 	}
// 	if _, ok := service.Attendance[Number(n)]; !ok {
// 		response(w, http.StatusNotFound, nil)
// 		return
// 	}
//
// 	record := Record{Number(n), service.Attendance[Number(n)]}
//
// 	response(w, http.StatusOK, record)
// }
//
// func response(w http.ResponseWriter, code int, data interface{}) {
// 	w.WriteHeader(code)
// 	if data == nil {
// 		return
// 	}
// 	if err := json.NewEncoder(w).Encode(Payload{Data: data}); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
//
// func closeBody(r *http.Request) {
// 	if err := r.Body.Close(); err != nil {
// 		log.Fatalln(err.Error())
// 	}
// }
