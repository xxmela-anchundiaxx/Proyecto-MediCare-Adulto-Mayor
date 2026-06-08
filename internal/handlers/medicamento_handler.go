package handlers

import (
	"encoding/json"
	"net/http"
	//"proyecto-medicare-adulto-mayor/internal/models"
	"proyecto-medicare-adulto-mayor/internal/models"
	"proyecto-medicare-adulto-mayor/internal/storage"
	//"strconv"
	//"github.com/go-chi/chi/v5"
)

type MedicamentoHandler struct {
	Storage storage.Almacen
}

func NewMedicamentoHandler(s storage.Almacen) *MedicamentoHandler {
	return &MedicamentoHandler{Storage: s}
}

func (s *MedicamentoHandler) ListarMedicacion(w http.ResponseWriter, _ *http.Request) {
	medicacion, err := s.Storage.ListarMedicacion()
	if err != nil {
		http.Error(w, "Error al listar medicación", http.StatusInternalServerError)
		return
	}		
	RespondJSON(w, http.StatusOK, medicacion)
}

func (s *MedicamentoHandler) CrearMedicacion(w http.ResponseWriter, r *http.Request) {
    var nueva models.Medicacion
    if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
        http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
        return
    }

    creada, err := s.Storage.CrearMedicacion(nueva)
    if err != nil {
        http.Error(w, "Error al crear medicación", http.StatusInternalServerError)
        return
    }

    RespondJSON(w, http.StatusCreated, creada)
}



