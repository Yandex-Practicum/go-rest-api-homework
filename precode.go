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

// Ниже напишите обработчики для каждого эндпоинта
func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	newTask.ID = fmt.Sprint(len(tasks) + 1)

	tasks[newTask.ID] = newTask

	w.WriteHeader(http.StatusCreated)
}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	task, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func deleteTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	_, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getAllTasksHandler)
	r.Post("/tasks", createTaskHandler)
	r.Get("/tasks/{id}", getTaskByIDHandler)
	r.Delete("/tasks/{id}", deleteTaskByIDHandler)

	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
