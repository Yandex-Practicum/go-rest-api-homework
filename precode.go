package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	errMsgNotFound      = "task not found with id:"
	errMsgIdEmpty       = "cannot recognised empty id"
	errMsgBodyEmpty     = "cannot recognised empty body"
	errMsgAlreadyExists = "entity already exists id:"
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

// Ниже напишите обработчики для каждого эндпоинта
// ...
func getTasksHandler(w http.ResponseWriter, r *http.Request) {

	list := make([]Task, 0, len(tasks))

	for _, e := range tasks {
		list = append(list, e)
	}

	body, err := json.Marshal(list)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body)
}

func getTaskByIdHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		msg := errMsgNotFound + id
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body)
}

func setTaskHandler(w http.ResponseWriter, r *http.Request) {

	if r.ContentLength == 0 {
		http.Error(w, errMsgBodyEmpty, http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	var task Task

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := task.ID

	if id == "" {
		http.Error(w, errMsgIdEmpty, http.StatusBadRequest)
		return
	}

	if _, ok := tasks[id]; ok {
		msg := errMsgAlreadyExists + id
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	tasks[id] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	if !ok {
		msg := errMsgNotFound + id
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...

	r.Get("/tasks", getTasksHandler)
	r.Get("/tasks/{id}", getTaskByIdHandler)
	r.Delete("/tasks/{id}", deleteTaskHandler)
	r.Post("/tasks", setTaskHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
