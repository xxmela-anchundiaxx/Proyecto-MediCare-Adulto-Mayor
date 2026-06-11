package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"          
    "gorm.io/gorm" 

	"proyecto-medicare-adulto-mayor/internal/models"
	"proyecto-medicare-adulto-mayor/internal/storage"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MedicamentoHandler struct {
	Storage storage.Almacen
}

func NewMedicamentoHandler(s storage.Almacen) *MedicamentoHandler {
	return &MedicamentoHandler{Storage: s}
}

//ListarMedicacion GET /api/v1/medicaciones
func (s *MedicamentoHandler) ListarMedicacion(w http.ResponseWriter, _ *http.Request) {
	medicacion, err := s.Storage.ListarMedicacion()
	if err != nil {
		http.Error(w, "Error al listar medicación", http.StatusInternalServerError)
		return
	}		
	RespondJSON(w, http.StatusOK, medicacion)
}

//ObtenerMedicacion GET /api/v1/medicaciones/{id}
func (s *MedicamentoHandler) ObtenerMedicacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	medicacion, econtrado := s.Storage.BuscarMedicacionPorID(id)
	if econtrado != nil {
		RespondError(w, http.StatusNotFound, "Medicacion no encontrada")
		return
	}
	RespondJSON(w, http.StatusOK, medicacion)
}

//CrearMedicacion POST /api/v1/medicaciones
func (s *MedicamentoHandler) CrearMedicacion(w http.ResponseWriter, r *http.Request) {
    var nueva models.Medicacion
    
    if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
        RespondError(w, http.StatusBadRequest, "Datos de medicación inválidos: "+err.Error())
        return
    }

    if strings.TrimSpace(nueva.Nombre) == "" {
        RespondError(w, http.StatusBadRequest, "El nombre de la medicación es requerido")
        return
    }
    if strings.TrimSpace(nueva.Dosis) == "" {
        RespondError(w, http.StatusBadRequest, "La dosis de la medicación es requerida")
        return
    }

    if strings.TrimSpace(nueva.Frecuencia) == "" {
        RespondError(w, http.StatusBadRequest, "La frecuencia de la medicación es requerida")
        return
    }

    if strings.TrimSpace(nueva.Hora_programada) == "" {
        RespondError(w, http.StatusBadRequest, "La hora programada de la medicación es requerida")
        return
    }

    if nueva.Inicio_tratamiento.IsZero() {
        RespondError(w, http.StatusBadRequest, "La fecha de inicio del tratamiento es requerida")
        return
    }

    if nueva.Fecha_creacion.IsZero() {
        nueva.Fecha_creacion = time.Now()
    }

    creada, err := s.Storage.CrearMedicacion(nueva)
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al guardar en el almacén: "+err.Error())
        return
    }

    RespondJSON(w, http.StatusCreated, creada)
}

//ActualizarMedicacion PUT /api/v1/medicaciones/{id}
func (s *MedicamentoHandler) ActualizarMedicacion(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    var datos models.Medicacion
    if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
        RespondError(w, http.StatusBadRequest, "Datos de medicación inválidos: "+err.Error())
        return
    }

    if strings.TrimSpace(datos.Nombre) == "" {
        RespondError(w, http.StatusBadRequest, "El nombre de la medicación es requerido")
        return
    }
    if strings.TrimSpace(datos.Dosis) == "" {
        RespondError(w, http.StatusBadRequest, "La dosis de la medicación es requerida")
        return
    }
    if strings.TrimSpace(datos.Frecuencia) == "" {
        RespondError(w, http.StatusBadRequest, "La frecuencia de la medicación es requerida")
        return
    }
    if strings.TrimSpace(datos.Hora_programada) == "" {
        RespondError(w, http.StatusBadRequest, "La hora programada de la medicación es requerida")
        return
    }

    actualizado, err := s.Storage.ActualizarMedicacion(id, datos)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        RespondError(w, http.StatusNotFound, "Medicacion no encontrada")
        return
    }
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al actualizar medicación")
        return
    }

    RespondJSON(w, http.StatusOK, actualizado)
}

//EliminarMedicacion DELETE /api/v1/medicaciones/{id}
func (s *MedicamentoHandler) EliminarMedicacion(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    ok, err := s.Storage.EliminarMedicacion(id)
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al eliminar medicacion")
        return
    }

    if !ok {
        RespondError(w, http.StatusNotFound, "Medicacion no encontrada")
        return
    }

    RespondJSON(w, http.StatusOK, map[string]interface{}{"message": "Medicacion eliminada correctamente"})
}

