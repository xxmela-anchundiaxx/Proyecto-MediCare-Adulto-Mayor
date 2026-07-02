package medicacion

import (
	"encoding/json"
	"medicare-adulto-mayor/internal/handlers/respond"
	"medicare-adulto-mayor/internal/models/medicacion"
	servicioMedicacion "medicare-adulto-mayor/internal/service/medicacion"
	"net/http"
)

type ManejadorMedicamento struct {
	Servicio *servicioMedicacion.ServicioMedicamento
}

func NuevoManejadorMedicamento(s *servicioMedicacion.ServicioMedicamento) *ManejadorMedicamento {
	return &ManejadorMedicamento{Servicio: s}
}

func (h *ManejadorMedicamento) RegistrarMedicamento(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	var req medicacion.CreateMedicamentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.ResponderError(w, http.StatusBadRequest, "cuerpo de petición inválido")
		return
	}

	m, err := h.Servicio.RegistrarMedicamiento(req)
	if err != nil {
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusCreated, m)
}

func (h *ManejadorMedicamento) ListarPorPaciente(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	pacienteID := r.URL.Query().Get("paciente_id")
	if pacienteID == "" {
		respond.ResponderError(w, http.StatusBadRequest, "paciente_id es obligatorio")
		return
	}

	medicamentos, err := h.Servicio.ListarPorPaciente(pacienteID)
	if err != nil {
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusOK, medicamentos)
}
