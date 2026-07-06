package monitoreo

import (
    "encoding/json"
    "errors"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "gorm.io/gorm"

    models "proyecto-medicare-adulto-mayor/internal/models/monitoreo"
    "proyecto-medicare-adulto-mayor/internal/response"
    monitoreoSvc "proyecto-medicare-adulto-mayor/internal/service/monitoreo"
)

type ManejadorMonitoreo struct {
    Servicio *monitoreoSvc.MonitoreoService
}

func NewManejadorMonitoreo(s *monitoreoSvc.MonitoreoService) *ManejadorMonitoreo {
    return &ManejadorMonitoreo{Servicio: s}
}

func (h *ManejadorMonitoreo) CrearRelacionHandler(w http.ResponseWriter, r *http.Request) {
    var relacion models.CuidadorPaciente
    if err := json.NewDecoder(r.Body).Decode(&relacion); err != nil {
        response.RespondError(w, http.StatusBadRequest, "Datos inválidos")
        return
    }
    if relacion.CuidadorID == 0 {
        response.RespondError(w, http.StatusBadRequest, "cuidador_id es requerido")
        return
    }
    if relacion.PacienteID == 0 {
        response.RespondError(w, http.StatusBadRequest, "paciente_id es requerido")
        return
    }
    if relacion.Relacion == "" {
        response.RespondError(w, http.StatusBadRequest, "relacion es requerida")
        return
    }
    nueva, err := h.Servicio.Crear(relacion)
    if err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al crear relación")
        return
    }
    response.RespondJSON(w, http.StatusCreated, nueva)
}

func (h *ManejadorMonitoreo) ObtenerRelacionesHandler(w http.ResponseWriter, r *http.Request) {
    relaciones, err := h.Servicio.Listar()
    if err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al listar relaciones")
        return
    }
    response.RespondJSON(w, http.StatusOK, relaciones)
}

func (h *ManejadorMonitoreo) ObtenerRelacionPorIDHandler(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }
    relacion, err := h.Servicio.Obtener(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            response.RespondError(w, http.StatusNotFound, "Relación no encontrada")
            return
        }
        response.RespondError(w, http.StatusInternalServerError, "Error al obtener relación")
        return
    }
    response.RespondJSON(w, http.StatusOK, relacion)
}

func (h *ManejadorMonitoreo) ActualizarRelacionHandler(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }
    var datos models.CuidadorPaciente
    if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
        response.RespondError(w, http.StatusBadRequest, "Datos inválidos")
        return
    }
    actualizada, err := h.Servicio.Actualizar(id, datos)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            response.RespondError(w, http.StatusNotFound, "Relación no encontrada")
            return
        }
        response.RespondError(w, http.StatusInternalServerError, "Error al actualizar relación")
        return
    }
    response.RespondJSON(w, http.StatusOK, actualizada)
}

func (h *ManejadorMonitoreo) EliminarRelacionHandler(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }
    ok, err := h.Servicio.Eliminar(id)
    if err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al eliminar relación")
        return
    }
    if !ok {
        response.RespondError(w, http.StatusNotFound, "Relación no encontrada")
        return
    }
    response.RespondJSON(w, http.StatusOK, map[string]string{
        "mensaje": "Relación eliminada correctamente",
    })
}