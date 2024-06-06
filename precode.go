package main

import (
	"bytes"
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

// handleTask обрабатывает GET запросы на получение всех элементов мапы
func handleTask(res http.ResponseWriter, req *http.Request) {
	dataJson, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(dataJson)
}

// handleAddId обрабатывает POST запросы на добавление элемента мапы по ID
func handleAddId(res http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	ind := chi.URLParam(req, "id")
	if _, ok := tasks[ind]; ok {
		http.Error(res, "Task exist", http.StatusBadRequest)
		return
	}
	tasks[ind] = task
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
}

// handleGetId обрабатывает GET запросы на получение элемента мапы по ID
func handleGetId(res http.ResponseWriter, req *http.Request) {
	ind := chi.URLParam(req, "id")
	task, ok := tasks[ind]
	if !ok {
		http.Error(res, "Task not exist", http.StatusBadRequest)
		return
	}

	dataJson, err := json.Marshal(task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(dataJson)
}

// handleDelId обрабатывает DELETE запросы на удаление элемента мапы по ID
func handleDelId(res http.ResponseWriter, req *http.Request) {
	ind := chi.URLParam(req, "id")
	if _, ok := tasks[ind]; !ok {
		http.Error(res, "Task not exist", http.StatusBadRequest)
		return
	}

	delete(tasks, ind)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", handleTask)
	r.Post("/tasks/{id}", handleAddId)
	r.Get("/tasks/{id}", handleGetId)
	r.Delete("/tasks/{id}", handleDelId)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
