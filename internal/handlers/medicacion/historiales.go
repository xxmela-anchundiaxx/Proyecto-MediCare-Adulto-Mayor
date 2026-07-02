package medicacion

import (
	"encoding/json"
	"medicare-adulto-mayor/internal/handlers/respond"
	"medicare-adulto-mayor/internal/models/medicacion"
	servicioMedicacion "medicare-adulto-mayor/internal/service/medicacion"
	"net/http"
)

type ManejadorHistorial struct {
	Servicio *servicioMedicacion.ServicioHistorial
}

func NuevoManejadorHistorial(s *servicioMedicacion.ServicioHistorial) *ManejadorHistorial {
	return &ManejadorHistorial{Servicio: s}
}

func (h *ManejadorHistorial) RegistrarToma(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	var req medicacion.RecordAdherenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.ResponderError(w, http.StatusBadRequest, "cuerpo de petición inválido")
		return
	}

	hist, err := h.Servicio.RegistrarToma(req)
	if err != nil {
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusCreated, hist)
}

func (h *ManejadorHistorial) ListarPorPaciente(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	pacienteID := r.URL.Query().Get("paciente_id")
	if pacienteID == "" {
		respond.ResponderError(w, http.StatusBadRequest, "paciente_id es obligatorio")
		return
	}

	historial, err := h.Servicio.ListarPorPaciente(pacienteID)
	if err != nil {
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusOK, historial)
}
