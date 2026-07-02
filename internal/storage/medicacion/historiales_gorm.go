package medicacion

import (
	"medicare-adulto-mayor/internal/models/medicacion"
	"gorm.io/gorm"
)

type RepositorioHistorial interface {
	CrearHistorial(h *medicacion.HistorialMedicacion) error
	ListarPorPaciente(pacienteID string) ([]medicacion.HistorialMedicacion, error)
}

type StorageHistorialGORM struct {
	DB *gorm.DB
}

func NuevoStorageHistorialGORM(db *gorm.DB) *StorageHistorialGORM {
	return &StorageHistorialGORM{DB: db}
}

func (s *StorageHistorialGORM) CrearHistorial(h *medicacion.HistorialMedicacion) error {
	return s.DB.Create(h).Error
}

func (s *StorageHistorialGORM) ListarPorPaciente(pacienteID string) ([]medicacion.HistorialMedicacion, error) {
	var lista []medicacion.HistorialMedicacion
	err := s.DB.Where("paciente_id = ?", pacienteID).Order("fecha_hora DESC").Find(&lista).Error
	return lista, err
}
