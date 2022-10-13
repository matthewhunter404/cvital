package router

import (
	"context"
	"cvital/domain/users"
	"errors"
	"net/http"
)

func (s *Server) authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenString, err := users.GetTokenFromBearerAuth(r.Header.Get("Authorization"))
		if err != nil {
			s.WriteJsonResponse(w, httpResponse{
				Code:   http.StatusBadRequest,
				Error:  "missing or invalid Authorization header, expecting Bearer Auth token",
				Result: nil,
			})
			return
		}

		claims, err := s.UsersUseCase.ValidateToken(tokenString)
		if err != nil {
			s.WriteJsonResponse(w, httpResponse{
				Code:   http.StatusUnauthorized,
				Error:  err.Error(),
				Result: nil,
			})
			return
		}
		SetContextClaims(ctx, claims)
		next.ServeHTTP(w, r)
	})
}

const ClaimsContextKey = "jwt_token_claims"

func SetContextClaims(ctx context.Context, claims *users.Claims) {
	ctx = context.WithValue(ctx, ClaimsContextKey, *claims)
}

// *jwt.Token?
func GetContextClaims(ctx context.Context) (*users.Claims, error) {
	jwtToken, ok := ctx.Value(ClaimsContextKey).(*users.Claims)
	if !ok {
		return nil, errors.New("access token claims missing")
	}
	return jwtToken, nil
}
