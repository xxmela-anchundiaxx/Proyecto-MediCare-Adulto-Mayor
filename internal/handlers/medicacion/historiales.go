package medicacion

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"proyecto-medicare-adulto-mayor/internal/response"
	models "proyecto-medicare-adulto-mayor/internal/models/medicacion"
)

// ListarHistorial GET api/v1/historial
func (h *Server) ListarHistorial(w http.ResponseWriter, r *http.Request) {
	historiales, err := h.Historial.Listar()
	if err != nil {
		response.RespondError(w, http.StatusInternalServerError, "Error al listar historiales")
		return
	}
	response.RespondJSON(w, http.StatusOK, historiales)
}

// BuscarPorID GET api/v1/historial/{id}
func (h *Server) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	registro, err := h.Historial.Obtener(id)
	if err != nil {
		response.RespondError(w, http.StatusNotFound, "Historial no encontrado")
		return
	}
	response.RespondJSON(w, http.StatusOK, registro)
}

// CrearHistorial POST api/v1/historial
func (h *Server) CrearHistorial(w http.ResponseWriter, r *http.Request) {
	var hMed models.HistorialMedicacion
	if err := json.NewDecoder(r.Body).Decode(&hMed); err != nil {
		response.RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	nuevo, err := h.Historial.Crear(hMed)
	if err != nil {
		response.RespondError(w, http.StatusInternalServerError, "Error al crear historial")
		return
	}

	response.RespondJSON(w, http.StatusCreated, nuevo)
}

// ActualizarHistorial PUT api/historial/{id}
func (h *Server) ActualizarHistorial(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var datos models.HistorialMedicacion
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		response.RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	actualizado, err := h.Historial.Actualizar(id, datos)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.RespondError(w, http.StatusNotFound, "Historial no encontrado")
			return
		}
		response.RespondError(w, http.StatusInternalServerError, "Error al actualizar historial")
		return
	}

	response.RespondJSON(w, http.StatusOK, actualizado)
}

// EliminarHistorial DELETE api/v1/historial/{id}
func (h *Server) EliminarHistorial(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    if err := h.Historial.Eliminar(id); err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al eliminar historial")
        return
    }
    response.RespondJSON(w, http.StatusNoContent, nil)
}

// ListarMedicacionPorPaciente GET /api/v1/pacientes/{id}/medicaciones
func (h *Server) ListarMedicacionPorPaciente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	pacienteID, _ := strconv.Atoi(idStr)

	medicaciones, err := h.MedicacionHistorial.ListarMedicacionPorPaciente(pacienteID)
	if err != nil {
		response.RespondError(w, http.StatusInternalServerError, "Error al listar medicaciones")
		return
	}
	response.RespondJSON(w, http.StatusOK, medicaciones)
}

// ListarHistorialPorPaciente GET /api/v1/pacientes/{id}/historial
func (h *Server) ListarHistorialPorPaciente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	pacienteID, _ := strconv.Atoi(idStr)

	historiales, err := h.MedicacionHistorial.ListarHistorialPorPaciente(pacienteID)
	if err != nil {
		response.RespondError(w, http.StatusInternalServerError, "Error al listar historial")
		return
	}
	response.RespondJSON(w, http.StatusOK, historiales)
}