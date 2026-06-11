package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/xxmela-anchundiaxx/Proyecto-MediCare-Adulto-Mayor/internal/models"
	"github.com/xxmela-anchundiaxx/Proyecto-MediCare-Adulto-Mayor/internal/storage"
)

func CrearRelacionHandler(w http.ResponseWriter, r *http.Request) {

	var relacion models.CuidadorPaciente

	err := json.NewDecoder(r.Body).Decode(&relacion)

	if err != nil {
		http.Error(w, "Datos invalidos", http.StatusBadRequest)
		return
	}

	storage.CrearRelacion(relacion)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(relacion)
}
