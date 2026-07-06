package models

import "time"

type Usuario struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreadoEn     time.Time `json:"creado_en"`
}