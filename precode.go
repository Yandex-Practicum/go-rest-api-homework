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

// Обработчик выдающий все задачи
func getTask(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //Ошибка 400 bad request
		return                                            //завершаем в случае ошибки
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //Статус 200 ОК
	w.Write(resp)
}

// Обработчик выдающий задачу по ID
func getTaskID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) //400 bad request
		return                                                                       //завершаем в случае ошибки
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //400 bad request
		return                                            //завершаем в случае ошибки
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //Status OK
	w.Write(resp)
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //Ошибка 400 bad request
		return                                            //завершаем в случае ошибки
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return //завершаем в случае ошибки
	}
	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //Статус ОК (как требуется в задании, так то я бы написал http.StatusCreated)
}

// обработчик удаляющий задачу по id
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, ok := tasks[id]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) //400 bad request
		return                                                                       //завершаем
	}
	delete(tasks, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //Статус 200 ОК
}

func main() {
	r := chi.NewRouter()                //роутер
	r.Get("/tasks", getTask)            //Подключаем получение всех задач
	r.Post("/tasks", postTask)          //Подключаем обработчик отправки задачи на сервер
	r.Get("/tasks/{id}", getTaskID)     //Подключаем обработчик получения задачи по ID
	r.Delete("/tasks/{id}", deleteTask) //Подключаем обработчик удаления задачи по ID

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
