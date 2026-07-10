package farmacia

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"proyecto-medicare-adulto-mayor/internal/models/farmacia"
	"proyecto-medicare-adulto-mayor/internal/response"
	servicioFarmacia "proyecto-medicare-adulto-mayor/internal/service/farmacia"
)

type ManejadorFarmacia struct {
	Servicio *servicioFarmacia.ServicioFarmacia
}

func NuevoManejadorFarmacia(s *servicioFarmacia.ServicioFarmacia) *ManejadorFarmacia {
	return &ManejadorFarmacia{Servicio: s}
}

// POST /api/v1/farmacias
func (h *ManejadorFarmacia) RegistrarFarmacia(w http.ResponseWriter, r *http.Request) {
	var req farmacia.Farmacia
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.Servicio.RegistrarFarmacia(&req); err != nil {
		response.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondJSON(w, http.StatusCreated, req)
}

// GET /api/v1/farmacias
func (h *ManejadorFarmacia) BuscarCercanas(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	if latStr == "" || lonStr == "" {
		lista, err := h.Servicio.ListarTodas()
		if err != nil {
			response.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.RespondJSON(w, http.StatusOK, lista)
		return
	}

	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)
	radioKM, _ := strconv.ParseFloat(r.URL.Query().Get("radio_km"), 64)

	cercanas, err := h.Servicio.BuscarCercanas(lat, lon, radioKM)
	if err != nil {
		response.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondJSON(w, http.StatusOK, cercanas)
}

// GET /api/v1/farmacias/{id}
func (h *ManejadorFarmacia) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	f, err := h.Servicio.BuscarPorID(id)
	if err != nil {
		response.RespondError(w, http.StatusNotFound, "Farmacia no encontrada")
		return
	}

	response.RespondJSON(w, http.StatusOK, f)
}

// PUT /api/v1/farmacias/{id}
func (h *ManejadorFarmacia) ActualizarFarmacia(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req farmacia.Farmacia
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.Servicio.ActualizarFarmacia(id, &req); err != nil {
		response.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondJSON(w, http.StatusOK, req)
}

// DELETE /api/v1/farmacias/{id}
func (h *ManejadorFarmacia) EliminarFarmacia(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.Servicio.EliminarFarmacia(id); err != nil {
		response.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondJSON(w, http.StatusOK, map[string]interface{}{"message": "Farmacia eliminada correctamente"})
}