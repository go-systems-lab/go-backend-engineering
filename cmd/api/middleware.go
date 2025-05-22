package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("no authorization header provided"))
				return
			}

			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 || authHeaderParts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid authorization header"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(authHeaderParts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid authorization header"))
				return
			}

			username := app.config.auth.basic.user
			password := app.config.auth.basic.password

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
