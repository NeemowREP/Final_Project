package api

import "net/http"

func Init(mux *http.ServeMux, cfg *APIConfig) {
	mux.HandleFunc("/api/signin", SigninHandler(cfg))
	mux.HandleFunc("/api/nextdate", Auth(cfg, NextDayHandler))
	mux.HandleFunc("/api/task", Auth(cfg, TaskHandler))
	mux.HandleFunc("/api/tasks", Auth(cfg, TasksHandler))
	mux.HandleFunc("/api/task/done", Auth(cfg, taskDoneHandler))
}
