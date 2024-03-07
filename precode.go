package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Текст ошибок
const (
	ErrorMessageBadRequest          = "Bad Request"
	ErrorMessageInternalServerError = "Internal Server Error"
	ErrorMessageTaskNotFound        = "Task Not Found"
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

	// Создаем слайс для хранения задач
	var taskList []Task

	// Передаём значения мапы tasks в слайс
	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	// Сортируем задачи по их ID
	sort.Slice(taskList, func(i, j int) bool {
		id1, _ := strconv.Atoi(taskList[i].ID)
		id2, _ := strconv.Atoi(taskList[j].ID)
		return id1 < id2
	})

	// Преобразуем список задач в формат JSON.
	out, err := json.Marshal(taskList)
	if err != nil {
		fmt.Printf("Ошибка при маршализации JSON: %s\n", err.Error())
		http.Error(w, ErrorMessageInternalServerError, http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type как application/json.
	w.Header().Set("Content-Type", "application/json")

	// Устанавливаем статус код HTTP как 200 OK.
	w.WriteHeader(http.StatusOK)

	// Отправляем данные в формате JSON клиенту.
	_, err = w.Write(out)

	if err != nil {
		fmt.Printf("Ошибка при отправке JSON: %s\n", err.Error())
		http.Error(w, ErrorMessageInternalServerError, http.StatusInternalServerError)
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

	// Декодируем JSON из тела запроса непосредственно в структуру Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, ErrorMessageBadRequest, http.StatusBadRequest)
		return
	}
	if task.ID == "" {
		http.Error(w, "Задача с пустым ID", http.StatusBadRequest)
		return
	}
	// Проверяем, существует ли задача с таким же ID
	if _, exists := tasks[task.ID]; exists {
		http.Error(w, "Задача с таким ID уже существует", http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

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
		http.Error(w, ErrorMessageTaskNotFound, http.StatusBadRequest)
		return
	}

	// Преобразуем список задач в формат JSON.
	out, err := json.Marshal(task)
	if err != nil {
		fmt.Printf("Ошибка при маршализации JSON: %s\n", err.Error())
		http.Error(w, ErrorMessageInternalServerError, http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type как application/json.
	w.Header().Set("Content-Type", "application/json")

	// Устанавливаем статус код HTTP как 200 OK.
	w.WriteHeader(http.StatusOK)

	// Отправляем данные в формате JSON клиенту.
	_, err = w.Write(out)
	if err != nil {
		fmt.Printf("Ошибка при отправке JSON: %s\n", err.Error())
		http.Error(w, ErrorMessageInternalServerError, http.StatusInternalServerError)
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
		http.Error(w, ErrorMessageTaskNotFound, http.StatusBadRequest)
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
