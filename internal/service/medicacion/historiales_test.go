package medicacion

import (
	"errors"
	"testing"

	models "proyecto-medicare-adulto-mayor/internal/models/medicacion"
	"proyecto-medicare-adulto-mayor/internal/service"
)

// --- MOCKS DE REPOSITORIO (STORAGE) ---

type mockHistorialRepository struct {
	onBuscarHistorialPorID func(id int) (models.HistorialMedicacion, error)
	onListarHistorial      func() ([]models.HistorialMedicacion, error)
	onCrearHistorial       func(h models.HistorialMedicacion) (models.HistorialMedicacion, error)
	onActualizarHistorial  func(id int, h models.HistorialMedicacion) (models.HistorialMedicacion, error)
	onEliminarHistorial    func(id int) error
}

func (m *mockHistorialRepository) BuscarHistorialPorID(id int) (models.HistorialMedicacion, error) {
	return m.onBuscarHistorialPorID(id)
}
func (m *mockHistorialRepository) ListarHistorial() ([]models.HistorialMedicacion, error) {
	return m.onListarHistorial()
}
func (m *mockHistorialRepository) CrearHistorial(h models.HistorialMedicacion) (models.HistorialMedicacion, error) {
	return m.onCrearHistorial(h)
}
func (m *mockHistorialRepository) ActualizarHistorial(id int, h models.HistorialMedicacion) (models.HistorialMedicacion, error) {
	return m.onActualizarHistorial(id, h)
}
func (m *mockHistorialRepository) EliminarHistorial(id int) error {
	return m.onEliminarHistorial(id)
}

type mockMedicacionHistorialRepository struct {
	onListarMedicacionPorPaciente func(pacienteID int) ([]models.Medicacion, error)
	onListarHistorialPorPaciente  func(pacienteID int) ([]models.HistorialMedicacion, error)
}

func (m *mockMedicacionHistorialRepository) ListarMedicacionPorPaciente(pacienteID int) ([]models.Medicacion, error) {
	return m.onListarMedicacionPorPaciente(pacienteID)
}
func (m *mockMedicacionHistorialRepository) ListarHistorialPorPaciente(pacienteID int) ([]models.HistorialMedicacion, error) {
	return m.onListarHistorialPorPaciente(pacienteID)
}

// --- PRUEBAS UNITARIAS ---

func TestHistorialService_Obtener_Exitoso(t *testing.T) {
	mockRepo := &mockHistorialRepository{
		onBuscarHistorialPorID: func(id int) (models.HistorialMedicacion, error) {
			return models.HistorialMedicacion{
				ID:           id,
				MedicacionID: 10,
				Tomada:       true,
				Observacion:  "Tomado a tiempo",
			}, nil
		},
	}

	svc := NewHistorialService(mockRepo)
	resultado, err := svc.Obtener(1)

	if err != nil {
		t.Fatalf("No se esperaba error, se obtuvo: %v", err)
	}
	if resultado.MedicacionID != 10 {
		t.Errorf("Se esperaba MedicacionID 10, se obtuvo: %d", resultado.MedicacionID)
	}
	if resultado.Observacion != "Tomado a tiempo" {
		t.Errorf("La observación no coincide")
	}
}

func TestHistorialService_Obtener_NoEncontrado(t *testing.T) {
	mockRepo := &mockHistorialRepository{
		onBuscarHistorialPorID: func(id int) (models.HistorialMedicacion, error) {
			return models.HistorialMedicacion{}, errors.New("db error")
		},
	}

	svc := NewHistorialService(mockRepo)
	_, err := svc.Obtener(99)

	if !errors.Is(err, service.ErrNoEncontrado) {
		t.Errorf("Se esperaba el error global service.ErrNoEncontrado, se obtuvo: %v", err)
	}
}

func TestHistorialService_Listar(t *testing.T) {
	mockRepo := &mockHistorialRepository{
		onListarHistorial: func() ([]models.HistorialMedicacion, error) {
			return []models.HistorialMedicacion{
				{ID: 1}, {ID: 2},
			}, nil
		},
	}

	svc := NewHistorialService(mockRepo)
	lista, err := svc.Listar()

	if err != nil {
		t.Fatalf("Error inesperado al listar: %v", err)
	}
	if len(lista) != 2 {
		t.Errorf("Se esperaban 2 elementos, se obtuvieron: %d", len(lista))
	}
}