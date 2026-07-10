package monitoreo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	models "proyecto-medicare-adulto-mayor/internal/models/monitoreo"
	servicio "proyecto-medicare-adulto-mayor/internal/service/monitoreo"
)

// mockMonitoreoRepo implementa la interfaz storage.RepositorioMonitoreo
// para inyectarla en el servicio durante las pruebas del Handler
type mockMonitoreoRepo struct {
	onListar     func() ([]models.CuidadorPaciente, error)
	onBuscarID   func(id int) (models.CuidadorPaciente, error)
	onCrear      func(rel models.CuidadorPaciente) (models.CuidadorPaciente, error)
	onActualizar func(id int, rel models.CuidadorPaciente) (models.CuidadorPaciente, error)
	onEliminar   func(id int) (bool, error)
}

func (m *mockMonitoreoRepo) ListarRelaciones() ([]models.CuidadorPaciente, error) { return m.onListar() }
func (m *mockMonitoreoRepo) BuscarRelacionPorID(id int) (models.CuidadorPaciente, error) { return m.onBuscarID(id) }
func (m *mockMonitoreoRepo) CrearRelacion(rel models.CuidadorPaciente) (models.CuidadorPaciente, error) { return m.onCrear(rel) }
func (m *mockMonitoreoRepo) ActualizarRelacion(id int, rel models.CuidadorPaciente) (models.CuidadorPaciente, error) { return m.onActualizar(id, rel) }
func (m *mockMonitoreoRepo) EliminarRelacion(id int) (bool, error) { return m.onEliminar(id) }

// --- PRUEBAS DE LOS HANDLERS ---

func TestCrearRelacionHandler_Exitoso(t *testing.T) {
	repo := &mockMonitoreoRepo{
		onCrear: func(rel models.CuidadorPaciente) (models.CuidadorPaciente, error) {
			rel.ID = 100
			return rel, nil
		},
	}
	svc := servicio.NuevoServicioMonitoreo(repo)
	handler := NewManejadorMonitoreo(svc)

	payload := models.CuidadorPaciente{
		CuidadorID: 1,
		PacienteID: 2,
		Relacion:   "Enfermero",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/monitoreo", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.CrearRelacionHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Se esperaba código 201, se obtuvo: %d", rr.Code)
	}

	var respondida models.CuidadorPaciente
	json.Unmarshal(rr.Body.Bytes(), &respondida)

	if respondida.ID != 100 {
		t.Errorf("Se esperaba ID asignado 100, llegó: %d", respondida.ID)
	}
}

func TestObtenerRelacionPorIDHandler_Exitoso(t *testing.T) {
	repo := &mockMonitoreoRepo{
		onBuscarID: func(id int) (models.CuidadorPaciente, error) {
			return models.CuidadorPaciente{ID: id, CuidadorID: 5, PacienteID: 10, Relacion: "Hijo"}, nil
		},
	}
	svc := servicio.NuevoServicioMonitoreo(repo)
	handler := NewManejadorMonitoreo(svc)

	req, _ := http.NewRequest(http.MethodGet, "/monitoreo/12", nil)
	
	// Como usamos Chi Router, debemos simular el parámetro URL {"id": "12"} usando su contexto
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "12")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	rr := httptest.NewRecorder()
	handler.ObtenerRelacionPorIDHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Se esperaba código 200, se obtuvo: %d", rr.Code)
	}

	var respondida models.CuidadorPaciente
	json.Unmarshal(rr.Body.Bytes(), &respondida)

	if respondida.ID != 12 || respondida.Relacion != "Hijo" {
		t.Errorf("Los datos devueltos por el handler no coinciden con el mock")
	}
}

func TestObtenerRelacionPorIDHandler_IDInvalido(t *testing.T) {
	svc := servicio.NuevoServicioMonitoreo(&mockMonitoreoRepo{})
	handler := NewManejadorMonitoreo(svc)

	req, _ := http.NewRequest(http.MethodGet, "/monitoreo/abc", nil)
	
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "abc") // Forzamos un string no convertible a int
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	rr := httptest.NewRecorder()
	handler.ObtenerRelacionPorIDHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Se esperaba un código 400 por ID inválido, llegó: %d", rr.Code)
	}
}

func TestEliminarRelacionHandler_Exitoso(t *testing.T) {
	repo := &mockMonitoreoRepo{
		onEliminar: func(id int) (bool, error) {
			return true, nil
		},
	}
	svc := servicio.NuevoServicioMonitoreo(repo)
	handler := NewManejadorMonitoreo(svc)

	req, _ := http.NewRequest(http.MethodDelete, "/monitoreo/1", nil)
	
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	rr := httptest.NewRecorder()
	handler.EliminarRelacionHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Se esperaba código 200, se obtuvo: %d", rr.Code)
	}
}