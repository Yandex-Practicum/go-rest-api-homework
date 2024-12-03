package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func tasksRead(res http.ResponseWriter, req *http.Request) {
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonData)
}

func tasksUpdate(res http.ResponseWriter, req *http.Request) {
	var task Task

	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(res, "Поле ID обязательно", http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
}

func tasksId(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if task, ok := tasks[id]; ok {
		jsonData, err := json.Marshal(task)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(jsonData)
	} else {
		http.Error(res, "Задача ID "+id+" не найдена", http.StatusBadRequest)
	}

}

func taskDelete(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if _, ok := tasks[id]; ok {
		delete(tasks, id)
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
	} else {
		http.Error(res, "Задача ID "+id+" не найдена", http.StatusBadRequest)
	}

}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", tasksRead)
	r.Post("/tasks", tasksUpdate)
	r.Get("/tasks/{id}", tasksId)
	r.Delete("/tasks/{id}", taskDelete)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
