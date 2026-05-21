package models

import "time"

type Usuario struct {
	ID                int       `json:"id"`
	Nombre_completo   string    `json:"nombre_completo"`
	Email             string    `json:"email"`
	Password          string    `json:"password,omitempty"` 
	Rol               string    `json:"rol"`                
	Fecha_registro    time.Time `json:"fecha_registro"`
}