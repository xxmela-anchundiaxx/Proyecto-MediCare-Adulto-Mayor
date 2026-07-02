package medicacion

import (
	"medicare-adulto-mayor/internal/models/medicacion"
	"gorm.io/gorm"
)

type RepositorioPaciente interface {
	BuscarPorID(id string) (*medicacion.Paciente, error)
	BuscarPorUsuarioID(usuarioID string) (*medicacion.Paciente, error)
	CrearPaciente(p *medicacion.Paciente) error
	ListarPorCuidador(cuidadorID string) ([]medicacion.Paciente, error)
}

type StoragePacienteGORM struct {
	DB *gorm.DB
}

func NuevoStoragePacienteGORM(db *gorm.DB) *StoragePacienteGORM {
	return &StoragePacienteGORM{DB: db}
}

func (s *StoragePacienteGORM) BuscarPorID(id string) (*medicacion.Paciente, error) {
	var p medicacion.Paciente
	err := s.DB.First(&p, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (s *StoragePacienteGORM) BuscarPorUsuarioID(usuarioID string) (*medicacion.Paciente, error) {
	var p medicacion.Paciente
	err := s.DB.First(&p, "usuario_id = ?", usuarioID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (s *StoragePacienteGORM) CrearPaciente(p *medicacion.Paciente) error {
	return s.DB.Create(p).Error
}

func (s *StoragePacienteGORM) ListarPorCuidador(cuidadorID string) ([]medicacion.Paciente, error) {
	var lista []medicacion.Paciente
	err := s.DB.Where("cuidador_id = ?", cuidadorID).Find(&lista).Error
	return lista, err
}
