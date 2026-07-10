package farmacia

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"proyecto-medicare-adulto-mayor/internal/models/farmacia"
)

// mockFarmaciaRepo implementa storage.RepositorioFarmacia para las pruebas
type mockFarmaciaRepo struct {
	onCrear    func(f *farmacia.Farmacia) error
	onBuscarID func(id string) (*farmacia.Farmacia, error)
	onListar   func() ([]farmacia.Farmacia, error)
	onCercanas func(lat, lon, radio float64) ([]farmacia.Farmacia, error)
	onAct      func(id string, f *farmacia.Farmacia) error
	onEliminar func(id string) error
}

func (m *mockFarmaciaRepo) CrearFarmacia(f *farmacia.Farmacia) error { return m.onCrear(f) }
func (m *mockFarmaciaRepo) BuscarPorID(id string) (*farmacia.Farmacia, error) { return m.onBuscarID(id) }
func (m *mockFarmaciaRepo) ListarTodas() ([]farmacia.Farmacia, error) { return m.onListar() }
func (m *mockFarmaciaRepo) BuscarCercanas(lat, lon, radioKM float64) ([]farmacia.Farmacia, error) {
	return m.onCercanas(lat, lon, radioKM)
}
func (m *mockFarmaciaRepo) ActualizarFarmacia(id string, f *farmacia.Farmacia) error { return m.onAct(id, f) }
func (m *mockFarmaciaRepo) EliminarFarmacia(id string) error { return m.onEliminar(id) }

// --- PRUEBAS ---

func TestRegistrarFarmacia_Exitoso(t *testing.T) {
	repo := &mockFarmaciaRepo{
		onCrear: func(f *farmacia.Farmacia) error {
			return nil
		},
	}
	svc := NuevoServicioFarmacia(repo)

	f := &farmacia.Farmacia{
		Nombre:    "Farmacia Central",
		Direccion: "Av. Siempre Viva 123",
	}

	err := svc.RegistrarFarmacia(f)

	if err != nil {
		t.Fatalf("No se esperaba un error, se obtuvo: %v", err)
	}

	// Validar que se haya generado un UUID automáticamente
	if _, parseErr := uuid.Parse(f.ID); parseErr != nil {
		t.Errorf("Se esperaba un ID con formato UUID válido, se obtuvo: %s", f.ID)
	}
}

func TestRegistrarFarmacia_ErroresValidacion(t *testing.T) {
	svc := NuevoServicioFarmacia(&mockFarmaciaRepo{})

	tests := []struct {
		name string
		f    *farmacia.Farmacia
	}{
		{"Nombre vacío", &farmacia.Farmacia{Nombre: "", Direccion: "Av. Central"}},
		{"Dirección vacía", &farmacia.Farmacia{Nombre: "Farmacia A", Direccion: ""}},
		{"Ambos vacíos", &farmacia.Farmacia{Nombre: "", Direccion: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.RegistrarFarmacia(tt.f)
			if err == nil {
				t.Error("Se esperaba un error de validación, pero el resultado fue exitoso")
			}
		})
	}
}

func TestBuscarPorID(t *testing.T) {
	targetID := "123-uuid"
	repo := &mockFarmaciaRepo{
		onBuscarID: func(id string) (*farmacia.Farmacia, error) {
			if id != targetID {
				return nil, errors.New("not found")
			}
			return &farmacia.Farmacia{ID: id, Nombre: "Farmacia Encontrada"}, nil
		},
	}
	svc := NuevoServicioFarmacia(repo)

	res, err := svc.BuscarPorID(targetID)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	if res.Nombre != "Farmacia Encontrada" {
		t.Errorf("Se esperaba la farmacia con nombre 'Farmacia Encontrada', llegó: %s", res.Nombre)
	}
}

func TestBuscarCercanas_RadioDefault(t *testing.T) {
	var radioUtilizado float64
	repo := &mockFarmaciaRepo{
		onCercanas: func(lat, lon, radio float64) ([]farmacia.Farmacia, error) {
			radioUtilizado = radio
			return []farmacia.Farmacia{}, nil
		},
	}
	svc := NuevoServicioFarmacia(repo)

	// Probamos pasando un radio <= 0 para verificar que use el fallback de 5.0
	_, _ = svc.BuscarCercanas(-12.04, -77.03, 0)

	if radioUtilizado != 5.0 {
		t.Errorf("Se esperaba un radio asignado por defecto de 5.0, se envió: %f", radioUtilizado)
	}
}

func TestActualizarFarmacia(t *testing.T) {
	repo := &mockFarmaciaRepo{
		onAct: func(id string, f *farmacia.Farmacia) error {
			return nil
		},
	}
	svc := NuevoServicioFarmacia(repo)

	err := svc.ActualizarFarmacia("id-123", &farmacia.Farmacia{Nombre: "Nombre Nuevo"})
	if err != nil {
		t.Errorf("No se esperaba error al actualizar: %v", err)
	}
}

func TestEliminarFarmacia(t *testing.T) {
	repo := &mockFarmaciaRepo{
		onEliminar: func(id string) error {
			return nil
		},
	}
	svc := NuevoServicioFarmacia(repo)

	err := svc.EliminarFarmacia("id-123")
	if err != nil {
		t.Errorf("No se esperaba error al eliminar: %v", err)
	}
}