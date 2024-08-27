package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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

// Обработчик для получения всех задач:
func getTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из мапы tasks
	resp, err := json.Marshal(tasks)

	// статус 500 Internal Server Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// запись в заголовок
	w.Header().Set("Content-Type", "application/json")

	// статус OK
	w.WriteHeader(http.StatusOK)

	// запись сериализованных данных json в тело ответа
	w.Write(resp)
}

// Обработчик для отправки задачи на сервер:
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем данные из мапы tasks
	if err = json.Unmarshal(buf.Bytes(), &tasks); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	// запись в заголовок
	w.Header().Set("Content-Type", "application/json")

	// статус 201 Created
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID:
func getTask(w http.ResponseWriter, r *http.Request) {
	// id из url праметра
	id := chi.URLParam(r, "id")

	// проверка на наличие задачи в мапе tasks
	task, ok := tasks[id]

	// статус 400 Bad Request
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	// сериализуем данные из мапы tasks
	resp, err := json.Marshal(task)

	// статус 400 Bad Request
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// запись в заголовок
	w.Header().Set("Content-Type", "application/json")

	// статус OK
	w.WriteHeader(http.StatusOK)

	// запись сериализованных данных json в тело ответа
	w.Write(resp)
}

// Обработчик удаления задачи по ID:
func deleteTask(w http.ResponseWriter, r *http.Request) {
	// id из url праметра
	id := chi.URLParam(r, "id")

	// проверка на наличие задачи в мапе tasks
	_, ok := tasks[id]

	// статус 400 Bad Request
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	//удаление задачи по id
	delete(tasks, id)

	// статус OK
	w.WriteHeader(http.StatusOK)
}

func main() {
	// роутер
	r := chi.NewRouter()

	// Обработчик для получения всех задач
	r.Get("/tasks", getTasks)

	// Обработчик для отправки задачи на сервер
	r.Post("/tasks", postTask)

	// Обработчик для получения задачи по id
	r.Get("/tasks/{id}", getTask)

	// Обработчик удаления задачи по id
	r.Delete("/tasks/{id}", deleteTask)

	// сервер
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
