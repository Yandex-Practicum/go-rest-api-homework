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

// Получение всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Ошибка при кодировании задач", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Отправка задачи на сервер
func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Некорректные данные в запросе", http.StatusBadRequest)
		return
	}
	if task.ID == "" {
		http.Error(w, "Отсутствует ID задачи", http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	w.WriteHeader(http.StatusCreated)
}

// Получение задачи по ID
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Ошибка при кодировании задачи", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Удаление задачи по ID
func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	if _, exists := tasks[id]; !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
