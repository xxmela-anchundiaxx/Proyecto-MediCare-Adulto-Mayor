package handlers

import (
	"encoding/json"
	"errors"
	"medicare-adulto-mayor/internal/handlers/respond"
	"medicare-adulto-mayor/internal/models"
	"medicare-adulto-mayor/internal/service"
	"net/http"
)

type ManejadorAuth struct {
	Servicio *service.ServicioAuth
}

func NuevoManejadorAuth(s *service.ServicioAuth) *ManejadorAuth {
	return &ManejadorAuth{Servicio: s}
}

func (h *ManejadorAuth) Registro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	var req models.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.ResponderError(w, http.StatusBadRequest, "cuerpo de petición inválido")
		return
	}

	u, err := h.Servicio.Registrar(req)
	if err != nil {
		if errors.Is(err, service.ErrUsuarioExistente) {
			respond.ResponderError(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, service.ErrDatosInvalidos) {
			respond.ResponderError(w, http.StatusBadRequest, err.Error())
			return
		}
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Devolvemos el usuario creado y un token de prueba (su propio ID)
	respond.ResponderJSON(w, http.StatusCreated, map[string]interface{}{
		"usuario: ": u,
		"token":     u.ID,
	})
}

func (h *ManejadorAuth) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.ResponderError(w, http.StatusMethodNotAllowed, "método no permitido")
		return
	}

	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.ResponderError(w, http.StatusBadRequest, "cuerpo de petición inválido")
		return
	}

	u, err := h.Servicio.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrCredencialesInvalidas) {
			respond.ResponderError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respond.ResponderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Enviar respuesta exitosa con el token (su ID de usuario)
	respond.ResponderJSON(w, http.StatusOK, map[string]interface{}{
		"id":     u.ID,
		"nombre": u.Nombre,
		"email":  u.Email,
		"rol":    u.Rol,
		"token":  u.ID,
	})
}
