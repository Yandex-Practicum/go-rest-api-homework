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

func getAllTasks(w http.ResponseWriter, _ *http.Request) {
	response, err := json.Marshal(tasks)

	if err != nil {
		log.Printf("Error json.Marshal: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(response)

	if err != nil {
		log.Printf("Error w.Write response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&newTask)
	if err != nil {
		log.Printf("Error dec.Decode: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	findedTask, ok := tasks[id]
	if !ok {
		log.Print("Error task not found")
		http.Error(w, "getTaskById: task not found", http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(findedTask)
	if err != nil {
		log.Printf("Error json.Marshal: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(response)
	if err != nil {
		log.Printf("Error w.Write response: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	findedTask, ok := tasks[id]

	if !ok {
		log.Print("Error task not found")
		http.Error(w, "deleteTaskById: task not found", http.StatusBadRequest)
		return
	}

	delete(tasks, findedTask.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskById)
	r.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}

//GET http://localhost:8080/tasks

//POST http://localhost:8080/tasks Body: {
//	"ID":          "3",
//	"Description": "3Протестировать финальное задание с помощью Postmen",
//	"Note":        "3Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
//	"Applications": ["The Kill","A Beautiful Lie","Attack","Live Like A Dream"]
//}

//GET http://127.0.0.1:8080/tasks/2

//DEL http://localhost:8080/tasks/1
