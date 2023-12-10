package main

import (
	"encoding/json"
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

/*func Logger(handler http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		return log.Println(r.URL.Path)
		handler.ServeHTTP(w,r)
	})
} */

func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]

	if !ok {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Error record: %v", err)
	}
}

func newTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	err := json.NewDecoder(r.Body).Decode(&task)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	_, ok := tasks[id]

	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("MY TODO"))
		w.WriteHeader(http.StatusOK)

	})

	r.Get("/tasks", getTasks)
	r.Post("/tasks", newTask)
	r.Get("/tasks/{id}", getTaskById)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Printf("Server startup error: %s", err.Error())
		return
	}
}
