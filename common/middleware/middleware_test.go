package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luyb177/meow-nook/common/middleware"
)

const testSecret = "test-secret-key"

func makeToken(userID int64, username, role string, exp time.Duration) string {
	claims := middleware.Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(testSecret))
	return s
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	mw := middleware.JWTMiddleware(testSecret)

	var calledWith *middleware.Claims
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledWith = middleware.ClaimsFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	token := makeToken(42, "alice", "volunteer", time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rw := httptest.NewRecorder()

	mw(next).ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	if calledWith == nil {
		t.Fatal("claims were not injected into context")
	}
	if calledWith.UserID != 42 {
		t.Fatalf("expected UserID 42, got %d", calledWith.UserID)
	}
	if calledWith.Role != "volunteer" {
		t.Fatalf("expected role 'volunteer', got %q", calledWith.Role)
	}
}

func TestJWTMiddleware_MissingToken(t *testing.T) {
	mw := middleware.JWTMiddleware(testSecret)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/info", nil)
	rw := httptest.NewRecorder()

	mw(next).ServeHTTP(rw, req)

	if rw.Code == http.StatusOK {
		t.Fatal("expected non-200 response when no token provided")
	}
}

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
	mw := middleware.JWTMiddleware(testSecret)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	token := makeToken(1, "bob", "user", -time.Hour) // already expired

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rw := httptest.NewRecorder()

	mw(next).ServeHTTP(rw, req)

	if rw.Code == http.StatusOK {
		t.Fatal("expected non-200 response for expired token")
	}
}

func TestJWTMiddleware_WrongSecret(t *testing.T) {
	mw := middleware.JWTMiddleware(testSecret)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	token := makeToken(1, "eve", "user", time.Hour)
	// Sign with a different secret above; overwrite part of the token to simulate wrong secret
	badToken := token + "tampered"

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/info", nil)
	req.Header.Set("Authorization", "Bearer "+badToken)
	rw := httptest.NewRecorder()

	mw(next).ServeHTTP(rw, req)

	if rw.Code == http.StatusOK {
		t.Fatal("expected non-200 response for tampered token")
	}
}

func TestClaimsContext(t *testing.T) {
	claims := &middleware.Claims{UserID: 99, Username: "carol", Role: "admin"}
	ctx := middleware.WithClaims(context.Background(), claims)

	got := middleware.ClaimsFromContext(ctx)
	if got == nil {
		t.Fatal("expected claims from context, got nil")
	}
	if got.UserID != 99 {
		t.Fatalf("want UserID 99, got %d", got.UserID)
	}
}

func TestClaimsContext_Empty(t *testing.T) {
	got := middleware.ClaimsFromContext(context.Background())
	if got != nil {
		t.Fatalf("expected nil claims from empty context, got %+v", got)
	}
}
