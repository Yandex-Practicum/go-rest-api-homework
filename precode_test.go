package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetTasks(t *testing.T) {
	cases := []struct {
		name         string
		tasks        map[string]Task
		wantRespBody string
		wantStatus   int
	}{
		{
			name:         "Empty map",
			tasks:        map[string]Task{},
			wantRespBody: `[]`,
			wantStatus:   http.StatusOK,
		},
		{
			name: "Happy case",
			tasks: map[string]Task{
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
			},
			wantRespBody: `[{"id":"1","description":"Сделать финальное задание темы REST API","note":"Если сегодня сделаю, то завтра будет свободный день. Ура!","applications":["VS Code","Terminal","git"]},{"id":"2","description":"Протестировать финальное задание с помощью Postmen","note":"Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер","applications":["VS Code","Terminal","git","Postman"]}]`,
			wantStatus:   http.StatusOK,
		},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			tasks = tc.tasks

			request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
			response := httptest.NewRecorder()

			getTasksHandler(response, request)

			if response.Result().StatusCode != tc.wantStatus {
				t.Errorf("Wrong status %d, want %d", response.Result().StatusCode, tc.wantStatus)
				return
			}

			got := response.Body.String()
			got = strings.Trim(got, "\n")

			if got != tc.wantRespBody {
				t.Errorf("Wrong answer %s, want %s", got, tc.wantRespBody)
				return
			}
		})
	}
}

func TestGetTaskById(t *testing.T) {
	cases := []struct {
		name         string
		tasks        map[string]Task
		id           string
		wantRespBody string
		wantStatus   int
	}{
		{
			name:         "ID not found",
			tasks:        map[string]Task{"1": {ID: "1"}},
			id:           "2",
			wantRespBody: errMsgNotFound + "2",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name: "Happy case",
			tasks: map[string]Task{
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
			},
			id:           "1",
			wantRespBody: `{"id":"1","description":"Сделать финальное задание темы REST API","note":"Если сегодня сделаю, то завтра будет свободный день. Ура!","applications":["VS Code","Terminal","git"]}`,
			wantStatus:   http.StatusOK,
		},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			tasks = tc.tasks

			request, _ := http.NewRequest(http.MethodGet, "/tasks/", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			response := httptest.NewRecorder()

			getTaskByIdHandler(response, request)

			if response.Result().StatusCode != tc.wantStatus {
				t.Errorf("Wrong status %d, want %d", response.Result().StatusCode, tc.wantStatus)
				return
			}

			got := response.Body.String()
			got = strings.Trim(got, "\n")

			if got != tc.wantRespBody {
				t.Errorf("Wrong answer %s, want %s", got, tc.wantRespBody)
				return
			}
		})
	}

}

func TestSetTask(t *testing.T) {
	cases := []struct {
		name          string
		tasks         map[string]Task
		pl            Task
		wantTasks     map[string]Task
		wantRespBody  string
		wantStatus    int
		emptyBodyTest bool
	}{
		{
			name:          "Empty pl",
			tasks:         map[string]Task{},
			wantTasks:     map[string]Task{},
			wantRespBody:  errMsgBodyEmpty,
			wantStatus:    http.StatusBadRequest,
			emptyBodyTest: true,
		},
		{
			name:         "Empty id",
			tasks:        map[string]Task{},
			pl:           Task{Description: "Without id"},
			wantTasks:    map[string]Task{},
			wantRespBody: errMsgIdEmpty,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name: "Id exists",
			tasks: map[string]Task{
				"1": {ID: "1"},
			},
			pl: Task{ID: "1"},
			wantTasks: map[string]Task{
				"1": {ID: "1"},
			},
			wantRespBody: errMsgAlreadyExists + "1",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name: "Happy case",
			tasks: map[string]Task{
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
			},
			pl: Task{
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
			wantTasks: map[string]Task{
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
			},
			wantRespBody: "",
			wantStatus:   http.StatusOK,
		},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			tasks = tc.tasks

			var body io.Reader

			if !tc.emptyBodyTest {
				pl, err := json.Marshal(tc.pl)
				if err != nil {
					t.Logf("Error on marshaling: %s", err.Error())
					return
				}

				body = bytes.NewReader(pl)
			} else {
				body = nil
			}

			request, _ := http.NewRequest(http.MethodPost, "/tasks", body)
			response := httptest.NewRecorder()

			setTaskHandler(response, request)

			if response.Result().StatusCode != tc.wantStatus {
				t.Errorf("Wrong status %d, want %d", response.Result().StatusCode, tc.wantStatus)
				return
			}

			got := response.Body.String()
			got = strings.Trim(got, "\n")

			if got != tc.wantRespBody {
				t.Errorf("Wrong answer %s, want %s", got, tc.wantRespBody)
				return
			}

			if !reflect.DeepEqual(tasks, tc.wantTasks) {
				t.Errorf("Wrong result %v, want %v", tasks, tc.wantTasks)
				return
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	cases := []struct {
		name         string
		tasks        map[string]Task
		id           string
		wantTasks    map[string]Task
		wantRespBody string
		wantStatus   int
	}{
		{
			name:         "ID not found",
			tasks:        map[string]Task{"1": {ID: "1"}},
			id:           "2",
			wantTasks:    map[string]Task{"1": {ID: "1"}},
			wantRespBody: errMsgNotFound + "2",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name: "Happy case",
			tasks: map[string]Task{
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
			},
			id: "1",
			wantTasks: map[string]Task{
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
			},
			wantRespBody: "",
			wantStatus:   http.StatusOK,
		},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			tasks = tc.tasks

			request, _ := http.NewRequest(http.MethodGet, "/tasks/", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			response := httptest.NewRecorder()

			deleteTaskHandler(response, request)

			if response.Result().StatusCode != tc.wantStatus {
				t.Errorf("Wrong status %d, want %d", response.Result().StatusCode, tc.wantStatus)
				return
			}

			got := response.Body.String()
			got = strings.Trim(got, "\n")

			if got != tc.wantRespBody {
				t.Errorf("Wrong answer %s, want %s", got, tc.wantRespBody)
				return
			}

			if !reflect.DeepEqual(tasks, tc.wantTasks) {
				t.Errorf("Wrong result %v, want %v", tasks, tc.wantTasks)
				return
			}
		})
	}
}
