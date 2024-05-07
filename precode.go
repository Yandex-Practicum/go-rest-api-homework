package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5" // Импорт пакета chi для маршрутизации
)

// Структура для задачи
type Task struct {
	ID           string   `json:"id"`           // ID задачи
	Description  string   `json:"description"`  // Описание задачи
	Note         string   `json:"note"`         // Примечание к задаче
	Applications []string `json:"applications"` // Приложения, используемые для выполнения задачи
}

// Мапа для хранения задач
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

// Обработчик для получения всех задач
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		return
	}
}

// Обработчик для создания новой задачи
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
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
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		return
	}
}

// Обработчик для удаления задачи по ID
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

// Главная функция, запускающая сервер и регистрирующая обработчики
func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasksHandler)           // Регистрация обработчика для получения всех задач
	r.Post("/tasks", createTaskHandler)        // Регистрация обработчика для создания новой задачи
	r.Get("/tasks/{id}", getTaskHandler)       // Регистрация обработчика для получения задачи по ID
	r.Delete("/tasks/{id}", deleteTaskHandler) // Регистрация обработчика для удаления задачи по ID

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
