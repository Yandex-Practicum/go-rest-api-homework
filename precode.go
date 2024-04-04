package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
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
func getTasks(w http.ResponseWriter, r *http.Request) {
	//сериализуем данные
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//записываем в заголовок тип контента, данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	//так как все успешно, то записываем StatusOk
	w.WriteHeader(http.StatusOK)
	//записываем данные ответа
	w.Write(resp)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	//получаем значение параметра из URL запроса и присваиваем переменной
	id := chi.URLParam(r, "id")

	//проверяем существование такой задачи в мапе
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	//сериализуем данные из мапы tasks
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	//так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	//записываем данные ответа
	w.Write(resp)
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	//парсим тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// десериализуем данные в структуру Task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	//добавляем в мапу tasks
	tasks[task.ID] = task
	//в заголовок записываем тип контента, у нас данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	//так как все успешно добавлено записываем статус Created
	w.WriteHeader(http.StatusCreated)
}

func delTask(w http.ResponseWriter, r *http.Request) {
	//получаем значение параметра из URL запроса и присваиваем переменной
	id := chi.URLParam(r, "id")

	//проверяем существование такой задачи в мапе
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	//удаляем задачу из мапы по id
	delete(tasks, id)

	//сериализуем данные из мапы tasks
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	//в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	//так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	//записываем данные ответа
	w.Write(resp)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	//регистрируем в роутере эндпоинт `/tasks`с методом GET, для которого используется обработчик `getTasks`
	r.Get("/tasks", getTasks)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTask`
	r.Get("/tasks/{id}", getTask)
	//регистрируем в роутере эндпоинт `/tasks`с методом POST, для которого используется обработчик `postTask`
	r.Post("/tasks", postTask)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `delTask`
	r.Delete("/tasks/{id}", delTask)
	//запускаем сервер
	if err := http.ListenAndServe(":8081", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
