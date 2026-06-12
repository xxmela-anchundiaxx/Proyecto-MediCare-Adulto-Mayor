package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"proyecto-medicare-adulto-mayor/internal/models"
	"proyecto-medicare-adulto-mayor/internal/storage"
)

func CreateMedicamento(w http.ResponseWriter, r *http.Request) {

	var medicamento models.Medicamento_farmacia

	err := json.NewDecoder(r.Body).Decode(&medicamento)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	medicamento.ID = storage.IDCounter
	storage.IDCounter++

	storage.Medicamentos = append(storage.Medicamentos, medicamento)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(medicamento)
}

func GetMedicamentos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(storage.Medicamentos)
}
func GetMedicamento(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	for _, medicamento := range storage.Medicamentos {

		if medicamento.ID == id {
			json.NewEncoder(w).Encode(medicamento)
			return
		}
	}

	http.Error(w, "Medicamento no encontrado", http.StatusNotFound)
}

func GetMedicamentosDisponibles(w http.ResponseWriter, r *http.Request) {

	var disponibles []models.Medicamento_farmacia

	for _, medicamento := range storage.Medicamentos {

		if medicamento.Disponible {
			disponibles = append(disponibles, medicamento)
		}
	}

	json.NewEncoder(w).Encode(disponibles)
}

func UpdateMedicamento(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var medicamentoActualizado models.Medicamento_farmacia

	err = json.NewDecoder(r.Body).Decode(&medicamentoActualizado)

	if err != nil {
		http.Error(w, "Datos invalidos", http.StatusBadRequest)
		return
	}

	for i, medicamento := range storage.Medicamentos {

		if medicamento.ID == id {

			medicamentoActualizado.ID = id

			storage.Medicamentos[i] = medicamentoActualizado

			json.NewEncoder(w).Encode(medicamentoActualizado)
			return
		}
	}

	http.Error(w, "Medicamento no encontrado", http.StatusNotFound)
}
func DeleteMedicamento(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	for i, medicamento := range storage.Medicamentos {

		if medicamento.ID == id {

			storage.Medicamentos = append(
				storage.Medicamentos[:i],
				storage.Medicamentos[i+1:]...,
			)

			w.Write([]byte("Medicamento eliminado"))
			return
		}
	}

	http.Error(w, "Medicamento no encontrado", http.StatusNotFound)
}

func GetMedicamentoMasBarato(w http.ResponseWriter, r *http.Request) {

	if len(storage.Medicamentos) == 0 {
		http.Error(w, "No hay medicamentos", http.StatusNotFound)
		return
	}

	masBarato := storage.Medicamentos[0]

	for _, medicamento := range storage.Medicamentos {

		if medicamento.Precio < masBarato.Precio {
			masBarato = medicamento
		}
	}

	json.NewEncoder(w).Encode(masBarato)
}

