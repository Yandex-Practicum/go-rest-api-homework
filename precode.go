package main

import (
	"encoding/json"
	"fmt"
	"log"
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
// ...
func getALLTasks(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("Ошибка json.Marshal: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("Ошибка w.Write(response): %v", err)
		return
	}
}
func postTasks(w http.ResponseWriter, r *http.Request) {
	var postTask Task
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&postTask)
	if err != nil {
		log.Printf("Ошибка dec.Decode: %v", err)
		return
	}
	tasks[postTask.ID] = postTask
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
func getTasksId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	findedTask, ok := tasks[id]
	if !ok {
		log.Print("Ошибка task не найден")
		return
	}
	response, err := json.Marshal(findedTask)
	if err != nil {
		log.Printf("Ошибка json.Marshal: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("Ошибка w.Write(response): %v", err)
		return
	}
}
func deleteTasksId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	findedTask, ok := tasks[id]
	if !ok {
		log.Print("Ошибка task не найден")
		return
	}
	delete(tasks, findedTask.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...
	//Обработчик для получения всех задач
	r.Get("/tasks", getALLTasks)
	//Обработчик для отправки задачи на сервер
	r.Post("/tasks", postTasks)
	//Обработчик для получения задачи по ID
	r.Get("/tasks/{id}", getTasksId)
	//Обработчик удаления задачи по ID
	r.Delete("/tasks/{id}", deleteTasksId)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
