package service

import (
    "proyecto-medicare-adulto-mayor/internal/models/monitoreo"
    "proyecto-medicare-adulto-mayor/internal/storage"
)

type MonitoreoService struct {
    repo storage.MonitoreoRepository
}

func NewMonitoreoService(repo storage.MonitoreoRepository) *MonitoreoService {
    return &MonitoreoService{repo: repo}
}

func (s *MonitoreoService) Listar() ([]monitoreo.CuidadorPaciente, error) {
    return s.repo.ListarRelaciones()
}

func (s *MonitoreoService) Obtener(id int) (monitoreo.CuidadorPaciente, error) {
    return s.repo.BuscarRelacionPorID(id)
}

func (s *MonitoreoService) Crear(rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error) {
    return s.repo.CrearRelacion(rel)
}

func (s *MonitoreoService) Actualizar(id int, rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error) {
    return s.repo.ActualizarRelacion(id, rel)
}

func (s *MonitoreoService) Eliminar(id int) (bool, error) {
    return s.repo.EliminarRelacion(id)
}