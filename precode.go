package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`          // ID задачи
	Description  string   `json:"description"` // Заголовок
	Note         string   `json:"note"`        // Описание задачи
	Applications []string `json:"application"` // Приложения, которыми будете пользоваться
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

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", postNewTask)
	r.Get("/task/{id}", getTask)
	r.Delete("/task/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

// Обработчик для получения всех задачS
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// Обработчик для отправки задачи на сервер
func postNewTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	// Десериализация JSON напрямую из тела запроса
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//Добавляем в мапу новую задачу, проверки пройдены
	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	//Проверка на наличие задачи в мапе по ID
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	//Сериализаиця запрошенной задачи
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик удаления задачи по ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	//Проверка на наличие задачи в мапе по ID
	if _, ok := tasks[id]; !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
