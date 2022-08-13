package router

import (
	"cvital/domain/users"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello people of the world!"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})
	r.Get("/user/login", handlerFunction(login))

	return r
}

type httpResponse struct {
	Code   int         `json:"-"`
	Error  string      `json:"error"`
	Result interface{} `json:"result"`
}

type httpHandler func(w http.ResponseWriter, r *http.Request) (*httpResponse, error)

func handlerFunction(h httpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := h(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //TODO stop leaking internal error messages
			return
		}
		responseJson, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //TODO stop leaking internal error messages
			return
		}
		w.Write([]byte(responseJson))
		w.WriteHeader(response.Code)

	}
}

func login(w http.ResponseWriter, r *http.Request) (*httpResponse, error) {
	var loginRequest users.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		return &httpResponse{
			Code:   http.StatusBadRequest,
			Error:  err.Error(),
			Result: nil,
		}, nil
	}

	err = users.Login(loginRequest)
	if err != nil {
		return nil, err
	}

	return &httpResponse{
		Code:   http.StatusAccepted,
		Error:  "",
		Result: nil,
	}, nil
}
