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

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tasks[newTask.ID] = newTask
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, exists := tasks[id]; !exists {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task deleted successfully"))
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
