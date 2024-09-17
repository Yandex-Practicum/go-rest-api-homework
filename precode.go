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
// getTask возвращает все задачи
func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Ошибка", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// postTask создает новую запись
func postTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Ошибка", http.StatusBadRequest)
		return
	}
	tasks[newTask.ID] = newTask
	w.WriteHeader(http.StatusCreated)
}

// getTaskID добавляет задачу по ID
func getTaskID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	task, ok := tasks[id]

	if !ok {
		http.Error(w, "Ошибка", http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Ошибка", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// deleteTaskID удаляет задачу по ID
func deleteTaskID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	if _, ok := tasks[id]; !ok {
		http.Error(w, "Ошибка", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...
	r.Get("/tasks", getTask)              // получение всех задач
	r.Post("/tasks", postTask)            // Создание новых задач
	r.Get("/tasks/{id}", getTaskID)       // Получение задач по ID
	r.Delete("/tasks/{id}", deleteTaskID) // Удаление задачи по ID

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
