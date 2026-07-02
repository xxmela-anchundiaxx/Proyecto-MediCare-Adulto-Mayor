package storage

import (
	"database/sql"
	"errors"
	"medicare-adulto-mayor/internal/models"
)

type RepositorioUsuario interface {
	BuscarPorEmail(email string) (*models.Usuario, error)
	CrearUsuario(u *models.Usuario) error
}

type StorageUsuario struct {
	DB *sql.DB
}

func NuevoStorageUsuario(db *sql.DB) *StorageUsuario {
	return &StorageUsuario{DB: db}
}

func (s *StorageUsuario) BuscarPorEmail(email string) (*models.Usuario, error) {
	query := `SELECT id, nombre, email, password_hash, rol, fecha_creacion FROM usuarios WHERE email = ? LIMIT 1`
	var u models.Usuario
	var rolStr string

	err := s.DB.QueryRow(query, email).Scan(&u.ID, &u.Nombre, &u.Email, &u.PasswordHash, &rolStr, &u.FechaCreacion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	u.Rol = models.Role(rolStr)
	return &u, nil
}

func (s *StorageUsuario) CrearUsuario(u *models.Usuario) error {
	query := `INSERT INTO usuarios (id, nombre, email, password_hash, rol, fecha_creacion) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.DB.Exec(query, u.ID, u.Nombre, u.Email, u.PasswordHash, string(u.Rol), u.FechaCreacion)
	return err
}
