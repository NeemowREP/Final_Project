package api

import (
	"Final_Project/internal/db"
	"net/http"
)

type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")

	var (
		tasks []*db.Task
		err   error
	)

	if search == "" {
		tasks, err = db.Tasks(50)
	} else {
		tasks, err = db.TasksSearch(search, 50)
	}

	if err != nil {
		writeError(w, err)
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}

