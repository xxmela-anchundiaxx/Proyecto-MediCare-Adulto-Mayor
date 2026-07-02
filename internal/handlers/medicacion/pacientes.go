package medicacion

import (
	"encoding/json"
	"medicare-adulto-mayor/internal/handlers/respond"
	"medicare-adulto-mayor/internal/middleware"
	"medicare-adulto-mayor/internal/models/medicacion"
	servicioMedicacion "medicare-adulto-mayor/internal/service/medicacion"
	"net/http"
)

type ManejadorPaciente struct {
	Servicio *servicioMedicacion.ServicioPaciente
}

func NuevoManejadorPaciente(s *servicioMedicacion.ServicioPaciente) *ManejadorPaciente {
	return &ManejadorPaciente{Servicio: s}
}

func (h *ManejadorPaciente) RegistrarPaciente(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	var req medicacion.CreatePacienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.ResponderError(w, http.StatusBadRequest, "cuerpo de petición inválido")
		return
	}

	p, err := h.Servicio.CrearPaciente(req)
	if err != nil {
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusCreated, p)
}

func (h *ManejadorPaciente) ObtenerPaciente(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	// Obtener id desde query param "?id=..." o desde el contexto (usuario logueado)
	id := r.URL.Query().Get("id")
	if id == "" {
		// Fallback: buscar paciente asociado al usuario logueado en el token
		usuarioID := middleware.ObtenerUsuarioIDContext(r.Context())
		p, err := h.Servicio.Repo.BuscarPorUsuarioID(usuarioID)
		if err != nil {
			respond.ResponderError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if p == nil {
			respond.ResponderError(w, http.StatusNotFound, "paciente asociado no encontrado")
			return
		}
		id = p.ID
	}

	p, err := h.Servicio.ObtenerPaciente(id)
	if err != nil {
		respond.ResponderError(w, http.StatusNotFound, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusOK, p)
}

func (h *ManejadorPaciente) ListarPorCuidador(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	cuidadorID := middleware.ObtenerUsuarioIDContext(r.Context())
	pacientes, err := h.Servicio.ListarPorCuidador(cuidadorID)
	if err != nil {
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusOK, pacientes)
}
