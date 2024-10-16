// Пакет main реализует аналог приложения TODO-LIST со списком задач.
//
// Доступный функционал:
// - вывод списка всех задач (Эндпоинт "/tasks". Метод GET.);
// - добавление задачи (Эндпоинт "/tasks". Метод POST.);
// - вывод определенной задачи  (Эндпоинт "/tasks/{ID}". Метод GET.);
// - удаление определенной задачи (Эндпоинт "/tasks/{ID}". Метод DELETE.).
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Структура Task (задача).
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

// Мапа tasks (задачи).
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

// getTasks обработчик для получения всех задач.
func getTasks(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// sendTask обработчик для отправки задачи на сервер.
func sendTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// getTask обработчик для получения задачи по ID.
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// deleteTask удаления задачи по ID.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Post("/tasks", sendTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
