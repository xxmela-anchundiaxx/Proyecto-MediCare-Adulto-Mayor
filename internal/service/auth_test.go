package service

import (
	"errors"
	"testing"
	"time"

	"proyecto-medicare-adulto-mayor/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepository simula el comportamiento de storage.UserRepository sin tocar la DB real
type mockUserRepository struct {
	onCrearUsuario           func(u models.Usuario) (models.Usuario, error)
	onBuscarUsuarioPorEmail func(email string) (models.Usuario, bool)
}

func (m *mockUserRepository) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	return m.onCrearUsuario(u)
}

func (m *mockUserRepository) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	return m.onBuscarUsuarioPorEmail(email)
}

// --- PRUEBAS DE REGISTRO ---

func TestAuthService_Registrar_Exitoso(t *testing.T) {
	mockRepo := &mockUserRepository{
		onBuscarUsuarioPorEmail: func(email string) (models.Usuario, bool) {
			return models.Usuario{}, false // El email no existe previamente
		},
		onCrearUsuario: func(u models.Usuario) (models.Usuario, error) {
			u.ID = 1 // Simulamos que la DB le asigna el ID 1
			return u, nil
		},
	}

	svc := NewAuthService(mockRepo)
	u, err := svc.Registrar("test@medicare.com", "password123")

	if err != nil {
		t.Fatalf("No se esperaba un error, se obtuvo: %v", err)
	}
	if u.ID != 1 {
		t.Errorf("Se esperaba ID 1, se obtuvo %d", u.ID)
	}
	if u.Email != "test@medicare.com" {
		t.Errorf("Se esperaba email formateado en minúsculas")
	}
}

func TestAuthService_Registrar_EmailEnUso(t *testing.T) {
	mockRepo := &mockUserRepository{
		onBuscarUsuarioPorEmail: func(email string) (models.Usuario, bool) {
			return models.Usuario{ID: 1, Email: email}, true // El email YA existe
		},
	}

	svc := NewAuthService(mockRepo)
	_, err := svc.Registrar("ya_existe@medicare.com", "password123")

	if !errors.Is(err, ErrEmailEnUso) {
		t.Errorf("Se esperaba ErrEmailEnUso, se obtuvo: %v", err)
	}
}

// --- PRUEBAS DE LOGIN ---

func TestAuthService_Login_Exitoso(t *testing.T) {
	passwordOriginal := "segura123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(passwordOriginal), bcrypt.DefaultCost)

	mockRepo := &mockUserRepository{
		onBuscarUsuarioPorEmail: func(email string) (models.Usuario, bool) {
			return models.Usuario{
				ID:           42,
				Email:        email,
				PasswordHash: string(hash),
			}, true
		},
	}

	svc := NewAuthService(mockRepo)
	token, err := svc.Login("test@medicare.com", passwordOriginal)

	if err != nil {
		t.Fatalf("Login falló inesperadamente: %v", err)
	}
	if token == "" {
		t.Error("El token generado no debería estar vacío")
	}
}

func TestAuthService_Login_CredencialesInvalidas(t *testing.T) {
	mockRepo := &mockUserRepository{
		onBuscarUsuarioPorEmail: func(email string) (models.Usuario, bool) {
			return models.Usuario{}, false // No existe el usuario
		},
	}

	svc := NewAuthService(mockRepo)
	_, err := svc.Login("no_existoc@medicare.com", "123")

	if !errors.Is(err, ErrCredencialesInvalidas) {
		t.Errorf("Se esperaba ErrCredencialesInvalidas, se obtuvo: %v", err)
	}
}

// --- PRUEBAS DE VERIFICACIÓN DE JWT TOKENS ---

func TestAuthService_VerificarToken_Valido(t *testing.T) {
	mockRepo := &mockUserRepository{}
	svc := NewAuthService(mockRepo)

	// Creamos un usuario ficticio para generarle un token legítimo primero
	usuarioFake := models.Usuario{ID: 100, Email: "token@medicare.com"}
	tokenStr, err := svc.generarToken(usuarioFake)
	if err != nil {
		t.Fatalf("Error preparando el token del test: %v", err)
	}

	// Probamos la verificación de dicho token
	uid, err := svc.VerificarToken(tokenStr)
	if err != nil {
		t.Fatalf("No se pudo verificar un token válido: %v", err)
	}
	if uid != 100 {
		t.Errorf("Se esperaba recuperar el UsuarioID 100, se obtuvo: %d", uid)
	}
}

func TestAuthService_VerificarToken_InvalidooExpirado(t *testing.T) {
	mockRepo := &mockUserRepository{}
	svc := NewAuthService(mockRepo)

	// Caso 1: Un token que es texto corrupto cualquiera
	_, err := svc.VerificarToken("un.token.totalmente.invalido")
	if !errors.Is(err, ErrCredencialesInvalidas) {
		t.Errorf("Se esperaba error por token corrupto")
	}

	// Caso 2: Un token legítimo firmado con otra firma / secreto ajeno
	claimsMalos := Claims{
		UsuarioID: 5,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tokenMalo := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsMalos)
	tokenMaloStr, _ := tokenMalo.SignedString([]byte("UN_SECRETO_HACKER_DIFERENTE"))

	_, err = svc.VerificarToken(tokenMaloStr)
	if !errors.Is(err, ErrCredencialesInvalidas) {
		t.Errorf("Se esperaba que fallara al ser firmado con una llave diferente")
	}
}