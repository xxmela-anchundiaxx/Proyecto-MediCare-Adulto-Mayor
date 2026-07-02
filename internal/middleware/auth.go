package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UsuarioContextKey contextKey = "usuario_id"

// AuthMiddleware intercepta las peticiones y valida que se envíe un ID de usuario válido en el Header Authorization
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "cabecera de autorización faltante"}`, http.StatusUnauthorized)
			return
		}

		partes := strings.Split(authHeader, " ")
		if len(partes) != 2 || strings.ToLower(partes[0]) != "bearer" {
			http.Error(w, `{"error": "formato de autorización inválido. Use 'Bearer <usuario_id>'"}`, http.StatusUnauthorized)
			return
		}

		usuarioID := partes[1]
		if usuarioID == "" {
			http.Error(w, `{"error": "token de usuario inválido"}`, http.StatusUnauthorized)
			return
		}

		// Almacena el usuarioID en el contexto de la petición
		ctx := context.WithValue(r.Context(), UsuarioContextKey, usuarioID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ObtenerUsuarioIDContext recupera el usuario_id del contexto de la petición de forma segura
func ObtenerUsuarioIDContext(ctx context.Context) string {
	val := ctx.Value(UsuarioContextKey)
	if val == nil {
		return ""
	}
	id, ok := val.(string)
	if !ok {
		return ""
	}
	return id
}
