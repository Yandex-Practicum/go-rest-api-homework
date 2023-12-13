package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

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

func getTasks(write http.ResponseWriter, read *http.Request) {

	response, err := json.Marshal(tasks)

	if err != nil {
		http.Error(write, err.Error(), 500)
		return
	}

	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(200)
	write.Write(response)
}

func postTask(write http.ResponseWriter, read *http.Request) {
	var newTask Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(read.Body)

	if err != nil {
		http.Error(write, err.Error(), 400)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &newTask); err != nil {
		http.Error(write, err.Error(), 400)
		return
	}

	tasks[newTask.ID] = newTask

	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(201)
}

func getTaskById(write http.ResponseWriter, read *http.Request) {
	id := chi.URLParam(read, "id")
	task, ok := tasks[id]

	if !ok {
		http.Error(write, "400 Bad Request", 400)
		return
	}

	response, err := json.Marshal(task)

	if err != nil {
		http.Error(write, err.Error(), 400)
		return
	}

	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(200)
	write.Write(response)
}

func deleteTaskById(write http.ResponseWriter, read *http.Request) {
	id := chi.URLParam(read, "id")
	_, ok := tasks[id]

	if !ok {
		http.Error(write, "400 Bad Request", 400)
		return
	}

	delete(tasks, id)
	write.WriteHeader(200)
}

func main() {
	// здесь регистрируйте ваши обработчики
	read := chi.NewRouter()
	read.Get("/tasks/{id}", getTaskById)
	read.Get("/tasks", getTasks)
	read.Post("/tasks", postTask)
	read.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", read); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
