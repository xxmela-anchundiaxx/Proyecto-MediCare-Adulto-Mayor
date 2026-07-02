package models

import "time"

type Role string

const (
	RolePaciente      Role = "paciente"
	RoleFamiliar      Role = "familiar"
	RoleMedico         Role = "medico"
	RoleAdministrador Role = "administrador"
)

type Usuario struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	Nombre        string    `json:"nombre"`
	Email         string    `json:"email" gorm:"unique"`
	PasswordHash  string    `json:"-"`
	Rol           Role      `json:"rol"`
	FechaCreacion time.Time `json:"fecha_creacion"`
}

func (Usuario) TableName() string {
	return "usuarios"
}


type UserRegisterRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Rol      Role   `json:"rol"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"nombre"`
	Email string ` Greenwood:"email"`
	Rol   string `json:"rol"`
}
