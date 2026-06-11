package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func ObtenerRelacionesHandler(w http.ResponseWriter, r *http.Request) {

	relaciones := storage.ObtenerRelaciones()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(relaciones)
}

func ObtenerRelacionPorIDHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	relacion, encontrado := storage.ObtenerRelacionPorID(id)

	if !encontrado {
		http.Error(w, "Relacion no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(relacion)
}

func ActualizarRelacionHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var nuevaRelacion models.CuidadorPaciente

	err = json.NewDecoder(r.Body).Decode(&nuevaRelacion)

	if err != nil {
		http.Error(w, "Datos invalidos", http.StatusBadRequest)
		return
	}

	actualizado := storage.ActualizarRelacion(id, nuevaRelacion)

	if !actualizado {
		http.Error(w, "Relacion no encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(nuevaRelacion)
}
