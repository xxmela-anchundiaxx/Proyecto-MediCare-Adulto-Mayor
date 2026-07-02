package monitoreo

import (
	"medicare-adulto-mayor/internal/models/monitoreo"
	"gorm.io/gorm"
)

type RepositorioMonitoreo interface {
	RegistrarSignos(s *monitoreo.MonitoreoSignos) error
	ListarPorPaciente(pacienteID string) ([]monitoreo.MonitoreoSignos, error)
}

type StorageMonitoreoGORM struct {
	DB *gorm.DB
}

func NuevoStorageMonitoreoGORM(db *gorm.DB) *StorageMonitoreoGORM {
	return &StorageMonitoreoGORM{DB: db}
}

func (s *StorageMonitoreoGORM) RegistrarSignos(m *monitoreo.MonitoreoSignos) error {
	return s.DB.Create(m).Error
}

func (s *StorageMonitoreoGORM) ListarPorPaciente(pacienteID string) ([]monitoreo.MonitoreoSignos, error) {
	var lista []monitoreo.MonitoreoSignos
	err := s.DB.Where("paciente_id = ?", pacienteID).Order("fecha_hora DESC").Find(&lista).Error
	return lista, err
}
