package medicacion

import (
	"encoding/json"
	"errors"
	"net/http"          
    "gorm.io/gorm" 

	"proyecto-medicare-adulto-mayor/internal/models/medicacion"
    "proyecto-medicare-adulto-mayor/internal/response"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ListarPaciente GET api/v1/pacientes
func (h *Server) ListarPacientes(w http.ResponseWriter, r *http.Request) {
    pacientes, err := h.Paciente.Listar()
    if err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al listar pacientes")
        return
    }
    response.RespondJSON(w, http.StatusOK, pacientes)
}

// BuscarPaciente por ID GET api/v1/pacientes/{id}
func (h *Server) BuscarPacientePorID(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    paciente, err := h.Paciente.Obtener(id)
    if err != nil {
        response.RespondError(w, http.StatusNotFound, "Paciente no encontrado")
        return
    }
    response.RespondJSON(w, http.StatusOK, paciente)
}

// CrearPaciente POST api/v1/pacientes
func (h *Server) CrearPaciente(w http.ResponseWriter, r *http.Request) {
    var p medicacion.Paciente
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        response.RespondError(w, http.StatusBadRequest, "JSON inválido")
        return
    }

    nuevo, err := h.Paciente.Crear(p)
    if err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al crear paciente")
        return
    }
    response.RespondJSON(w, http.StatusCreated, nuevo)
}

// ActualzarPaciente PUT api/v1/pacientes/{id}
func (h *Server) ActualizarPaciente(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    var p medicacion.Paciente
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        response.RespondError(w, http.StatusBadRequest, "JSON inválido")
        return
    }

    actualizado, err := h.Paciente.Actualizar(id, p)
    if err != nil {
        // Si el error es que no existe el registro, puedes devolver 404
        if errors.Is(err, gorm.ErrRecordNotFound) {
            response.RespondError(w, http.StatusNotFound, "Paciente no encontrado")
            return
        }
        response.RespondError(w, http.StatusInternalServerError, "Error al actualizar paciente")
        return
    }
    response.RespondJSON(w, http.StatusOK, actualizado)
}


// EliminarPaciente DELETE api/v1/pacientes/{id}
func (h *Server) EliminarPaciente(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    if err := h.Paciente.Eliminar(id); err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al eliminar paciente")
        return
    }
    response.RespondJSON(w, http.StatusNoContent, nil)
}
