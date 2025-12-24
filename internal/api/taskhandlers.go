package api

import (
	"Final_Project/internal/db"
	"Final_Project/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "неверный формат now", 400)
			return
		}
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	next, err := service.NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, next)
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
		return fmt.Errorf("uncorrect date")
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
		writeError(w, fmt.Errorf("json deserialization error"))
		return
	}

	if task.Title == "" {
		writeError(w, fmt.Errorf("task title not specified"))
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeError(w, err)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}


func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, fmt.Errorf("id not specified"))
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, fmt.Errorf("task not found"))
		return
	}

	writeJSON(w, task)
}


func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		writeError(w, fmt.Errorf("invalid json format"))
		return
	}

	if t.ID == "" {
		writeError(w, fmt.Errorf("id not specified"))
		return
	}

	if t.Title == "" {
		writeError(w, fmt.Errorf("title not specified"))
		return
	}

	if t.Date != "" {
		if _, err := time.Parse("20060102", t.Date); err != nil {
			writeError(w, fmt.Errorf("invalid date format"))
			return
		}
	}

	if t.Repeat != "" {
		if !isValidRepeat(t.Repeat) {
			writeError(w, fmt.Errorf("invalid repeat format"))
			return
		}
	}

	err = db.UpdateTask(&t)
	if err != nil {
		writeError(w, err)
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
		writeError(w, fmt.Errorf("id not specified"))
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, fmt.Errorf("task not found"))
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, err)
			return
		}
	} else {
		now, _ := time.Parse("20060102", task.Date)
		next, err := service.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, err)
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {
			writeError(w, err)
			return
		}
	}

	writeJSON(w, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, fmt.Errorf("id not specified"))
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]string{})
}



