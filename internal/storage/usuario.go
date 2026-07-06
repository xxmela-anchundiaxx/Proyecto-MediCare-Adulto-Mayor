package storage

import (
    "proyecto-medicare-adulto-mayor/internal/models"
    "time"
)

type AlmacenUsuario struct {
    usuarios []models.Usuario
    contador int
}

// Constructor
func NuevoAlmacenUsuario() *AlmacenUsuario {
    return &AlmacenUsuario{
        usuarios: []models.Usuario{},
        contador: 1,
    }
}

// CrearUsuario agrega un nuevo usuario
func (a *AlmacenUsuario) CrearUsuario(email, passwordHash string) models.Usuario {
    u := models.Usuario{
        ID:           a.contador,
        Email:        email,
        PasswordHash: passwordHash,
        CreadoEn:     time.Now(),
    }
    a.contador++
    a.usuarios = append(a.usuarios, u)
    return u
}

// BuscarUsuarioPorEmail busca un usuario por su email
func (a *AlmacenUsuario) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
    for _, u := range a.usuarios {
        if u.Email == email {
            return u, true
        }
    }
    return models.Usuario{}, false
}

// ListarUsuarios devuelve todos los usuarios
func (a *AlmacenUsuario) ListarUsuarios() []models.Usuario {
    return a.usuarios
}

// EliminarUsuario elimina un usuario por ID
func (a *AlmacenUsuario) EliminarUsuario(id int) bool {
    for i, u := range a.usuarios {
        if u.ID == id {
            a.usuarios = append(a.usuarios[:i], a.usuarios[i+1:]...)
            return true
        }
    }
    return false
}
