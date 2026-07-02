package farmacia

import (
	"bytes"
	"encoding/json"
	"medicare-adulto-mayor/internal/middleware"
	"medicare-adulto-mayor/internal/models/farmacia"
	servicioFarmacia "medicare-adulto-mayor/internal/service/farmacia"
	"net/http"
	"net/http/httptest"
	"testing"
)

/* =========================
   FAKE REPOSITORY
========================= */

type FakeRepo struct{}

func (f *FakeRepo) CrearFarmacia(fr *farmacia.Farmacia) error {
	fr.ID = "test-id"
	return nil
}

func (f *FakeRepo) ListarTodas() ([]farmacia.Farmacia, error) {
	return []farmacia.Farmacia{
		{ID: "1", Nombre: "Farmacia A"},
	}, nil
}

func (f *FakeRepo) BuscarCercanas(lat, lon, radio float64) ([]farmacia.Farmacia, error) {
	return []farmacia.Farmacia{}, nil
}

/* =========================
   HELPER
========================= */

func newHandler() *ManejadorFarmacia {
	repo := &FakeRepo{}
	serv := servicioFarmacia.NuevoServicioFarmacia(repo)
	return NuevoManejadorFarmacia(serv)
}

/* =========================
   TEST 1: POST correcto
========================= */

func TestRegistrarFarmaciaHandler(t *testing.T) {

	handler := newHandler()

	body := farmacia.Farmacia{
		Nombre:    "Cruz Azul",
		Direccion: "Manta",
	}

	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/farmacias", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.RegistrarFarmacia(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Se esperaba 201, se obtuvo %d", rr.Code)
	}
}

/* =========================
   TEST 2: Método incorrecto
========================= */

func TestRegistrarFarmaciaMetodoIncorrecto(t *testing.T) {

	handler := newHandler()

	req := httptest.NewRequest("GET", "/farmacias", nil)
	rr := httptest.NewRecorder()

	handler.RegistrarFarmacia(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Se esperaba 405")
	}
}

/* =========================
   TEST 3: 401 Unauthorized
========================= */

func TestMiddlewareUnauthorized(t *testing.T) {

	mux := http.NewServeMux()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	protected := middleware.AuthMiddleware(mux)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	protected.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba 401 Unauthorized")
	}
}
