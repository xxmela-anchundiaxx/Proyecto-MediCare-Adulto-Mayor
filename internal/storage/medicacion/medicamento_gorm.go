package medicacion

import (
	"medicare-adulto-mayor/internal/models/medicacion"
	"gorm.io/gorm"
)

type RepositorioMedicamento interface {
	BuscarPorID(id string) (*medicacion.Medicamento, error)
	ListarPorPaciente(pacienteID string) ([]medicacion.Medicamento, error)
	CrearMedicamento(m *medicacion.Medicamento) error
	ActualizarStock(id string, nuevoStock int) error
}

type StorageMedicamentoGORM struct {
	DB *gorm.DB
}

func NuevoStorageMedicamentoGORM(db *gorm.DB) *StorageMedicamentoGORM {
	return &StorageMedicamentoGORM{DB: db}
}

func (s *StorageMedicamentoGORM) BuscarPorID(id string) (*medicacion.Medicamento, error) {
	var m medicacion.Medicamento
	err := s.DB.First(&m, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (s *StorageMedicamentoGORM) ListarPorPaciente(pacienteID string) ([]medicacion.Medicamento, error) {
	var lista []medicacion.Medicamento
	err := s.DB.Where("paciente_id = ?", pacienteID).Find(&lista).Error
	return lista, err
}

func (s *StorageMedicamentoGORM) CrearMedicamento(m *medicacion.Medicamento) error {
	return s.DB.Create(m).Error
}

func (s *StorageMedicamentoGORM) ActualizarStock(id string, nuevoStock int) error {
	return s.DB.Model(&medicacion.Medicamento{}).Where("id = ?", id).Update("stock", nuevoStock).Error
}
