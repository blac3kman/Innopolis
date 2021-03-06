package handler_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blac3kman/Innopolis/internal/demo_app/entities"
	"github.com/blac3kman/Innopolis/internal/demo_app/handler"
	usecase_user "github.com/blac3kman/Innopolis/internal/demo_app/usecase"
	"github.com/blac3kman/Innopolis/internal/demo_app/usecase/mocks"
)

type args struct {
	w *httptest.ResponseRecorder
	r *http.Request
}

func setUpArgs(method string, route string, payload string) args {
	req, err := http.NewRequest(method, route, strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	return args{
		w: httptest.NewRecorder(),
		r: req,
	}
}

func TestNew(t *testing.T) {
	type args struct {
		us usecase_user.User
	}
	tests := []struct {
		name string
		args args
		want handler.Handler
	}{
		{
			name: `Success`,
			args: args{us: &mocks.User{}},
			want: handler.New(&mocks.User{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handler.New(tt.args.us); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_AddUser(t *testing.T) {
	route := `/new/user`

	type fields struct {
		us usecase_user.User
	}
	tests := []struct {
		name           string
		fields         func() fields
		args           args
		wantStatusCode int
		wantBody       string
	}{
		{
			name: `Success`,
			fields: func() fields {

				mock := mocks.User{}
				mock.On(`Create`, context.TODO(), `gopher`, `gopher@kaliningrad.ru`).Return(entities.User{
					ID:    1,
					Name:  "gopher",
					Email: "gopher@kaliningrad.ru",
				}, nil)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"name": "gopher", "email": "gopher@kaliningrad.ru"}`),
			wantStatusCode: http.StatusOK,
			wantBody:       `{"id":1,"name":"gopher","email":"gopher@kaliningrad.ru"}`,
		},
		{
			name: `Bad request - empty email`,
			fields: func() fields {

				mock := mocks.User{}
				mock.On(`Create`, context.TODO(), `gopher`, `gopher@kaliningrad.ru`).Return(entities.User{
					ID:    1,
					Name:  "gopher",
					Email: "gopher@kaliningrad.ru",
				}, nil)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"name": "gopher"}`),
			wantStatusCode: http.StatusBadRequest,
			wantBody:       http.StatusText(http.StatusBadRequest),
		},
		{
			name: `Bad request - empty name`,
			fields: func() fields {

				mock := mocks.User{}
				mock.On(`Create`, context.TODO(), `gopher`, `gopher@kaliningrad.ru`).Return(entities.User{
					ID:    1,
					Name:  "gopher",
					Email: "gopher@kaliningrad.ru",
				}, nil)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"email": "gopher@kaliningrad.ru"}`),
			wantStatusCode: http.StatusBadRequest,
			wantBody:       http.StatusText(http.StatusBadRequest),
		},
		{
			name: `Bad request - empty name`,
			fields: func() fields {

				mock := mocks.User{}
				mock.On(`Create`, context.TODO(), `gopher`, `gopher@kaliningrad.ru`).Return(entities.User{}, errors.New(`some error`))

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"name": "gopher", "email": "gopher@kaliningrad.ru"}`),
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       http.StatusText(http.StatusInternalServerError),
		},
		{
			name: `Bad request - empty payload`,
			fields: func() fields {

				mock := mocks.User{}

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, ``),
			wantStatusCode: http.StatusBadRequest,
			wantBody:       http.StatusText(http.StatusBadRequest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(handler.New(tt.fields().us).AddUser)
			h.ServeHTTP(tt.args.w, tt.args.r)

			gotBody := strings.TrimSpace(tt.args.w.Body.String())

			assert.Equal(t, tt.wantStatusCode, tt.args.w.Code)
			assert.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func Test_handler_EditUser(t *testing.T) {

	route := `user/edit`

	type fields struct {
		us usecase_user.User
	}
	tests := []struct {
		name           string
		fields         func() fields
		args           args
		wantStatusCode int
		wantBody       string
	}{
		{
			name: `Success`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`UpdateEmail`, context.TODO(), int64(1), `newgopher@kaliningrad.ru`).Return(entities.User{
					ID:    1,
					Name:  `gopher`,
					Email: `newgopher@kaliningrad.ru`,
				}, nil)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"user_id": 1, "email":"newgopher@kaliningrad.ru"}`),
			wantStatusCode: http.StatusOK,
			wantBody:       `{"id":1,"name":"gopher","email":"newgopher@kaliningrad.ru"}`,
		},
		{
			name: `Error`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`UpdateEmail`, context.TODO(), int64(1), `newgopher@kaliningrad.ru`).Return(entities.User{}, errors.New(`some error`))

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"user_id": 1, "email":"newgopher@kaliningrad.ru"}`),
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       http.StatusText(http.StatusInternalServerError),
		},
		{
			name: `Bad_payload`,
			fields: func() fields {
				mock := mocks.User{}
				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, ``),
			wantStatusCode: http.StatusBadRequest,
			wantBody:       http.StatusText(http.StatusBadRequest),
		},
		{
			name: `Bad_payload_empty_email`,
			fields: func() fields {
				mock := mocks.User{}

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"user_id": 1}`),
			wantStatusCode: http.StatusBadRequest,
			wantBody:       http.StatusText(http.StatusBadRequest),
		},
		{
			name: `Bad_payload_empty_user_id`,
			fields: func() fields {
				mock := mocks.User{}

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"email": "some@email.ru"}`),
			wantStatusCode: http.StatusBadRequest,
			wantBody:       http.StatusText(http.StatusBadRequest),
		},
		{
			name: `Error not found`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`UpdateEmail`, context.TODO(), int64(99), `newgopher@kaliningrad.ru`).Return(entities.User{}, sql.ErrNoRows)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"user_id": 99, "email":"newgopher@kaliningrad.ru"}`),
			wantStatusCode: http.StatusNotFound,
			wantBody:       http.StatusText(http.StatusNotFound),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(handler.New(tt.fields().us).EditUser)
			h.ServeHTTP(tt.args.w, tt.args.r)

			gotBody := strings.TrimSpace(tt.args.w.Body.String())

			assert.Equal(t, tt.wantStatusCode, tt.args.w.Code)
			assert.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func Test_handler_GetUser(t *testing.T) {
	route := `/user`
	type fields struct {
		us usecase_user.User
	}
	tests := []struct {
		name           string
		fields         func() fields
		args           args
		wantStatusCode int
		wantBody       string
	}{
		{
			name: `Success`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`Get`, context.TODO(), int64(1)).Return(entities.User{
					ID:    1,
					Name:  `gopher`,
					Email: `newgopher@kaliningrad.ru`,
				}, nil)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodGet, route, `{"user_id": 1}'`),
			wantStatusCode: http.StatusOK,
			wantBody:       `{"id":1,"name":"gopher","email":"newgopher@kaliningrad.ru"}`,
		},
		{
			name: `Error_not_found`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`Get`, context.TODO(), int64(999)).Return(entities.User{}, sql.ErrNoRows)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodGet, route, `{"user_id": 999}`),
			wantStatusCode: http.StatusNotFound,
			wantBody:       http.StatusText(http.StatusNotFound),
		},
		{
			name: `Error_bad_payload`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`Get`, context.TODO(), int64(999)).Return(entities.User{}, sql.ErrNoRows)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodGet, route, `{}`),
			wantStatusCode: http.StatusBadRequest,
			wantBody:       http.StatusText(http.StatusBadRequest),
		},
		{
			name: `Error_wrong_payload`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`Get`, context.TODO(), int64(999)).Return(entities.User{}, sql.ErrNoRows)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodGet, route, ``),
			wantStatusCode: http.StatusBadRequest,
			wantBody:       http.StatusText(http.StatusBadRequest),
		},
		{
			name: `Error_500`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`Get`, context.TODO(), int64(999)).Return(entities.User{}, errors.New(`some error`))

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodGet, route, `{"user_id": 999}`),
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       http.StatusText(http.StatusInternalServerError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(handler.New(tt.fields().us).GetUser)
			h.ServeHTTP(tt.args.w, tt.args.r)

			gotBody := strings.TrimSpace(tt.args.w.Body.String())

			assert.Equal(t, tt.wantStatusCode, tt.args.w.Code)
			assert.Equal(t, tt.wantBody, gotBody)
		})
	}
}

