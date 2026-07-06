package medicacion

import (
    "encoding/json"
    "errors"
    "net/http"
    "strconv"
    "strings"
    "time"
    "log"

    "github.com/go-chi/chi/v5"
    "gorm.io/gorm"

    "proyecto-medicare-adulto-mayor/internal/models/medicacion"
    "proyecto-medicare-adulto-mayor/internal/response"
)

// ListarMedicacion GET /api/v1/medicaciones
/*func (s *Server) ListarMedicacion(w http.ResponseWriter, _ *http.Request) {
    meds, err := s.Medicacion.Listar()
    if err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al listar medicación")
        return
    }
    response.RespondJSON(w, http.StatusOK, meds)
}*/

func (s *Server) ListarMedicacion(w http.ResponseWriter, _ *http.Request) {
    meds, err := s.Medicacion.Listar()
    if err != nil {
        // ¡Esto imprimirá el culpable real en tu terminal!
        log.Printf("[ERROR CRÍTICO] Falló ListarMedicacion en la Base de Datos: %v", err)
        
        response.RespondError(w, http.StatusInternalServerError, "Error al listar medicación")
        return
    }
    response.RespondJSON(w, http.StatusOK, meds)
}

// ObtenerMedicacion GET /api/v1/medicaciones/{id}
func (s *Server) ObtenerMedicacion(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    med, err := s.Medicacion.Obtener(id)
    if err != nil {
        response.RespondError(w, http.StatusNotFound, "Medicacion no encontrada")
        return
    }

    response.RespondJSON(w, http.StatusOK, med)
}

// CrearMedicacion POST /api/v1/medicaciones
func (s *Server) CrearMedicacion(w http.ResponseWriter, r *http.Request) {
    var nueva medicacion.Medicacion

    if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
        response.RespondError(w, http.StatusBadRequest, "Datos de medicación inválidos: "+err.Error())
        return
    }

    if strings.TrimSpace(nueva.Nombre) == "" {
        response.RespondError(w, http.StatusBadRequest, "El nombre de la medicación es requerido")
        return
    }
    if strings.TrimSpace(nueva.Dosis) == "" {
        response.RespondError(w, http.StatusBadRequest, "La dosis de la medicación es requerida")
        return
    }
    if strings.TrimSpace(nueva.Frecuencia) == "" {
        response.RespondError(w, http.StatusBadRequest, "La frecuencia de la medicación es requerida")
        return
    }
    if strings.TrimSpace(nueva.Hora_programada) == "" {
        response.RespondError(w, http.StatusBadRequest, "La hora programada de la medicación es requerida")
        return
    }
    if nueva.Inicio_tratamiento.IsZero() {
        response.RespondError(w, http.StatusBadRequest, "La fecha de inicio del tratamiento es requerida")
        return
    }
    if nueva.Fecha_creacion.IsZero() {
        nueva.Fecha_creacion = time.Now()
    }

    creada, err := s.Medicacion.Crear(nueva)
    if err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al guardar en el almacén: "+err.Error())
        return
    }

    response.RespondJSON(w, http.StatusCreated, creada)
}

// ActualizarMedicacion PUT /api/v1/medicaciones/{id}
func (s *Server) ActualizarMedicacion(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    var datos medicacion.Medicacion
    if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
        response.RespondError(w, http.StatusBadRequest, "Datos de medicación inválidos: "+err.Error())
        return
    }

    if strings.TrimSpace(datos.Nombre) == "" {
        response.RespondError(w, http.StatusBadRequest, "El nombre de la medicación es requerido")
        return
    }
    if strings.TrimSpace(datos.Dosis) == "" {
        response.RespondError(w, http.StatusBadRequest, "La dosis de la medicación es requerida")
        return
    }
    if strings.TrimSpace(datos.Frecuencia) == "" {
        response.RespondError(w, http.StatusBadRequest, "La frecuencia de la medicación es requerida")
        return
    }
    if strings.TrimSpace(datos.Hora_programada) == "" {
        response.RespondError(w, http.StatusBadRequest, "La hora programada de la medicación es requerida")
        return
    }

    actualizado, err := s.Medicacion.Actualizar(id, datos)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        response.RespondError(w, http.StatusNotFound, "Medicacion no encontrada")
        return
    }
    if err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al actualizar medicación")
        return
    }

    response.RespondJSON(w, http.StatusOK, actualizado)
}

// EliminarMedicacion DELETE /api/v1/medicaciones/{id}
func (s *Server) EliminarMedicacion(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    if err := s.Medicacion.Eliminar(id); err != nil {
        response.RespondError(w, http.StatusInternalServerError, "Error al eliminar medicacion")
        return
    }

    response.RespondJSON(w, http.StatusOK, map[string]interface{}{"message": "Medicacion eliminada correctamente"})
}
