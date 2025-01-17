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

func postTasks(w http.ResponseWriter, r *http.Request) {
var task Task
decoder := json.NewDecoder(r.Body)
    defer r.Body.Close()
 if err := decoder.Decode(&task); err != nil {

        http.Error(w, "Некорректный JSON: "+err.Error(), http.StatusBadRequest)
        return

    }

    tasks[task.ID] = task

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
}

func getPerIdTasks(w http.ResponseWriter, r *http.Request) {

id := chi.URLParam(r, "id")

task, ok := tasks[id]

if !ok {

	http.Error(w, "Задание не найдено", http.StatusNoContent)
        
	return

}

resp, err := json.Marshal(task)

if err != nil {

http.Error(w, err.Error(), http.StatusInternalServerError)

}

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write(resp)

}

func deletePerIdTasks(w http.ResponseWriter, r *http.Request) {

id := chi.URLParam(r, "id")

	 if r.Method != http.MethodDelete {

        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return

    }

		 if _, exists := tasks[id]; !exists {

        http.Error(w, "Задание не найдено", http.StatusNotFound)
        return
				
    }

    delete(tasks, id)


    w.WriteHeader(http.StatusNoContent)
}






func main() {

	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTasks)
	r.Get("/tasks/", getPerIdTasks)
	r.Delete("/tasks/", deletePerIdTasks)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
