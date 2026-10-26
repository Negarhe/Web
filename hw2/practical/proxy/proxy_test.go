package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func testConfig(upstream string) config {
	return config{
		upstream: upstream,
		secret:   []byte("test-secret-value"),
		username: "student",
		password: "student",
		tokenTTL: time.Hour,
	}
}

func TestLoginAndProxy(t *testing.T) {
	var gotUser, gotAuth, gotBody, gotQuery string
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser = r.Header.Get(userHeader)
		gotAuth = r.Header.Get("Authorization")
		gotQuery = r.URL.RawQuery
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer backend.Close()

	handler := mustHandler(t, testConfig(backend.URL))

	login := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"student","password":"student"}`))
	loginRec := httptest.NewRecorder()
	handler.ServeHTTP(loginRec, login)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login status = %d, body = %s", loginRec.Code, loginRec.Body.String())
	}
	var issued struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(loginRec.Body.Bytes(), &issued); err != nil || issued.Token == "" {
		t.Fatalf("login response: %s", loginRec.Body.String())
	}

	req := httptest.NewRequest(http.MethodPost, "/anything?x=1", strings.NewReader("hello"))
	req.Header.Set("Authorization", "Bearer "+issued.Token)
	req.Header.Set(userHeader, "forged")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("proxy status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if gotUser != "student" || gotAuth == "" || gotQuery != "x=1" || gotBody != "hello" {
		t.Fatalf("upstream saw user=%q auth=%q query=%q body=%q", gotUser, gotAuth, gotQuery, gotBody)
	}
}

func TestRejectsMissingAndInvalidTokens(t *testing.T) {
	called := false
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	cfg := testConfig(backend.URL)
	handler := mustHandler(t, cfg)

	noToken := httptest.NewRequest(http.MethodGet, "/get", nil)
	noTokenRec := httptest.NewRecorder()
	handler.ServeHTTP(noTokenRec, noToken)
	if noTokenRec.Code != http.StatusUnauthorized || called {
		t.Fatalf("missing token status=%d called=%v", noTokenRec.Code, called)
	}

	bad := httptest.NewRequest(http.MethodGet, "/get", nil)
	bad.Header.Set("Authorization", "Bearer not-a-jwt")
	badRec := httptest.NewRecorder()
	handler.ServeHTTP(badRec, bad)
	if badRec.Code != http.StatusUnauthorized || called {
		t.Fatalf("bad token status=%d called=%v", badRec.Code, called)
	}

	other := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "student",
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	signed, err := other.SignedString([]byte("another-secret-value"))
	if err != nil {
		t.Fatal(err)
	}
	wrongKey := httptest.NewRequest(http.MethodGet, "/get", nil)
	wrongKey.Header.Set("Authorization", "Bearer "+signed)
	wrongRec := httptest.NewRecorder()
	handler.ServeHTTP(wrongRec, wrongKey)
	if wrongRec.Code != http.StatusUnauthorized || called {
		t.Fatalf("wrong key status=%d called=%v", wrongRec.Code, called)
	}
}

func TestRejectsExpiredToken(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("upstream should not be called")
	}))
	defer backend.Close()

	cfg := testConfig(backend.URL)
	handler := mustHandler(t, cfg)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "student",
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
	})
	signed, err := token.SignedString(cfg.secret)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestLoginIsNeverProxied(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("upstream should not be called")
	}))
	defer backend.Close()

	cfg := testConfig(backend.URL)
	handler := mustHandler(t, cfg)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "student",
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	signed, err := token.SignedString(cfg.secret)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestLoginFailureIsNotProxied(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("upstream should not be called")
	}))
	defer backend.Close()

	handler := mustHandler(t, testConfig(backend.URL))
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"student","password":"wrong"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func mustHandler(t *testing.T, cfg config) http.Handler {
	t.Helper()
	handler, err := newHandler(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}
