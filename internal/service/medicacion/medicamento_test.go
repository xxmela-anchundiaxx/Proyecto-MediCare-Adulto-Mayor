package medicacion

import (
	"medicare-adulto-mayor/internal/models/medicacion"
	"testing"
)

// MockRepositorioMedicamento implements RepositorioMedicamento for testing business rules
type MockRepositorioMedicamento struct {
	CrearMedicamentoCalled bool
	BuscarPorIDCalled      bool
	ListarPorPacienteCalled bool
	ActualizarStockCalled   bool
}

func (m *MockRepositorioMedicamento) BuscarPorID(id string) (*medicacion.Medicamento, error) {
	m.BuscarPorIDCalled = true
	return nil, nil
}

func (m *MockRepositorioMedicamento) ListarPorPaciente(pacienteID string) ([]medicacion.Medicamento, error) {
	m.ListarPorPacienteCalled = true
	return nil, nil
}

func (m *MockRepositorioMedicamento) CrearMedicamento(med *medicacion.Medicamento) error {
	m.CrearMedicamentoCalled = true
	return nil
}

func (m *MockRepositorioMedicamento) ActualizarStock(id string, nuevoStock int) error {
	m.ActualizarStockCalled = true
	return nil
}

func TestRegistrarMedicamento_ReglaDeNegocio_DatosInvalidos(t *testing.T) {
	// Arrange
	mockRepo := &MockRepositorioMedicamento{}
	servicio := NuevoServicioMedicamento(mockRepo)

	// Caso inválido: Nombre vacío (regla de negocio debe rechazar y no llamar al repositorio)
	reqInvalido := medicacion.CreateMedicamentoRequest{
		PacienteID:  "paciente-123",
		Nombre:      "", // Inválido!
		Dosis:       "500mg",
		Frecuencia:  "Cada 8 horas",
		Stock:       10,
	}

	// Act
	res, err := servicio.RegistrarMedicamiento(reqInvalido)

	// Assert
	if err == nil {
		t.Fatal("se esperaba un error de validación, pero no ocurrió")
	}

	if err.Error() != "campos obligatorios incompletos" {
		t.Errorf("se esperaba el error 'campos obligatorios incompletos', pero se obtuvo: %v", err)
	}

	if res != nil {
		t.Errorf("el resultado de medicamento debería ser nil, pero se obtuvo: %+v", res)
	}

	// Verificar regla de negocio: NO se debió llamar al método del repositorio para persistir
	if mockRepo.CrearMedicamentoCalled {
		t.Error("¡Violación de regla de negocio! Se llamó a CrearMedicamento del repositorio con datos inválidos")
	}
}

func TestRegistrarMedicamento_DatosValidos(t *testing.T) {
	// Arrange
	mockRepo := &MockRepositorioMedicamento{}
	servicio := NuevoServicioMedicamento(mockRepo)

	reqValido := medicacion.CreateMedicamentoRequest{
		PacienteID:  "paciente-123",
		Nombre:      "Paracetamol",
		Dosis:       "500mg",
		Frecuencia:  "Cada 8 horas",
		Stock:       10,
	}

	// Act
	res, err := servicio.RegistrarMedicamiento(reqValido)

	// Assert
	if err != nil {
		t.Fatalf("no se esperaba un error, pero ocurrió: %v", err)
	}

	if res == nil {
		t.Fatal("se esperaba un medicamento registrado, pero se obtuvo nil")
	}

	if res.Nombre != "Paracetamol" {
		t.Errorf("nombre incorrecto, se obtuvo: %s", res.Nombre)
	}

	// Se debió llamar al repositorio
	if !mockRepo.CrearMedicamentoCalled {
		t.Error("se esperaba que se llamara a CrearMedicamento en el repositorio")
	}
}
