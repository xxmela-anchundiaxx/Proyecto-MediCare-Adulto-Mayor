package handlers

import (
    "encoding/json"
    "net/http"

    "proyecto-medicare-adulto-mayor/internal/response"
    //"proyecto-medicare-adulto-mayor/internal/service"
)

type credenciales struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (s *Server) Registrar(w http.ResponseWriter, r *http.Request) {
    var c credenciales
    if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
        response.RespondError(w, http.StatusBadRequest, "JSON inválido")
        return
    }
    usuario, err := s.Auth.Registrar(c.Email, c.Password)
    if err != nil {
        response.RespondError(w, statusDeError(err), err.Error())
        return
    }
    response.RespondJSON(w, http.StatusCreated, usuario)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
    var c credenciales
    if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
        response.RespondError(w, http.StatusBadRequest, "JSON inválido")
        return
    }
    token, err := s.Auth.Login(c.Email, c.Password)
    if err != nil {
        response.RespondError(w, statusDeError(err), err.Error())
        return
    }
    response.RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}

