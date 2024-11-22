package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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

// функция выдает в консоль сообщение об ошибке, если метод w.Write не был завершен корректно
func writeErrorLog(b int, err error) {
	if err != nil {
		fmt.Printf("Ошибка при обработке запроса: %v\n", err)
	}
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeErrorLog(w.Write(resp))
}

func postTask(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")

	//проверка на дубликат
	_, idExist := tasks[task.ID]
	if idExist {
		fmt.Printf("таск с ID: %s уже существует\n", task.ID)
		w.WriteHeader(http.StatusNotModified)
	} else {
		tasks[task.ID] = task
		w.WriteHeader(http.StatusCreated)
	}
}

func getOneTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Нет задачи с таким ID", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeErrorLog(w.Write(resp))
}

func deleteOneTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	bufferTask, ok := tasks[id]
	if !ok {
		//StatusnoContent вроде подходит лучше, но строгая формулировка просит вернуть код 400
		http.Error(w, "Нет задачи с таким ID", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Удален таск с ID=%s (%s)\n", id, bufferTask.Description)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getOneTask)
	r.Delete("/tasks/{id}", deleteOneTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
