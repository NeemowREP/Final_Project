package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Final_Project/internal/db"
	"Final_Project/internal/service"
)

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if dstart == "" {
		writeError(w, fmt.Errorf("date not specified"), http.StatusBadRequest)
		return
	}

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			writeError(w, fmt.Errorf("invalid now format"), http.StatusBadRequest)
			return
		}
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	next, err := service.NextDate(now, dstart, repeat)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	 w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(next))
}

func checkDate(task *db.Task) error {
	now := time.Now().Truncate(24 * time.Hour)
	nowStr := now.Format(DateFormat)

	if task.Date == "" || task.Date == "today" {
		task.Date = nowStr
		return nil
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date")
	}

	if task.Repeat != "" {
		_, err := service.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = nowStr
		} else {
			next, err := service.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, fmt.Errorf("json deserialization error"), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, fmt.Errorf("task title not specified"), http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, fmt.Errorf("id not specified"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, fmt.Errorf("task not found"), http.StatusNotFound)
		return
	}

	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		writeError(w, fmt.Errorf("invalid json format"), http.StatusBadRequest)
		return
	}

	if t.ID == "" {
		writeError(w, fmt.Errorf("id not specified"), http.StatusBadRequest)
		return
	}

	if t.Title == "" {
		writeError(w, fmt.Errorf("title not specified"), http.StatusBadRequest)
		return
	}

	if t.Date != "" {
		if _, err := time.Parse(DateFormat, t.Date); err != nil {
			writeError(w, fmt.Errorf("invalid date format"), http.StatusBadRequest)
			return
		}
	}

	if t.Repeat != "" {
		if !isValidRepeat(t.Repeat) {
			writeError(w, fmt.Errorf("invalid repeat format"), http.StatusBadRequest)
			return
		}
	}

	err = db.UpdateTask(&t)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{})
}

func isValidRepeat(s string) bool {
	if len(s) < 3 {
		return false
	}
	if s[0] != 'd' && s[0] != 'w' && s[0] != 'm' {
		return false
	}
	if s[1] != ' ' {
		return false
	}
	for _, c := range s[2:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, fmt.Errorf("id not specified"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, fmt.Errorf("task not found"), http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		now, _ := time.Parse(DateFormat, task.Date)
		next, err := service.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, fmt.Errorf("id not specified"), http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{})
}
