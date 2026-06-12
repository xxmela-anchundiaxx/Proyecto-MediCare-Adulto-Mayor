package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

    "github.com/go-chi/chi/v5" 

	"proyecto-medicare-adulto-mayor/internal/models"
    "proyecto-medicare-adulto-mayor/internal/storage"
	
)

func CrearRelacionHandler(w http.ResponseWriter, r *http.Request) {
	var relacion models.CuidadorPaciente
	err := json.NewDecoder(r.Body).Decode(&relacion)
	if err != nil {
		http.Error(w, "Datos invalidos", http.StatusBadRequest)
		return
	}

	relacion = storage.CrearRelacion(relacion)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(relacion)
}

func ObtenerRelacionesHandler(w http.ResponseWriter, r *http.Request) {
	relaciones := storage.ObtenerRelaciones()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relaciones)
}

func ObtenerRelacionPorIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
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

	idStr := chi.URLParam(r, "id")
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

func EliminarRelacionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	eliminado := storage.EliminarRelacion(id)
	if !eliminado {
		http.Error(w, "Relacion no encontrada", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "mensaje": "Relación eliminada correctamente",
    })
}
