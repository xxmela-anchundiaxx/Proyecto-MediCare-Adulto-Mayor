package storage

import (
    "proyecto-medicare-adulto-mayor/internal/models"
    "gorm.io/gorm"
)

type AlmacenUsuario struct {
    db *gorm.DB
}

// NuevoAlmacenUsuario ahora recibe la conexión ya abierta.
func NuevoAlmacenUsuario(db *gorm.DB) *AlmacenUsuario {
    // La migración es mejor hacerla en el main.go junto con las otras tablas
    return &AlmacenUsuario{db: db}
}

func (a *AlmacenUsuario) CrearUsuario(u models.Usuario) (models.Usuario, error) {
    if err := a.db.Create(&u).Error; err != nil {
        return models.Usuario{}, err
    }
    return u, nil
}

func (a *AlmacenUsuario) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
    var u models.Usuario
    // Usamos First para buscar; si no encuentra nada, GORM retorna error
    if err := a.db.Where("email = ?", email).First(&u).Error; err != nil {
        return models.Usuario{}, false
    }
    return u, true
}