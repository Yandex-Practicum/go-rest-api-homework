package main

import (
	"bytes"
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

// Обработчик для всех задач (метод GET) endpoint /tasks
func allTasks(res http.ResponseWriter, req *http.Request) {

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(resp)
}

// Обработчик для отаправки запросов на сервер (метод Post) endpoint /tasks
func postTasks(res http.ResponseWriter, req *http.Request) {

	var task Task
	var buff bytes.Buffer

	_, err := buff.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buff.Bytes(), &task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	//  Здесь проверяем начличие задачи
	if _, exists := tasks[task.ID]; exists {
		http.Error(res, "Эта задача уже существует", http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задач по ID (метод GET) endpoint /tasks/{id}
func idTasks(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	task, ok := tasks[id]
	if !ok {
		//  Здесь проверяем начличие задачи
		http.Error(res, "Задача не найдена", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		// Здесь возвращаем ошибку в соответствии с заданием
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-type", "application/json")
	res.Write(resp)
}

// Обработчик удаления задач по ID (метод DELETE) endpoint /tasks/{id}
func delTasks(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	_, ok := tasks[id]
	if !ok {
		// Здесь проверяем наличие задачи
		http.Error(res, "Задача не найдена", http.StatusNoContent)
		return
	}

	delete(tasks, id)
	res.Header().Set("Content-type", "application/json")
	// Возвращем статус в соответствии с заданием
	res.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...
	r.Get("/tasks", allTasks)
	r.Post("/tasks", postTasks)
	r.Get("/tasks/{id}", idTasks)
	r.Delete("/tasks/{id}", delTasks)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
