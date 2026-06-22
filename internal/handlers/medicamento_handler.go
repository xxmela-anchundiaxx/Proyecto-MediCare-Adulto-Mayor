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


//metodos de pacientes
type PacienteHandler struct {
    Storage storage.Almacen
}

func NewPacienteHandler(s storage.Almacen) *PacienteHandler{
    return &PacienteHandler{Storage: s}
}

// ListarPaciente GET api/v1/pacientes
func (h *PacienteHandler) ListarPacientes(w http.ResponseWriter, r *http.Request) {
    pacientes, err := h.Storage.ListarPacientes()
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al listar pacientes")
        return
    }
    RespondJSON(w, http.StatusOK, pacientes)
}

// BuscarPaciente por ID GET api/v1/pacientes/{id}
func (h *PacienteHandler) BuscarPacientePorID(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    paciente, err := h.Storage.BuscarPacientePorID(id)
    if err != nil {
        RespondError(w, http.StatusNotFound, "Paciente no encontrado")
        return
    }
    RespondJSON(w, http.StatusOK, paciente)
}

// CrearPaciente POST api/v1/pacientes
func (h *PacienteHandler) CrearPaciente(w http.ResponseWriter, r *http.Request) {
    var p models.Paciente
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        RespondError(w, http.StatusBadRequest, "JSON inválido")
        return
    }

    nuevo, err := h.Storage.CrearPaciente(p)
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al crear paciente")
        return
    }
    RespondJSON(w, http.StatusCreated, nuevo)
}

// ActualzarPaciente PUT api/v1/pacientes/{id}
func (h *PacienteHandler) ActualizarPaciente(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    var p models.Paciente
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        RespondError(w, http.StatusBadRequest, "JSON inválido")
        return
    }

    actualizado, err := h.Storage.ActualizarPaciente(id, p)
    if err != nil {
        // Si el error es que no existe el registro, puedes devolver 404
        if errors.Is(err, gorm.ErrRecordNotFound) {
            RespondError(w, http.StatusNotFound, "Paciente no encontrado")
            return
        }
        RespondError(w, http.StatusInternalServerError, "Error al actualizar paciente")
        return
    }
    RespondJSON(w, http.StatusOK, actualizado)
}


// EliminarPaciente DELETE api/v1/pacientes/{id}
func (h *PacienteHandler) EliminarPaciente(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    ok, err := h.Storage.EliminarPaciente(id)
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al eliminar paciente")
        return
    }
    if !ok {
        RespondError(w, http.StatusNotFound, "Paciente no encontrado")
        return
    }
    RespondJSON(w, http.StatusNoContent, nil)
}


type HistorialHandler struct {
    Storage storage.Almacen
}

func NewHistorialHandler(s storage.Almacen) *HistorialHandler {
    return &HistorialHandler{Storage: s}
}

// ListarHistorial GET api/v1/historial
func (h *HistorialHandler) ListarHistorial(w http.ResponseWriter, r *http.Request) {
    historiales, err := h.Storage.ListarHistorial()
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al listar historiales")
        return
    }
    RespondJSON(w, http.StatusOK, historiales)
}

// Buscarhistorial por ID GET api/v1/historial/{id}
func (h *HistorialHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    registro, err := h.Storage.BuscarHistorialPorID(id)
    if err != nil {
        RespondError(w, http.StatusNotFound, "Historial no encontrado")
        return
    }
    RespondJSON(w, http.StatusOK, registro)
}


// CrearHistorial POST api/v1/historial
func (h *HistorialHandler) Crear(w http.ResponseWriter, r *http.Request) {
    var hMed models.HistorialMedicacion
    if err := json.NewDecoder(r.Body).Decode(&hMed); err != nil {
        RespondError(w, http.StatusBadRequest, "JSON inválido")
        return
    }

    nuevo, err := h.Storage.CrearHistorial(hMed)
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al crear historial")
        return
    }
    RespondJSON(w, http.StatusCreated, nuevo)
}

// ActualizarHistorial PUT api/historial/{id}
func (h *HistorialHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    var datos models.HistorialMedicacion
    if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
        RespondError(w, http.StatusBadRequest, "JSON inválido")
        return
    }

    actualizado, err := h.Storage.ActualizarHistorial(id, datos)
    if err != nil {
        // Si el error es que no existe el registro
        if errors.Is(err, gorm.ErrRecordNotFound) {
            RespondError(w, http.StatusNotFound, "Historial no encontrado")
            return
        }
        RespondError(w, http.StatusInternalServerError, "Error al actualizar historial")
        return
    }
    RespondJSON(w, http.StatusOK, actualizado)
}

// EliminarHistorial DELETE api/v1/historial/{id}
func (h *HistorialHandler) Eliminar(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        RespondError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    ok, err := h.Storage.EliminarHistorial(id)
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al eliminar historial")
        return
    }
    if !ok {
        RespondError(w, http.StatusNotFound, "Historial no encontrado")
        return
    }
    RespondJSON(w, http.StatusNoContent, nil)
}

// ObtenerMedicacion de un paciente GET /api/v1/pacientes/{id}/medicaciones
func (h *MedicamentoHandler) ListarPorPaciente(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    pacienteID, _ := strconv.Atoi(idStr)

    medicaciones, err := h.Storage.ListarMedicacionPorPaciente(pacienteID)
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al listar medicaciones")
        return
    }
    RespondJSON(w, http.StatusOK, medicaciones)
}

// ObtenerHistorial de un paciente GET /api/v1/pacientes/{id}/historial
func (h *HistorialHandler) ListarPorPaciente(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    pacienteID, _ := strconv.Atoi(idStr)

    historiales, err := h.Storage.ListarHistorialPorPaciente(pacienteID)
    if err != nil {
        RespondError(w, http.StatusInternalServerError, "Error al listar historial")
        return
    }
    RespondJSON(w, http.StatusOK, historiales)
}

