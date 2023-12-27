package handlers

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/mustthink/go-storage-like-redis/config"
)

func BaseAuth(Handler http.HandlerFunc, auth config.BaseAuthConfig) http.HandlerFunc {
	if auth.User == "" && auth.Pass == "" {
		return Handler
	}

	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Basic ") {
			http.Error(w, "Invalid authorization method", http.StatusUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authHeader[6:])
		if err != nil {
			http.Error(w, "Invalid base64 encoding", http.StatusUnauthorized)
			return
		}

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			http.Error(w, "Invalid authorization credentials", http.StatusUnauthorized)
			return
		}

		if creds[0] != auth.User || creds[1] != auth.Pass {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		Handler(w, r)
	}
}
