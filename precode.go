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

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasksList := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		tasksList = append(tasksList, task)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasksList); err != nil {
		fmt.Println("Ошибка при сериализации:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newTask.ID = fmt.Sprintf("%d", len(tasks)+1)

	tasks[newTask.ID] = newTask

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		fmt.Println("Ошибка при сериализации:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	task, ok := tasks[taskID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		fmt.Println("Ошибка при сериализации:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	_, ok := tasks[taskID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// Регистрация обработчиков
	r.Get("/tasks", GetTasks)
	r.Post("/tasks", CreateTask)
	r.Get("/tasks/{id}", GetTaskByID)
	r.Delete("/tasks/{id}", DeleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
