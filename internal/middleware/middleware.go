package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/config"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/utility"
)

type Middleware struct {
	cfg *config.Config
}

type UserIDKey struct{}
type UserRoleKey struct{}

func (m *Middleware) CORS(next http.Handler) (handler http.Handler) {
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		for _, origin := range m.cfg.AllowedOrigins {
			headers.Add("Access-Control-Allow-Origin", origin)
		}
		headers.Add("Access-Control-Allow-Methods", m.cfg.AllowedMethods)
		headers.Add("Access-Control-Allow-Headers", m.cfg.AllowedHeaders)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
	return
}

func (m *Middleware) Auth(next http.Handler) (handler http.Handler) {
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if raw == "" || !ok {
			utility.WriteError(w, http.StatusUnauthorized, map[string]string{"message": "login required"})
			return
		}

		claims := new(model.UserClaims)
		token, err := jwt.ParseWithClaims(
			raw,
			claims,
			func(t *jwt.Token) (secret any, err error) {
				secret = m.cfg.JWTSecret
				return
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
			jwt.WithExpirationRequired(),
		)
		if err != nil {
			utility.WriteError(w, http.StatusUnauthorized, map[string]string{"message": err.Error()})
			return
		}
		if !token.Valid {
			utility.WriteError(w, http.StatusUnauthorized, map[string]string{"message": "invalid token"})
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey{}, claims.ID)
		ctx = context.WithValue(ctx, UserRoleKey{}, *claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	return
}

func (m *Middleware) Autho(next http.Handler, requiredRole string) (handler http.Handler) {
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		role, ok := ctx.Value("UserRoleKey").(string)
		if role == "" || !ok {
			utility.WriteError(w, http.StatusUnauthorized, map[string]string{"message": "login required"})
			return
		}
		if role != requiredRole {
			utility.WriteError(w, http.StatusForbidden, map[string]string{"message": "forbidden"})
			return
		}

		next.ServeHTTP(w, r)
	})
	return
}
