package router

import (
	"cvital/db"
	"cvital/domain"
	"cvital/domain/profiles"
	"cvital/domain/users"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
)

type Server struct {
	DB              *db.PostgresDB
	UsersUseCase    users.UseCase
	ProfilesUseCase profiles.UseCase
	Logger          zerolog.Logger
}

func NewRouter(s *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello people of the world!"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})
	r.Post("/user/login", s.handlerFunction(s.login))
	r.Post("/user", s.handlerFunction(s.createUser))
	r.Post("/cvprofile", s.handlerFunction(s.createCVProfile))
	return r
}

type httpResponse struct {
	Code   int         `json:"-"`
	Error  string      `json:"error"`
	Result interface{} `json:"result"`
}

type httpHandler func(r *http.Request) (*httpResponse, error)

func (s *Server) handlerFunction(h httpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := h(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			responseJson, err := json.Marshal(httpResponse{
				Error: ErrInternal,
			})
			if err != nil {
				s.Logger.Error().Err(err).Msg("Error marshalling json response")
				http.Error(w, err.Error(), http.StatusInternalServerError) //TODO stop leaking internal error messages
				return
			}
			w.Write(responseJson)
			return
		}
		responseJson, err := json.Marshal(response)
		if err != nil {
			s.Logger.Error().Err(err).Msg("Error marshalling json response")
			http.Error(w, err.Error(), http.StatusInternalServerError) //TODO stop leaking internal error messages
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Code)
		w.Write(responseJson)

	}
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
			s.Logger.Debug().Err(err).Msg("login failure")
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
