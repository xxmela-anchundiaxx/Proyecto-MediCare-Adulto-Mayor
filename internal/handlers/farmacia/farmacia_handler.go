package farmacia

import (
	"encoding/json"
	"medicare-adulto-mayor/internal/handlers/respond"
	"medicare-adulto-mayor/internal/models/farmacia"
	servicioFarmacia "medicare-adulto-mayor/internal/service/farmacia"
	"net/http"
	"strconv"
)

type ManejadorFarmacia struct {
	Servicio *servicioFarmacia.ServicioFarmacia
}

func NuevoManejadorFarmacia(s *servicioFarmacia.ServicioFarmacia) *ManejadorFarmacia {
	return &ManejadorFarmacia{Servicio: s}
}

func (h *ManejadorFarmacia) RegistrarFarmacia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	var req farmacia.Farmacia
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.ResponderError(w, http.StatusBadRequest, "cuerpo de petición inválido")
		return
	}

	if err := h.Servicio.RegistrarFarmacia(&req); err != nil {
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusCreated, req)
}

func (h *ManejadorFarmacia) BuscarCercanas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	radioStr := r.URL.Query().Get("radio_km")

	if latStr == "" || lonStr == "" {
		// Retornar todas las farmacias si no hay coordenadas de búsqueda
		lista, err := h.Servicio.ListarTodas()
		if err != nil {
			respond.ResponderError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respond.ResponderJSON(w, http.StatusOK, lista)
		return
	}

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)
	if err1 != nil || err2 != nil {
		respond.ResponderError(w, http.StatusBadRequest, "coordenadas latitud/longitud inválidas")
		return
	}

	radioKM := 5.0
	if radioStr != "" {
		if rVal, err := strconv.ParseFloat(radioStr, 64); err == nil {
			radioKM = rVal
		}
	}

	cercanas, err := h.Servicio.BuscarCercanas(lat, lon, radioKM)
	if err != nil {
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.ResponderJSON(w, http.StatusOK, cercanas)
}
