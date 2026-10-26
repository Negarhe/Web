package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer     = "jwt-proxy"
	userHeader = "X-Authenticated-User"
)

func newHandler(cfg config) (http.Handler, error) {
	target, err := url.Parse(cfg.upstream)
	if err != nil || target.Scheme == "" || target.Host == "" {
		return nil, errors.New("UPSTREAM_URL must be an absolute URL")
	}

	upstream := httputil.NewSingleHostReverseProxy(target)
	upstream.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("upstream %s %s: %v", r.Method, r.URL.Path, err)
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "bad gateway"})
	}

	proxy := &authProxy{cfg: cfg, upstream: upstream}
	mux := http.NewServeMux()
	mux.HandleFunc("/login", proxy.login)
	mux.Handle("/", proxy.requireJWT(upstream))
	return mux, nil
}

type authProxy struct {
	cfg      config
	upstream http.Handler
}

func (p *authProxy) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16))
	if err := decoder.Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if !sameSecret(body.Username, p.cfg.username) || !sameSecret(body.Password, p.cfg.password) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   p.cfg.username,
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(p.cfg.tokenTTL)),
	})
	signed, err := token.SignedString(p.cfg.secret)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not issue token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token":      signed,
		"token_type": "Bearer",
		"expires_in": int(p.cfg.tokenTTL.Seconds()),
	})
}

func (p *authProxy) requireJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, err := bearerToken(r.Header.Get("Authorization"))
		if err != nil {
			writeUnauthorized(w)
			return
		}
		subject, err := p.subject(raw)
		if err != nil {
			writeUnauthorized(w)
			return
		}

		r.Header.Del(userHeader)
		r.Header.Set(userHeader, subject)
		next.ServeHTTP(w, r)
	})
}

func (p *authProxy) subject(raw string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return p.cfg.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(issuer),
	)
	if err != nil || !token.Valid || claims.Subject == "" {
		return "", errors.New("invalid token")
	}
	return claims.Subject, nil
}

func bearerToken(header string) (string, error) {
	scheme, token, ok := strings.Cut(strings.TrimSpace(header), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
		return "", errors.New("missing bearer token")
	}
	return strings.TrimSpace(token), nil
}

func sameSecret(got, want string) bool {
	gotHash := sha256.Sum256([]byte(got))
	wantHash := sha256.Sum256([]byte(want))
	return subtle.ConstantTimeCompare(gotHash[:], wantHash[:]) == 1
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="jwt-proxy"`)
	writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
