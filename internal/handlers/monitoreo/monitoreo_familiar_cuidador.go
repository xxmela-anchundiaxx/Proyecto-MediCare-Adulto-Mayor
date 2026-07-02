package monitoreo

import (
	"encoding/json"
	"medicare-adulto-mayor/internal/handlers/respond"
	"medicare-adulto-mayor/internal/models/monitoreo"
	servicioMonitoreo "medicare-adulto-mayor/internal/service/monitoreo"
	"net/http"
)

type ManejadorMonitoreo struct {
	Servicio *servicioMonitoreo.ServicioMonitoreo
}

func NuevoManejadorMonitoreo(s *servicioMonitoreo.ServicioMonitoreo) *ManejadorMonitoreo {
	return &ManejadorMonitoreo{Servicio: s}
}

func (h *ManejadorMonitoreo) RegistrarSignos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	var req monitoreo.RegistrarSignosRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.ResponderError(w, http.StatusBadRequest, "cuerpo de petición inválido")
		return
	}

	m, err := h.Servicio.RegistrarMedicion(req)
	if err != nil {
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusCreated, m)
}

func (h *ManejadorMonitoreo) ListarPorPaciente(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	pacienteID := r.URL.Query().Get("paciente_id")
	if pacienteID == "" {
		respond.ResponderError(w, http.StatusBadRequest, "paciente_id es obligatorio")
		return
	}

	lista, err := h.Servicio.ListarPorPaciente(pacienteID)
	if err != nil {
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusOK, lista)
}
