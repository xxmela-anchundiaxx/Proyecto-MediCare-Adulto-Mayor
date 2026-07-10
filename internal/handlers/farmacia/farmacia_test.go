package farmacia

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"proyecto-medicare-adulto-mayor/internal/models/farmacia"
	servicio "proyecto-medicare-adulto-mayor/internal/service/farmacia"
)


type mockStorage struct{}

func (m *mockStorage) CrearFarmacia(f *farmacia.Farmacia) error { 
	return nil 
}
func (m *mockStorage) BuscarPorID(id string) (*farmacia.Farmacia, error) { 
	return &farmacia.Farmacia{ID: id, Nombre: "Farmacia Test"}, nil 
}
func (m *mockStorage) ListarTodas() ([]farmacia.Farmacia, error) { 
	return []farmacia.Farmacia{{ID: "1", Nombre: "Farmacia A"}}, nil 
}
func (m *mockStorage) BuscarCercanas(lat float64, lon float64, radioKM float64) ([]farmacia.Farmacia, error) { 
	return []farmacia.Farmacia{{ID: "1", Nombre: "Farmacia Cercana"}}, nil 
}
func (m *mockStorage) ActualizarFarmacia(id string, f *farmacia.Farmacia) error { 
	return nil 
}
func (m *mockStorage) EliminarFarmacia(id string) error { 
	return nil 
}

func setupManejador() *ManejadorFarmacia {
	repo := &mockStorage{}
	srv := servicio.NuevoServicioFarmacia(repo)
	return NuevoManejadorFarmacia(srv)
}

// 1. Test Registrar Farmacia
func TestRegistrarFarmacia(t *testing.T) {
	h := setupManejador()

	nueva := farmacia.Farmacia{Nombre: "Nueva Farmacia", Direccion: "Av. Central"}
	body, _ := json.Marshal(nueva)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/farmacias", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.RegistrarFarmacia(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusCreated, w.Code)
	}
}

// 2. Test Buscar Cercanas / Listar Todas
func TestBuscarCercanas(t *testing.T) {
	h := setupManejador()

	// Probamos el caso en que lista todas (sin query params)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/farmacias", nil)
	w := httptest.NewRecorder()

	h.BuscarCercanas(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusOK, w.Code)
	}
}

// 3. Test Obtener por ID usando Chi
func TestObtenerPorID(t *testing.T) {
	h := setupManejador()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/farmacias/xyz-123", nil)
	w := httptest.NewRecorder()

	// Simulamos el parámetro de ruta de Chi {id}
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "xyz-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	h.ObtenerPorID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusOK, w.Code)
	}
}

// 4. Test Eliminar Farmacia
func TestEliminarFarmacia(t *testing.T) {
	h := setupManejador()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/farmacias/xyz-123", nil)
	w := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "xyz-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	h.EliminarFarmacia(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusOK, w.Code)
	}
}