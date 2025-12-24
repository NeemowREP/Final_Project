package api

import "net/http"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/signin", SigninHandler)
	mux.HandleFunc("/api/nextdate", auth(NextDayHandler))
	mux.HandleFunc("/api/task", auth(TaskHandler))
	mux.HandleFunc("/api/tasks", auth(TasksHandler))
	mux.HandleFunc("/api/task/done", auth(taskDoneHandler))
}

