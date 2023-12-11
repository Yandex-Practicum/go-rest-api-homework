package main

import (
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
		Description: "Протестировать финальное задание с помощью Postman",
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

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasksHandler)
	r.Post("/tasks", createTaskHandler)
	r.Get("/tasks/{id}", getTaskByIDHandler)
	r.Delete("/tasks/{id}", deleteTaskByIDHandler)

	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)

	if err := http.ListenAndServe(port, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s\n", err.Error())
		return
	}
}
