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

// Ниже напишите обработчики для каждого эндпоинта
func getTasks(writer http.ResponseWriter, request *http.Request) {
	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

func addTask(writer http.ResponseWriter, request *http.Request) {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var task Task
	if err = json.Unmarshal(buffer.Bytes(), &task); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func getTaskById(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(writer, "Task not found", http.StatusNoContent)
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

func deleteTaskById(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	if _, ok := tasks[id]; ok {
		delete(tasks, id)
		writer.WriteHeader(http.StatusOK)
		return
	}

	http.Error(writer, "Task not found", http.StatusBadRequest)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", addTask)
	r.Get("/tasks/{id}", getTaskById)
	r.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
