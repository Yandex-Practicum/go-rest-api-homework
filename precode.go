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
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	// Код в комментариях написал на случай, если потребуется
	// возвращать после вызова метода мапу Tasks
	// resp, err := json.Marshal(tasks)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// w.Write(resp)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, exists := tasks[id]; exists {
		resp, err := json.Marshal(tasks[id])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return
	}
	http.Error(w, "Задача отсутствует в мапе", http.StatusBadRequest)
}

func delTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, exists := tasks[id]; exists {
		delete(tasks, id)
		//Код в комментариях написал на случай, если потребуется
		//возвращать после вызова метода мапу Tasks
		// resp, err := json.Marshal(tasks)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//w.Write(resp)
		return
	}
	http.Error(w, "Задача отсутствует в мапе", http.StatusBadRequest)

}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", delTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера", err.Error())
		return
	}
}
