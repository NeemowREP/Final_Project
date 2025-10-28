package api

import (
	"Final_Project/internal/db"
	"Final_Project/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func checkDate(task *db.Task) error {
	now := time.Now().Truncate(24 * time.Hour)
	nowStr := now.Format(DateFormat)

	if task.Date == "" || task.Date == "today" {
		task.Date = nowStr
		return nil
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректная дата")
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
		json.NewEncoder(w).Encode(map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"id": fmt.Sprintf("%d", id)})
}
