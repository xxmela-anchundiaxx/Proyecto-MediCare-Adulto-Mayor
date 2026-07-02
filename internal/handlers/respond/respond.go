package respond

import (
	"encoding/json"
	"net/http"
)

// ResponderJSON serializa un payload en formato JSON y lo envía al cliente con el código HTTP correspondiente
func ResponderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// En caso de error de serialización, enviar error interno
			http.Error(w, `{"error": "error de serialización interna"}`, http.StatusInternalServerError)
		}
	}
}

// ResponderError envía un error formateado en JSON al cliente
func ResponderError(w http.ResponseWriter, status int, mensaje string) {
	ResponderJSON(w, status, map[string]string{"error": mensaje})
}
