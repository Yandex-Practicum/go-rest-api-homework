package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type Task struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Note        string   `json:"note"`
	Application []string `json:"application"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Application: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postman",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Application: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

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

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTask.ID = generateUniqueID()
	tasks[newTask.ID] = newTask

	w.WriteHeader(http.StatusCreated)
}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	task, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	_, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)

	w.WriteHeader(http.StatusOK)
}

func generateUniqueID() string {
	// id задач был выбран просто по порядку, уникальные айди каким-либо образом не вычислял
	return fmt.Sprintf("%d", len(tasks)+1)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
