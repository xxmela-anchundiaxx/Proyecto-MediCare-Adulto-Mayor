package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt" // Necesario para generar el hash de prueba
	models "proyecto-medicare-adulto-mayor/internal/models"
	"proyecto-medicare-adulto-mayor/internal/service"
)

// mockUserRepository implementa la interfaz necesaria para el servicio de autenticación
type mockUserRepository struct {
	onCrearUsuario          func(u models.Usuario) (models.Usuario, error)
	onBuscarUsuarioPorEmail func(email string) (models.Usuario, bool)
}

func (m *mockUserRepository) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	return m.onCrearUsuario(u)
}

func (m *mockUserRepository) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	return m.onBuscarUsuarioPorEmail(email)
}

// --- PRUEBAS ---

func TestRegistrar_Exitoso(t *testing.T) {
	mockRepo := &mockUserRepository{
		onCrearUsuario: func(u models.Usuario) (models.Usuario, error) {
			u.ID = 1
			return u, nil
		},
		onBuscarUsuarioPorEmail: func(email string) (models.Usuario, bool) {
			return models.Usuario{}, false
		},
	}

	authSvc := service.NewAuthService(mockRepo)
	server := &Server{Auth: authSvc}

	payload := credenciales{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/registrar", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	server.Registrar(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Se esperaba código 201 Created, se obtuvo: %d", rr.Code)
	}
}

func TestLogin_Exitoso(t *testing.T) {
	passwordTest := "password123"
	// Generamos un hash real para que bcrypt.CompareHashAndPassword pueda validarlo
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordTest), bcrypt.DefaultCost)

	mockRepo := &mockUserRepository{
		onCrearUsuario: func(u models.Usuario) (models.Usuario, error) {
			return u, nil
		},
		onBuscarUsuarioPorEmail: func(email string) (models.Usuario, bool) {
			return models.Usuario{
				ID:           1,
				Email:        email,
				PasswordHash: string(hashedPassword), // Usamos el hash generado
			}, true
		},
	}

	authSvc := service.NewAuthService(mockRepo)
	server := &Server{Auth: authSvc}

	payload := credenciales{
		Email:    "login@example.com",
		Password: passwordTest, // Enviamos el password en plano
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	server.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Se esperaba código 200 OK, se obtuvo: %d", rr.Code)
	}
}