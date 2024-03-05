package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Task представляет собой структуру для описания задачи
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

// tasks содержит мапу задач, где ключами являются их идентификаторы, а значениями - сами задачи (Подобие БД)
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

// getTasks отправляет список задач в формате JSON в ответ на запрос GET
//
// w - объект http.ResponseWriter для отправки ответа клиенту
//
// r - объект http.Request, который содержит информацию о запросе клиента
func getTasks(w http.ResponseWriter, r *http.Request) {

	// Преобразуем список задач в формат JSON.
	out, err := json.Marshal(&tasks)

	if err != nil {
		// Если произошла ошибка при маршализации JSON, отправляем клиенту ответ с кодом статуса 500
		http.Error(w, "Ошибка при маршализации JSON", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type как application/json.
	w.Header().Set("Content-Type", "application/json")

	// Устанавливаем статус код HTTP как 200 OK.
	w.WriteHeader(http.StatusOK)

	// Отправляем данные в формате JSON клиенту.
	_, err = w.Write(out)

	if err != nil {
		// Если произошла ошибка при отправке JSON, отправляем клиенту ответ с кодом статуса 500
		http.Error(w, "Ошибка при отправке JSON", http.StatusInternalServerError)
		return
	}
}

// postTask принимает задачу в формате JSON и записывает её в общий список задач
//
// w - объект http.ResponseWriter для отправки ответа клиенту
//
// r - объект http.Request, который содержит информацию о запросе клиента
func postTask(w http.ResponseWriter, r *http.Request) {
	// Инициализируем переменную для хранения новой задачи.
	var task Task

	// Создаем буфер для чтения тела запроса.
	var buf bytes.Buffer

	// Читаем тело запроса в буфер.
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		// Если произошла ошибка при чтении тела запроса, отправляем клиенту ответ с кодом статуса 400 Bad Request.
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}

	// Распаковываем JSON данные из буфера
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		// Если произошла ошибка при маршализации JSON, отправляем клиенту ответ с кодом статуса 400 Bad Request.
		http.Error(w, "Ошибка при маршализации JSON", http.StatusBadRequest)
		return
	}

	// Создаём новый уникальный идентификатор для задачи.
	count := len(tasks) + 1
	// Добавляем задачу в мапу tasks с использованием уникального идентификатора.
	tasks[strconv.Itoa(count)] = task

	// Отправляем клиенту статус 201 Created
	w.WriteHeader(http.StatusCreated)
}

// getTasks отправляет список задач в формате JSON в ответ на запрос GET
//
// w - объект http.ResponseWriter для отправки ответа клиенту
//
// r - объект http.Request, который содержит информацию о запросе клиента
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	// Преобразуем список задач в формат JSON.
	out, err := json.Marshal(task)

	if err != nil {
		// Если произошла ошибка при маршализации JSON, отправляем клиенту ответ с кодом статуса 500
		http.Error(w, "Ошибка при маршализации JSON", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type как application/json.
	w.Header().Set("Content-Type", "application/json")

	// Устанавливаем статус код HTTP как 200 OK.
	w.WriteHeader(http.StatusOK)

	// Отправляем данные в формате JSON клиенту.
	_, err = w.Write(out)

	if err != nil {
		// Если произошла ошибка при отправке JSON, отправляем клиенту ответ с кодом статуса 500
		http.Error(w, "Ошибка при отправке JSON", http.StatusInternalServerError)
		return
	}
}

// deleteTask удаляет задачу из списка задач по её идентификатору.
//
// w - объект http.ResponseWriter для отправки ответа клиенту.
//
// r - объект http.Request, содержащий информацию о запросе клиента.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	// Извлекаем идентификатор задачи из URL-параметра.
	id := chi.URLParam(r, "id")

	// Проверяем, существует ли задача с указанным идентификатором в списке задач.
	_, ok := tasks[id]
	if !ok {
		// Если задачи с указанным идентификатором нет, возвращаем статус 400 Bad Request с сообщением об ошибке.
		http.Error(w, "Ошибка: Такой задачи нет в списке", http.StatusBadRequest)
		return
	}

	// Удаляем задачу из списка задач.
	delete(tasks, id)

	// Отправляем клиенту статус 200 OK, чтобы указать успешное удаление задачи.
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Создаем новый роутер
	r := chi.NewRouter()

	// Назначаем обработчик getTasks для GET-запроса на путь /tasks.
	r.Get("/tasks", getTasks)
	// Назначаем обработчик getTask для GET-запроса на путь /tasks/{id}.
	r.Get("/tasks/{id}", getTask)

	// Назначаем обработчик postTask для POST-запроса на путь /tasks.
	r.Post("/tasks", postTask)

	// Назначаем обработчик deleteTask для DELETE-запроса на путь /tasks/{id}.
	r.Delete("/tasks/{id}", deleteTask)

	// Запускаем HTTP-сервер на порту :8080 с нашим роутером chi
	if err := http.ListenAndServe(":8080", r); err != nil {
		// Если при запуске сервера возникла ошибка, выводим сообщение об ошибке.
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
