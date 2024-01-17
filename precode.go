package main

import (
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

	// сериализуем данные из слайса artists
	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTasksId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	findedTask, ok := tasks[id]
	if !ok {
		http.Error(w, "getTasksId: task not found", http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(findedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func postTasks(w http.ResponseWriter, r *http.Request) {
	var task Task
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func deleteTaskId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	findedTask, ok := tasks[id]

	if !ok {
		http.Error(w, "deleteTaskId: task not found", http.StatusBadRequest)
		return
	}

	delete(tasks, findedTask.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()
	//Обработчик должен вернуть все задачи, которые хранятся в мапе.
	// регистрируем в роутере эндпоинт `/artists` с методом GET, для которого используется обработчик `getArtists`
	r.Get("/tasks", getTasks)
	//Обработчик должен вернуть задачу с указанным в запросе пути ID, если такая есть в мапе.
	r.Get("/tasks/{id}", getTasksId)
	// регистрируем в роутере эндпоинт `/artists` с методом POST, для которого используется обработчик `postArtist`
	//Обработчик должен принимать задачу в теле запроса и сохранять ее в мапе.
	r.Post("/tasks", postTasks)
	//Обработчик должен удалить задачу из мапы по её ID.
	r.Delete("/tasks/{id}", deleteTaskId)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
