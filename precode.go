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

// CreateTaskHandler обработчик для создания новой задачи
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при декодировании JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Генерация уникального ID (в данном случае, просто увеличиваем счетчик)
	newTaskID := fmt.Sprintf("%d", len(tasks)+1)
	newTask.ID = newTaskID

	// Добавление новой задачи в карту задач
	tasks[newTaskID] = newTask

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

// GetTasksHandler обработчик для получения списка всех задач
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Преобразование карты задач в срез задач для ответа
	allTasks := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		allTasks = append(allTasks, task)
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allTasks)
}

// GetTaskHandler обработчик для получения информации о конкретной задаче
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	task, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteTaskHandler обработчик для удаления задачи
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if _, ok := tasks[taskID]; !ok {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Удаление задачи из карты задач
	delete(tasks, taskID)

	// Отправка ответа
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// Регистрация обработчиков
	r.Post("/tasks", CreateTaskHandler)
	r.Get("/tasks", GetTasksHandler)
	r.Get("/tasks/{id}", GetTaskHandler)
	r.Delete("/tasks/{id}", DeleteTaskHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
