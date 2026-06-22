package storage

import (
	"proyecto-medicare-adulto-mayor/internal/models"
)

type Almacen interface {
	ListarMedicacion() ([]models.Medicacion, error)
	BuscarMedicacionPorID(id int) (models.Medicacion, error)
	CrearMedicacion(medicacion models.Medicacion) (models.Medicacion, error)
	ActualizarMedicacion(id int, medicacion models.Medicacion) (models.Medicacion, error)
	EliminarMedicacion(id int) (bool, error)

    ListarPacientes() ([]models.Paciente, error)
    BuscarPacientePorID(id int) (models.Paciente, error)
    CrearPaciente(p models.Paciente) (models.Paciente, error)
    ActualizarPaciente(id int, p models.Paciente) (models.Paciente, error)
    EliminarPaciente(id int) (bool, error)
	
	BuscarHistorialPorID(id int) (models.HistorialMedicacion, error)
	ListarHistorial() ([]models.HistorialMedicacion, error)
	CrearHistorial(h models.HistorialMedicacion) (models.HistorialMedicacion, error)
	ActualizarHistorial(id int, h models.HistorialMedicacion) (models.HistorialMedicacion, error)
	EliminarHistorial(id int) (bool, error)

	ListarMedicacionPorPaciente(pacienteID int) ([]models.Medicacion, error)
	ListarHistorialPorPaciente(pacienteID int) ([]models.HistorialMedicacion, error)
}




