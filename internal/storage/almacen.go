package storage

import (
	"proyecto-medicare-adulto-mayor/internal/models/farmacia"
	"proyecto-medicare-adulto-mayor/internal/models/medicacion"
	"proyecto-medicare-adulto-mayor/internal/models/monitoreo"
	
	// SOLUCIÓN: Importamos el paquete de usuarios/auth asignándole el alias "models"
	models "proyecto-medicare-adulto-mayor/internal/models" 
)

// Interfaces del módulo Medicación
type MedicacionRepository interface {
	ListarMedicacion() ([]medicacion.Medicacion, error)
	BuscarMedicacionPorID(id int) (medicacion.Medicacion, error)
	CrearMedicacion(m medicacion.Medicacion) (medicacion.Medicacion, error)
	ActualizarMedicacion(id int, m medicacion.Medicacion) (medicacion.Medicacion, error)
	EliminarMedicacion(id int) error
}

// Interfaces de Pacientes
type PacientesRepository interface {
	ListarPacientes() ([]medicacion.Paciente, error)
	BuscarPacientePorID(id int) (medicacion.Paciente, error)
	CrearPaciente(p medicacion.Paciente) (medicacion.Paciente, error)
	ActualizarPaciente(id int, p medicacion.Paciente) (medicacion.Paciente, error)
	EliminarPaciente(id int) error
}

// Interfaces de Historial
type HistorialRepository interface {
	BuscarHistorialPorID(id int) (medicacion.HistorialMedicacion, error)
	ListarHistorial() ([]medicacion.HistorialMedicacion, error)
	CrearHistorial(h medicacion.HistorialMedicacion) (medicacion.HistorialMedicacion, error)
	ActualizarHistorial(id int, h medicacion.HistorialMedicacion) (medicacion.HistorialMedicacion, error)
	EliminarHistorial(id int) error
}

// Interfaz para consultas cruzadas de medicación/historial por paciente
type MedicacionHistorialRepository interface {
	ListarMedicacionPorPaciente(pacienteID int) ([]medicacion.Medicacion, error)
	ListarHistorialPorPaciente(pacienteID int) ([]medicacion.HistorialMedicacion, error)
}

// Interfaces del módulo Farmacia
type RepositorioFarmacia interface {
	ListarTodas() ([]farmacia.Farmacia, error)
	BuscarPorID(id string) (*farmacia.Farmacia, error)
	CrearFarmacia(f *farmacia.Farmacia) error
	ActualizarFarmacia(id string, f *farmacia.Farmacia) error
	EliminarFarmacia(id string) error
	BuscarCercanas(lat, lon, radioKM float64) ([]farmacia.Farmacia, error)
}

// Interfaces del módulo Monitoreo
type RepositorioMonitoreo interface {
	ListarRelaciones() ([]monitoreo.CuidadorPaciente, error)
	BuscarRelacionPorID(id int) (monitoreo.CuidadorPaciente, error)
	CrearRelacion(rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error)
	ActualizarRelacion(id int, rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error)
	EliminarRelacion(id int) (bool, error)
}


type UserRepository interface {
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
}

// Interfaz compuesta final
type Almacen interface {
	MedicacionRepository
	PacientesRepository
	HistorialRepository
	MedicacionHistorialRepository
	RepositorioFarmacia
	RepositorioMonitoreo
	UserRepository
}