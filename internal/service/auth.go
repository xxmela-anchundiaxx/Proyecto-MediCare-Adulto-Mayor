package service

import (
	"errors"
	"strings"
	"time"

	"proyecto-medicare-adulto-mayor/internal/models"
	"proyecto-medicare-adulto-mayor/internal/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretoJWT = []byte("medicare-secret-2024")
var duracionToken = time.Hour * 24

// Descomentado para evitar el error de compilación
//var ErrEmailEnUso = errors.New("el email ya está en uso") 
var ErrCredencialesInvalidos = errors.New("credenciales inválidas")
var ErrEmailVacio = errors.New("email o password vacíos")

type Claims struct {
	UsuarioID int `json:"uid"`
	jwt.RegisteredClaims
}

// 1. Definimos la estructura AuthService que espera tu main.go
type AuthService struct {
	repo *storage.AlmacenUsuario // O cambia esto a tu repositorio GORM definitivo si lo requieres
}

// 2. Definimos el constructor NewAuthService que llamas en main.go
func NewAuthService(repo *storage.AlmacenUsuario) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

// 3. Convertimos Registrar y Login en métodos de AuthService
func (s *AuthService) Registrar(email, password string) (models.Usuario, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || strings.TrimSpace(password) == "" {
		return models.Usuario{}, ErrEmailVacio
	}

	// Usamos s.repo en lugar de la variable global antigua
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}

	return s.repo.CrearUsuario(email, string(hash)), nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	// Buscar usuario en el repositorio inyectado
	usuario, existe := s.repo.BuscarUsuarioPorEmail(email)
	if !existe {
		return "", ErrCredencialesInvalidos
	}

	// Verificar contraseña
	if bcrypt.CompareHashAndPassword([]byte(usuario.PasswordHash), []byte(password)) != nil {
		return "", ErrCredencialesInvalidos
	}

	// Generar token JWT
	return s.generarToken(usuario)
}

func (s *AuthService) generarToken(u models.Usuario) (string, error) {
	claims := Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracionToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretoJWT)
}

func (s *AuthService) VerificarToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidos
		}
		return secretoJWT, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrCredencialesInvalidos
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, ErrCredencialesInvalidos
	}
	return claims.UsuarioID, nil
}