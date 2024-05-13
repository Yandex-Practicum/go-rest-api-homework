package main

import (
	"fmt"
	"net/http"
	"strings"
	"encoding/json"

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

// обработчики для каждого эндпоинта
func getTasks(w http.ResponseWriter, r *http.Request) {
	var taskErray []Task
	for _, task := range tasks {
		taskErray = append(taskErray, task)
	}
	resp, err := json.Marshal(taskErray)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// записываем тип контента
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем данные в тело ответа
	w.Write(resp)
}
func getUnusedIDs() []string {
	const maxID = 5
	unusedIDs := make([]string, 0, maxID)
	for i := 1; i <= maxID; i++ {
		id := fmt.Sprintf("%d", i)
		if _, exists := tasks[id]; !exists {
			unusedIDs = append(unusedIDs, id)
		}
	}
	return unusedIDs
}
func addTasks(w http.ResponseWriter, r *http.Request) {
	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	getUnusedIDs()
	if _, idExist := tasks[task.ID]; idExist {
		alreadyExId := make([]string, 0, len(tasks))
		for id := range tasks {
			alreadyExId = append(alreadyExId, id)
		}
		errMsg := fmt.Sprintf("Задача с id %s уже есть.", task.ID)
		unusedIDs := getUnusedIDs()
		if len(unusedIDs) > 0 {
			errMsg += " Но можно заюзать эти: " + strings.Join(unusedIDs, ", ")
		}
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	if len(task.Applications) == 0 {
		task.Applications = append(task.Applications, r.UserAgent())
	}
	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}
func getTasksID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Такой задачи нет", http.StatusNoContent)
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
func deleteTasksID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Такой задачи нет", http.StatusNotFound)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()
	// обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", addTasks)
	r.Get("/tasks/{id}", getTasksID)
	r.Delete("/tasks/{id}", deleteTasksID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		//http.Handle("/redirect", http.RedirectHandler("https://yandex.ru/video/preview/4828233716525781863", http.StatusMultiStatus))
		fmt.Printf("Неправильно. Попробуйте еще раз: %s", err.Error())
		return
	}
}
