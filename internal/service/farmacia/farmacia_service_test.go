package farmacia

import (
	"testing"

	model "medicare-adulto-mayor/internal/models/farmacia"
)

// Mock del repositorio
type MockRepositorio struct{}

func (m *MockRepositorio) CrearFarmacia(f *model.Farmacia) error {
	return nil
}

func (m *MockRepositorio) ListarTodas() ([]model.Farmacia, error) {
	return []model.Farmacia{}, nil
}

func (m *MockRepositorio) BuscarCercanas(lat, lon, radio float64) ([]model.Farmacia, error) {
	return []model.Farmacia{}, nil
}

// ===========================
// TEST 1: Registro correcto
// ===========================

func TestRegistrarFarmaciaCorrecta(t *testing.T) {

	repo := &MockRepositorio{}
	servicio := NuevoServicioFarmacia(repo)

	f := &model.Farmacia{
		Nombre:    "Farmacia Cruz Azul",
		Direccion: "Av. 4 de Noviembre",
	}

	err := servicio.RegistrarFarmacia(f)

	if err != nil {
		t.Fatalf("No debía devolver error, pero devolvió: %v", err)
	}

	if f.ID == "" {
		t.Error("Se esperaba que se genere un ID automáticamente")
	}
}

// ===========================
// TEST 2: Nombre vacío
// ===========================

func TestRegistrarFarmaciaSinNombre(t *testing.T) {

	repo := &MockRepositorio{}
	servicio := NuevoServicioFarmacia(repo)

	f := &model.Farmacia{
		Nombre:    "",
		Direccion: "Manta",
	}

	err := servicio.RegistrarFarmacia(f)

	if err == nil {
		t.Error("Se esperaba un error porque el nombre está vacío")
	}
}

// ===========================
// TEST 3: Dirección vacía
// ===========================

func TestRegistrarFarmaciaSinDireccion(t *testing.T) {

	repo := &MockRepositorio{}
	servicio := NuevoServicioFarmacia(repo)

	f := &model.Farmacia{
		Nombre:    "Farmacia Cruz Azul",
		Direccion: "",
	}

	err := servicio.RegistrarFarmacia(f)

	if err == nil {
		t.Error("Se esperaba un error porque la dirección está vacía")
	}
}

// ===========================
// TEST 4: Buscar cercanas
// ===========================

func TestBuscarCercanas(t *testing.T) {

	repo := &MockRepositorio{}
	servicio := NuevoServicioFarmacia(repo)

	_, err := servicio.BuscarCercanas(-0.95, -80.73, 0)

	if err != nil {
		t.Errorf("No debía devolver error: %v", err)
	}
}

// ===========================
// TEST 5: Listar farmacias
// ===========================

func TestListarTodas(t *testing.T) {

	repo := &MockRepositorio{}
	servicio := NuevoServicioFarmacia(repo)

	farmacias, err := servicio.ListarTodas()

	if err != nil {
		t.Fatalf("No debía devolver error")
	}

	if farmacias == nil {
		t.Error("Se esperaba una lista vacía, no nil")
	}
}
