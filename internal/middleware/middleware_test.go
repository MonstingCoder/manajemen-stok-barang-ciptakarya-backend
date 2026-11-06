package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/config"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

func TestAuthRejectsTokenWithoutRoleClaim(t *testing.T) {
	secret := "test-secret"
	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.UserClaims{
		UserID: 1,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{JWTSecret: &secret}
	nextCalled := false
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+raw)
	res := httptest.NewRecorder()

	NewMiddlewares(cfg).Auth(next).ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
	}
	if nextCalled {
		t.Fatal("next handler was called for a token without a role claim")
	}
}

func TestAuthAcceptsValidToken(t *testing.T) {
	secret := "test-secret"
	role := "admin"
	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.UserClaims{
		UserID: 1,
		Role:   &role,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{JWTSecret: &secret}
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = r.Context().Value(UserRoleKey{}) == role
		w.WriteHeader(http.StatusNoContent)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+raw)
	res := httptest.NewRecorder()

	NewMiddlewares(cfg).Auth(next).ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNoContent)
	}
	if !nextCalled {
		t.Fatal("next handler was not called with the authenticated role")
	}
}
