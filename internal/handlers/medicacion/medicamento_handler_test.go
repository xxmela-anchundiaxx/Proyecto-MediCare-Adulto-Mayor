package medicacion

import (
	"bytes"
	"encoding/json"
	"medicare-adulto-mayor/internal/middleware"
	"medicare-adulto-mayor/internal/models/medicacion"
	servicioMedicacion "medicare-adulto-mayor/internal/service/medicacion"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// FakeRepositorioMedicamento acts as an in-memory test double (fake) that actually saves data
type FakeRepositorioMedicamento struct {
	mu           sync.RWMutex
	medicamentos []medicacion.Medicamento
}

func (f *FakeRepositorioMedicamento) BuscarPorID(id string) (*medicacion.Medicamento, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, m := range f.medicamentos {
		if m.ID == id {
			return &m, nil
		}
	}
	return nil, nil
}

func (f *FakeRepositorioMedicamento) ListarPorPaciente(pacienteID string) ([]medicacion.Medicamento, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	var res []medicacion.Medicamento
	for _, m := range f.medicamentos {
		if m.PacienteID == pacienteID {
			res = append(res, m)
		}
	}
	return res, nil
}

func (f *FakeRepositorioMedicamento) CrearMedicamento(m *medicacion.Medicamento) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.medicamentos = append(f.medicamentos, *m)
	return nil
}

func (f *FakeRepositorioMedicamento) ActualizarStock(id string, nuevoStock int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, m := range f.medicamentos {
		if m.ID == id {
			f.medicamentos[i].Stock = nuevoStock
			return nil
		}
	}
	return nil
}

func TestRegistrarMedicamento_401_SinToken(t *testing.T) {
	// Arrange
	fakeRepo := &FakeRepositorioMedicamento{}
	servicio := servicioMedicacion.NuevoServicioMedicamento(fakeRepo)
	manejador := NuevoManejadorMedicamento(servicio)

	// Crear el handler envuelto por el Middleware de Autenticación
	handlerProtegido := middleware.AuthMiddleware(http.HandlerFunc(manejador.RegistrarMedicamento))

	body, _ := json.Marshal(medicacion.CreateMedicamentoRequest{
		PacienteID: "paciente-abc",
		Nombre:     "Metformina",
		Dosis:      "850mg",
		Frecuencia: "Cada 24 horas",
		Stock:      30,
	})

	req, err := http.NewRequest(http.MethodPost, "/medicamentos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// NO agregamos cabecera Authorization! (Prueba de error 401)
	rr := httptest.NewRecorder()

	// Act
	handlerProtegido.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("se esperaba código 401 Unauthorized, pero se obtuvo: %d", rr.Code)
	}

	var errorResponse map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("no se pudo decodificar el cuerpo de respuesta de error: %v", err)
	}

	if val, ok := errorResponse["error"]; !ok || val == "" {
		t.Error("se esperaba un mensaje de error en la respuesta")
	}
}

func TestRegistrarMedicamento_Exitoso_ConFake(t *testing.T) {
	// Arrange
	fakeRepo := &FakeRepositorioMedicamento{}
	servicio := servicioMedicacion.NuevoServicioMedicamento(fakeRepo)
	manejador := NuevoManejadorMedicamento(servicio)

	handlerProtegido := middleware.AuthMiddleware(http.HandlerFunc(manejador.RegistrarMedicamento))

	reqBody := medicacion.CreateMedicamentoRequest{
		PacienteID: "paciente-abc",
		Nombre:     "Metformina",
		Dosis:      "850mg",
		Frecuencia: "Cada 24 horas",
		Stock:      30,
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(http.MethodPost, "/medicamentos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Agregamos cabecera Authorization válida (Prueba exitosa)
	req.Header.Set("Authorization", "Bearer user-fake-id")
	rr := httptest.NewRecorder()

	// Act
	handlerProtegido.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusCreated {
		t.Fatalf("se esperaba código 201 Created, pero se obtuvo: %d. Detalle: %s", rr.Code, rr.Body.String())
	}

	var mRes medicacion.Medicamento
	if err := json.NewDecoder(rr.Body).Decode(&mRes); err != nil {
		t.Fatalf("error al decodificar respuesta exitosa: %v", err)
	}

	if mRes.ID == "" {
		t.Error("el medicamento registrado debería tener un ID auto-generado")
	}

	if mRes.Nombre != reqBody.Nombre {
		t.Errorf("se esperaba el nombre %s, pero se obtuvo %s", reqBody.Nombre, mRes.Nombre)
	}

	// Verificar en el fake en memoria que el elemento realmente se guardó
	fakeRepo.mu.RLock()
	defer fakeRepo.mu.RUnlock()
	if len(fakeRepo.medicamentos) != 1 {
		t.Errorf("se esperaba que el fake tuviera 1 medicamento guardado, pero tiene: %d", len(fakeRepo.medicamentos))
	} else if fakeRepo.medicamentos[0].Nombre != "Metformina" {
		t.Errorf("nombre incorrecto guardado en el fake: %s", fakeRepo.medicamentos[0].Nombre)
	}
}
