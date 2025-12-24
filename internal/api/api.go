package api

import "net/http"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", NextDayHandler)
	mux.HandleFunc("/api/task", TaskHandler)
	mux.HandleFunc("/api/tasks", TasksHandler)
	mux.HandleFunc("/api/task/done", taskDoneHandler)
}
