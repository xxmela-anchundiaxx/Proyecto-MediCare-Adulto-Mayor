package middleware

import (
	"context"
	"net/http"
	"strings"

	"proyecto-medicare-adulto-mayor/internal/service"
)

type claveContexto string

const ClaveUsuarioID claveContexto = "userID"

// Modificamos la función para que reciba el AuthService como parámetro
func Autenticacion(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encabezado := r.Header.Get("Authorization")
			partes := strings.SplitN(encabezado, " ", 2)

			if len(partes) != 2 || partes[0] != "Bearer" {
				http.Error(w, `{"error":"token requerido"}`, http.StatusUnauthorized)
				return
			}

			// Ahora llamamos al método VerificarToken desde la instancia authService
			userID, err := authService.VerificarToken(partes[1])
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"token inválido"}`))
				return
			}
			
			ctx := context.WithValue(r.Context(), ClaveUsuarioID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}