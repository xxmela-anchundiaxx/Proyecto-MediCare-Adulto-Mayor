package storage

import (
    "proyecto-medicare-adulto-mayor/internal/models/medicacion"
    //"proyecto-medicare-adulto-mayor/internal/models/farmacia"
    "proyecto-medicare-adulto-mayor/internal/models/monitoreo"
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
type Medicacion_HistorialRepository interface {
	ListarMedicacionPorPaciente(pacienteID int) ([]medicacion.Medicacion, error)
	ListarHistorialPorPaciente(pacienteID int) ([]medicacion.HistorialMedicacion, error)
}

/*type FarmaciaRepository interface {
    ListarFarmacias() ([]farmacia.Farmacia, error)
    BuscarFarmaciaPorID(id string) (farmacia.Farmacia, error)
    CrearFarmacia(f farmacia.Farmacia) (farmacia.Farmacia, error)
    ActualizarFarmacia(id string, f farmacia.Farmacia) (farmacia.Farmacia, error)
    EliminarFarmacia(id string) (bool, error)
}*/

type MonitoreoRepository interface {
    ListarRelaciones() ([]monitoreo.CuidadorPaciente, error)
    BuscarRelacionPorID(id int) (monitoreo.CuidadorPaciente, error)
    CrearRelacion(rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error)
    ActualizarRelacion(id int, rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error)
    EliminarRelacion(id int) (bool, error)
}

type MedicacionHistorialRepository interface {
    ListarMedicacionPorPaciente(pacienteID int) ([]medicacion.Medicacion, error)
    ListarHistorialPorPaciente(pacienteID int) ([]medicacion.HistorialMedicacion, error)
}

// Interfaz compuesta
type Almacen interface {
    MedicacionRepository
    PacientesRepository
    HistorialRepository
    MedicacionHistorialRepository
    //FarmaciaRepository
    //MonitoreoRepository
}

