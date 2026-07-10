package monitoreo

import (
	"errors"
	"proyecto-medicare-adulto-mayor/internal/models/monitoreo"
	"proyecto-medicare-adulto-mayor/internal/storage"
)

type ServicioMonitoreo struct {
	Repo storage.RepositorioMonitoreo
}

func NuevoServicioMonitoreo(repo storage.RepositorioMonitoreo) *ServicioMonitoreo {
	return &ServicioMonitoreo{Repo: repo}
}

func (s *ServicioMonitoreo) ListarRelaciones() ([]monitoreo.CuidadorPaciente, error) {
	return s.Repo.ListarRelaciones()
}

func (s *ServicioMonitoreo) BuscarRelacionPorID(id int) (monitoreo.CuidadorPaciente, error) {
	if id <= 0 {
		return monitoreo.CuidadorPaciente{}, errors.New("el ID debe ser mayor a 0")
	}
	return s.Repo.BuscarRelacionPorID(id)
}

func (s *ServicioMonitoreo) CrearRelacion(rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error) {
	if rel.CuidadorID <= 0 || rel.PacienteID <= 0 {
		return monitoreo.CuidadorPaciente{}, errors.New("el id del cuidador y del paciente son obligatorios")
	}
	if rel.Relacion == "" {
		return monitoreo.CuidadorPaciente{}, errors.New("la relación (ej. Hijo, Enfermero) es obligatoria")
	}
	return s.Repo.CrearRelacion(rel)
}

func (s *ServicioMonitoreo) ActualizarRelacion(id int, rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error) {
	if id <= 0 {
		return monitoreo.CuidadorPaciente{}, errors.New("el ID para actualizar es inválido")
	}
	if rel.CuidadorID <= 0 || rel.PacienteID <= 0 || rel.Relacion == "" {
		return monitoreo.CuidadorPaciente{}, errors.New("los datos de actualización no pueden estar vacíos")
	}
	return s.Repo.ActualizarRelacion(id, rel)
}

func (s *ServicioMonitoreo) EliminarRelacion(id int) (bool, error) {
	if id <= 0 {
		return false, errors.New("el ID debe ser mayor a 0")
	}
	return s.Repo.EliminarRelacion(id)
}