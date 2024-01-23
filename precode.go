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

// Обработчик для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	//Серриализируем данные из мапы tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// записываем в заголовок тип контента (json)
	w.Header().Set("Content-Type", "application/json")
	//статус 200 OK
	w.WriteHeader(http.StatusOK)
	//записываем серриализованные данные в тело ответа
	w.Write(resp)
}

// Обработчик для отправки задачи на сревер
func postTasks(w http.ResponseWriter, r *http.Request) {
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

// Обработчик для получения задачи по ID
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задачи не обнаружено", http.StatusNoContent)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик удаления задачи из Мапы
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //получили id
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNoContent)
		return
	}

	// Удаляем из Мапы элемент по ID
	delete(tasks, id)
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики

	//  Регистрируем обработчик для получения всех задач
	r.Get("/tasks", getTasks)

	//  Регистрируем обработчик для отправки задачи на сревер
	r.Post("/tasks", postTasks)

	// Регистрируем обработчик для получения задачи по ID
	r.Get("/task/{id}", getTask)

	// Регистрируем обработчик удаления задачи по ID
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
