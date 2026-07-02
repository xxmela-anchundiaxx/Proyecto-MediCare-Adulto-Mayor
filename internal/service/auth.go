package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"medicare-adulto-mayor/internal/models"
	"medicare-adulto-mayor/internal/storage"
	"time"

	"github.com/google/uuid"
)

type ServicioAuth struct {
	Repo storage.RepositorioUsuario
}

func NuevoServicioAuth(repo storage.RepositorioUsuario) *ServicioAuth {
	return &ServicioAuth{Repo: repo}
}

func (s *ServicioAuth) Registrar(req models.UserRegisterRequest) (*models.Usuario, error) {
	if req.Nombre == "" || req.Email == "" || req.Password == "" {
		return nil, ErrDatosInvalidos
	}

	existente, err := s.Repo.BuscarPorEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existente != nil {
		return nil, ErrUsuarioExistente
	}

	hash := hashPassword(req.Password)

	u := &models.Usuario{
		ID:            uuid.New().String(),
		Nombre:        req.Nombre,
		Email:         req.Email,
		PasswordHash:  hash,
		Rol:           req.Rol,
		FechaCreacion: time.Now(),
	}

	if err := s.Repo.CrearUsuario(u); err != nil {
		return nil, fmt.Errorf("error al guardar usuario: %w", err)
	}

	return u, nil
}

func (s *ServicioAuth) Login(req models.UserLoginRequest) (*models.Usuario, error) {
	u, err := s.Repo.BuscarPorEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrCredencialesInvalidas
	}

	hash := hashPassword(req.Password)
	if u.PasswordHash != hash {
		return nil, ErrCredencialesInvalidas
	}

	return u, nil
}

// hashPassword aplica un hash SHA-256 para demostración y simplicidad sin dependencias complejas de bcrypt
func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + "medicare-salt-12345"))
	return hex.EncodeToString(hasher.Sum(nil))
}
