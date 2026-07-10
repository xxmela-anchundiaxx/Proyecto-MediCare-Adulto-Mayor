package monitoreo

import (
	"testing"
	"proyecto-medicare-adulto-mayor/internal/models/monitoreo"
)

// mockMonitoreoRepo implementa la interfaz storage.RepositorioMonitoreo
type mockMonitoreoRepo struct {
	onListar     func() ([]monitoreo.CuidadorPaciente, error)
	onBuscarID   func(id int) (monitoreo.CuidadorPaciente, error)
	onCrear      func(rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error)
	onActualizar func(id int, rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error)
	onEliminar   func(id int) (bool, error)
}

func (m *mockMonitoreoRepo) ListarRelaciones() ([]monitoreo.CuidadorPaciente, error) { 
	return m.onListar() 
}
func (m *mockMonitoreoRepo) BuscarRelacionPorID(id int) (monitoreo.CuidadorPaciente, error) { 
	return m.onBuscarID(id) 
}
func (m *mockMonitoreoRepo) CrearRelacion(rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error) { 
	return m.onCrear(rel) 
}
func (m *mockMonitoreoRepo) ActualizarRelacion(id int, rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error) { 
	return m.onActualizar(id, rel) 
}
func (m *mockMonitoreoRepo) EliminarRelacion(id int) (bool, error) { 
	return m.onEliminar(id) 
}

// --- PRUEBAS UNITARIAS ---

func TestCrearRelacion_Exitoso(t *testing.T) {
	repo := &mockMonitoreoRepo{
		onCrear: func(rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error) {
			rel.ID = 1
			return rel, nil
		},
	}
	
	// Ya no lleva el prefijo "servicio." porque estás en el mismo package
	svc := NuevoServicioMonitoreo(repo)

	nuevaRel := monitoreo.CuidadorPaciente{
		CuidadorID: 2,
		PacienteID: 3,
		Relacion:   "Hijo",
	}

	resultado, err := svc.CrearRelacion(nuevaRel)
	if err != nil {
		t.Fatalf("No se esperaba error en la creación, se obtuvo: %v", err)
	}

	if resultado.ID != 1 {
		t.Errorf("Se esperaba que el ID devuelto fuera 1, se obtuvo: %d", resultado.ID)
	}
}

func TestCrearRelacion_ErroresValidacion(t *testing.T) {
	svc := NuevoServicioMonitoreo(&mockMonitoreoRepo{})

	tests := []struct {
		name string
		rel  monitoreo.CuidadorPaciente
	}{
		{"CuidadorID vacío", monitoreo.CuidadorPaciente{CuidadorID: 0, PacienteID: 1, Relacion: "Enfermero"}},
		{"PacienteID vacío", monitoreo.CuidadorPaciente{CuidadorID: 2, PacienteID: 0, Relacion: "Familiar"}},
		{"Relación vacía", monitoreo.CuidadorPaciente{CuidadorID: 2, PacienteID: 3, Relacion: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CrearRelacion(tt.rel)
			if err == nil {
				t.Error("Se esperaba un error de validación en el servicio, pero no ocurrió")
			}
		})
	}
}

func TestBuscarRelacionPorID_Invalido(t *testing.T) {
	svc := NuevoServicioMonitoreo(&mockMonitoreoRepo{})

	_, err := svc.BuscarRelacionPorID(0)
	if err == nil {
		t.Error("Se esperaba un error al buscar un ID menor o igual a cero")
	}
}

func TestActualizarRelacion_Exitoso(t *testing.T) {
	repo := &mockMonitoreoRepo{
		onActualizar: func(id int, rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error) {
			rel.ID = id
			return rel, nil
		},
	}
	svc := NuevoServicioMonitoreo(repo)

	relActualizada := monitoreo.CuidadorPaciente{CuidadorID: 4, PacienteID: 5, Relacion: "Tutor"}
	resultado, err := svc.ActualizarRelacion(10, relActualizada)

	if err != nil {
		t.Fatalf("Error inesperado al actualizar: %v", err)
	}

	if resultado.ID != 10 {
		t.Errorf("Se esperaba que mantuviera el ID 10, llegó: %d", resultado.ID)
	}
}

func TestEliminarRelacion_NotFound(t *testing.T) {
	repo := &mockMonitoreoRepo{
		onEliminar: func(id int) (bool, error) {
			return false, nil
		},
	}
	svc := NuevoServicioMonitoreo(repo)

	eliminado, err := svc.EliminarRelacion(99)
	if err != nil {
		t.Fatalf("No se esperaba un error del sistema, se obtuvo: %v", err)
	}

	if eliminado {
		t.Error("Se esperaba que 'eliminado' fuera false para una relación inexistente")
	}
}