func Test_handler_RemoveUser(t *testing.T) {
	route := `user/delete`

	type fields struct {
		us usecase_user.User
	}
	tests := []struct {
		name           string
		fields         func() fields
		args           args
		wantStatusCode int
		wantBody       string
	}{
		{
			name: `Success`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`Delete`, context.TODO(), int64(999)).Return(nil)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"user_id": 999}`),
			wantStatusCode: http.StatusOK,
			wantBody:       http.StatusText(http.StatusOK),
		},
		{
			name: `Error bad request`,
			fields: func() fields {
				mock := mocks.User{}
				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, ``),
			wantStatusCode: http.StatusBadRequest,
			wantBody:       http.StatusText(http.StatusBadRequest),
		},
		{
			name: `Error wrong request`,
			fields: func() fields {
				mock := mocks.User{}
				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"name": "some name"}`),
			wantStatusCode: http.StatusBadRequest,
			wantBody:       http.StatusText(http.StatusBadRequest),
		},
		{
			name: `Error wrong request`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`Delete`, context.TODO(), int64(999)).Return(sql.ErrNoRows)

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"user_id": 999}`),
			wantStatusCode: http.StatusNotFound,
			wantBody:       http.StatusText(http.StatusNotFound),
		},
		{
			name: `Error 500`,
			fields: func() fields {
				mock := mocks.User{}

				mock.On(`Delete`, context.TODO(), int64(999)).Return(errors.New(`some error`))

				return fields{
					us: &mock,
				}
			},
			args:           setUpArgs(http.MethodPost, route, `{"user_id": 999}`),
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       http.StatusText(http.StatusInternalServerError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(handler.New(tt.fields().us).RemoveUser)
			h.ServeHTTP(tt.args.w, tt.args.r)

			gotBody := strings.TrimSpace(tt.args.w.Body.String())

			assert.Equal(t, tt.wantStatusCode, tt.args.w.Code)
			assert.Equal(t, tt.wantBody, gotBody)
		})
	}
}
