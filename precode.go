package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Task основная структура
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

// Мапа для хранения полученных данных с POST
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
// ...

// Функция для нахождения максимального значения в слайсе
func Max(arr []int) int {
	max := arr[0]
	for _, value := range arr {
		if value > max {
			max = value
		}
	}
	return max
}

func getTask(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса artists
	resp, err := json.MarshalIndent(tasks, " ", "   ")
	/*resp, err := json.Marshal(tasks)*/
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

/*
Для тестирования в отладке использовано
curl http://localhost:8080/tasks --include --header "Content-Type: application/json" --request "POST" --data '{"id":"3","description":"end","note":"end uau","applications":["VS Code","Terminal","git","Postman"]}'
*/
func postTask(w http.ResponseWriter, r *http.Request) {
	var taskPost Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &taskPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// обработка ключа мапы, увеличим его на max+1, если повторяется отправляемый ID
	if _, ok := tasks[taskPost.ID]; ok {
		//http.Error(w, "Уже существует", http.StatusNoContent)
		alreadyExit := []int{}
		for id := range tasks {
			strId, _ := strconv.Atoi(id)
			alreadyExit = append(alreadyExit, strId)
		}
		idNew := strconv.Itoa(Max(alreadyExit) + 1)
		taskPost.ID = idNew
		tasks[idNew] = taskPost
	} else {
		tasks[taskPost.ID] = taskPost
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getTaskId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tskid, ok := tasks[id]
	if !ok {
		http.Error(w, "Не найдено", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(tskid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func delTaskId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tskid, ok := tasks[id]
	if !ok {
		http.Error(w, "Не найдено", http.StatusNoContent)
		return
	} else {
		delete(tasks, id)
	}

	resp, err := json.Marshal(tskid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	r := chi.NewRouter()
	// здесь регистрируйте ваши обработчики
	// ...

	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTask`
	r.Get("/tasks", getTask)

	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTask`
	r.Post("/tasks", postTask)

	// регистрируем в роутере эндпоинт `/task/{id}` с методом GET, для которого используется обработчик `getTaskId`
	r.Get("/tasks/{id}", getTaskId)

	// регистрируем в роутере эндпоинт `/task/{id}` с методом DELETE, для которого используется обработчик `delTaskId`
	r.Get("/tasks/{id}", delTaskId)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
