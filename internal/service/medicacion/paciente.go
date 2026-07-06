package medicacion

import (
	models "proyecto-medicare-adulto-mayor/internal/models/medicacion"
	"proyecto-medicare-adulto-mayor/internal/service"
	"proyecto-medicare-adulto-mayor/internal/storage"
)

// paciente
type PacienteService struct {
	repo storage.PacientesRepository
}

func NewPacienteService(repo storage.PacientesRepository) *PacienteService {
	return &PacienteService{repo: repo}
}

func (s *PacienteService) Listar() ([]models.Paciente, error) {
	return s.repo.ListarPacientes()
}

func (s *PacienteService) Obtener(id int) (models.Paciente, error) {
	paciente, err := s.repo.BuscarPacientePorID(id)
	if err != nil {
		return models.Paciente{}, service.ErrNoEncontrado
	}
	return paciente, nil
}

func (s *PacienteService) Crear(p models.Paciente) (models.Paciente, error) {
	return s.repo.CrearPaciente(p)
}

func (s *PacienteService) Actualizar(id int, p models.Paciente) (models.Paciente, error) {
	return s.repo.ActualizarPaciente(id, p)
}

func (s *PacienteService) Eliminar(id int) error {
	return s.repo.EliminarPaciente(id)
}