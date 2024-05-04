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
		Description: "Протестировать финальное задание с помощью Postman",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик для получения всех задач
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Обработчик для отправки задачи на сервер
func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask

	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskID := chi.URLParam(r, "id")

	task, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}

	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Обработчик удаления задачи по ID
func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskID := chi.URLParam(r, "id")

	_, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
