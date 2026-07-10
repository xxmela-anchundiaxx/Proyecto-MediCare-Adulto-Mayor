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
	// Corregido: El tipo exacto en tu paquete service es ServicioMonitoreo
	Servicio *monitoreoSvc.ServicioMonitoreo
}

func NewManejadorMonitoreo(s *monitoreoSvc.ServicioMonitoreo) *ManejadorMonitoreo {
	return &ManejadorMonitoreo{Servicio: s}
}

func (h *ManejadorMonitoreo) CrearRelacionHandler(w http.ResponseWriter, r *http.Request) {
	var relacion models.CuidadorPaciente
	if err := json.NewDecoder(r.Body).Decode(&relacion); err != nil {
		response.RespondError(w, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Corregido: Llamada al método correcto CrearRelacion
	nueva, err := h.Servicio.CrearRelacion(relacion)
	if err != nil {
		response.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	response.RespondJSON(w, http.StatusCreated, nueva)
}

func (h *ManejadorMonitoreo) ObtenerRelacionesHandler(w http.ResponseWriter, r *http.Request) {
	// Corregido: Llamada al método correcto ListarRelaciones
	relaciones, err := h.Servicio.ListarRelaciones()
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

	// Corregido: Llamada al método correcto BuscarRelacionPorID
	relacion, err := h.Servicio.BuscarRelacionPorID(id)
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

	// Corregido: Llamada al método correcto ActualizarRelacion
	actualizada, err := h.Servicio.ActualizarRelacion(id, datos)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.RespondError(w, http.StatusNotFound, "Relación no encontrada")
			return
		}
		response.RespondError(w, http.StatusBadRequest, err.Error())
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

	// Corregido: Llamada al método correcto EliminarRelacion
	ok, err := h.Servicio.EliminarRelacion(id)
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