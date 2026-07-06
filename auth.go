package handlers

import (
    "encoding/json"
    "net/http"
    "proyecto-medicare-adulto-mayor/internal/service"
    "proyecto-medicare-adulto-mayor/internal/response"
)

type AuthHandler struct {
    auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
    return &AuthHandler{auth: auth}
}

type credenciales struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *AuthHandler) Registrar(w http.ResponseWriter, r *http.Request) {
    var creds credenciales
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        response.RespondError(w, http.StatusBadRequest, "JSON inválido")
        return
    }
    usuario, err := h.auth.Registrar(creds.Email, creds.Password)
    if err != nil {
        response.RespondError(w, http.StatusBadRequest, err.Error())
        return
    }
    response.RespondJSON(w, http.StatusCreated, usuario)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var creds credenciales
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        response.RespondError(w, http.StatusBadRequest, "JSON inválido")
        return
    }
    token, err := h.auth.Login(creds.Email, creds.Password)
    if err != nil {
        response.RespondError(w, http.StatusUnauthorized, err.Error())
        return
    }
    response.RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}
