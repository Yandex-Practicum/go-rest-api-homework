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

func getTasks(w http.ResponseWriter, r *http.Request) {

	jsonData, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Ошибка конвертации данных в JSON: %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Ошибка записи JSON в ответ: %s", err.Error())
		return
	}
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	var buff bytes.Buffer

	if _, err := buff.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Ошибка чтения тела запроса: %s", err.Error())
		return
	}

	err := json.Unmarshal(buff.Bytes(), &newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Ошибка преобразования json в Task: %s", err.Error())
		return
	}

	tasks[newTask.ID] = newTask

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func getTask(w http.ResponseWriter, r *http.Request) {

	taskId := chi.URLParam(r, "id")

	task, found := tasks[taskId]
	if !found {
		errInfo := "Задача не найдена по ID"
		http.Error(w, errInfo, http.StatusBadRequest)
		fmt.Printf(errInfo+": %s.", taskId)
		return
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Ошибка преобразования Task в json: %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Ошибка записи JSON в ответ: %s", err.Error())
	}

}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	taskId := chi.URLParam(r, "id")

	_, found := tasks[taskId]
	if !found {
		errInfo := "Задача не найдена по ID"
		http.Error(w, errInfo, http.StatusBadRequest)
		fmt.Printf(errInfo+": %s.", taskId)
		return
	}

	delete(tasks, taskId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
