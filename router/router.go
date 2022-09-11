package router

import (
	"cvital/db"
	"cvital/domain"
	"cvital/domain/profiles"
	"cvital/domain/users"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(s *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello people of the world!"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})
	r.Get("/user/login", handlerFunction(s.login))
	r.Post("/user", handlerFunction(s.createUser))
	r.Post("/cvprofile", handlerFunction(s.createCVProfile))
	return r
}

type httpResponse struct {
	Code   int         `json:"-"`
	Error  string      `json:"error"`
	Result interface{} `json:"result"`
}

type httpHandler func(r *http.Request) (*httpResponse, error)

func handlerFunction(h httpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := h(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			responseJson, err := json.Marshal(httpResponse{
				Error: ErrInternal,
			})
			if err != nil {
				log.Printf("Error marshalling json response: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError) //TODO stop leaking internal error messages
				return
			}
			w.Write(responseJson)
			return
		}
		responseJson, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling json response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError) //TODO stop leaking internal error messages
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Code)
		w.Write(responseJson)

	}
}

type Server struct {
	DB              *db.PostgresDB
	UsersUseCase    users.UseCase
	ProfilesUseCase profiles.UseCase
}

// https://www.rfc-editor.org/rfc/rfc6749#section-5.1
type loginResult struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

func (s *Server) login(r *http.Request) (*httpResponse, error) {
	var loginRequest users.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		return &httpResponse{
			Code:   http.StatusBadRequest,
			Error:  err.Error(),
			Result: nil,
		}, nil
	}

	jwt, expiryTime, err := s.UsersUseCase.Login(r.Context(), loginRequest)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrLoginFailed):
			log.Printf("login failure: %v\n", err)
			return &httpResponse{
				Code:  http.StatusUnauthorized,
				Error: ErrLoginFailed,
			}, nil
		default:
			return nil, err
		}
	}

	return &httpResponse{
		Code: http.StatusAccepted,
		Result: loginResult{
			AccessToken: *jwt,
			TokenType:   "jwt",
			ExpiresIn:   expiryTime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (s *Server) createUser(r *http.Request) (*httpResponse, error) {
	var request users.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return &httpResponse{
			Code:   http.StatusBadRequest,
			Error:  err.Error(),
			Result: nil,
		}, nil
	}

	newUser, err := s.UsersUseCase.CreateUser(r.Context(), request)
	if err != nil {
		switch err {
		case domain.ErrAlreadyExists:
			return &httpResponse{
				Code:  http.StatusBadRequest,
				Error: ErrAlreadyExists,
			}, nil
		default:
			return nil, err
		}
	}

	return &httpResponse{
		Code:   http.StatusAccepted,
		Error:  "",
		Result: newUser,
	}, nil
}

func (s *Server) createCVProfile(r *http.Request) (*httpResponse, error) {
	var request profiles.CreateCVProfileRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return &httpResponse{
			Code:   http.StatusBadRequest,
			Error:  err.Error(),
			Result: nil,
		}, nil
	}

	tokenString, err := users.GetTokenFromBearerAuth(r.Header.Get("Authorization"))
	if err != nil {
		return &httpResponse{
			Code:   http.StatusBadRequest,
			Error:  "missing or invalid Authorization header, expecting Bearer Auth token",
			Result: nil,
		}, nil
	}

	claims, err := s.UsersUseCase.ValidateToken(tokenString)
	if err != nil {
		return &httpResponse{
			Code:   http.StatusUnauthorized,
			Error:  err.Error(),
			Result: nil,
		}, nil
	}

	newUser, err := s.ProfilesUseCase.CreateCVProfile(r.Context(), request, claims.Email)
	if err != nil {
		switch err {
		case domain.ErrAlreadyExists:
			return &httpResponse{
				Code:  http.StatusBadRequest,
				Error: ErrAlreadyExists,
			}, nil
		default:
			return nil, err
		}
	}

	return &httpResponse{
		Code:   http.StatusAccepted,
		Result: newUser,
	}, nil
}
