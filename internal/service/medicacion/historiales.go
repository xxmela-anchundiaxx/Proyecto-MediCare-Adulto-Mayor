package medicacion

import (
	models "proyecto-medicare-adulto-mayor/internal/models/medicacion"
	"proyecto-medicare-adulto-mayor/internal/service"
	"proyecto-medicare-adulto-mayor/internal/storage"
)

// historial
type HistorialService struct {
	repo storage.HistorialRepository
}

func NewHistorialService(repo storage.HistorialRepository) *HistorialService {
	return &HistorialService{repo: repo}
}

func (s *HistorialService) Listar() ([]models.HistorialMedicacion, error) {
	return s.repo.ListarHistorial()
}

func (s *HistorialService) Obtener(id int) (models.HistorialMedicacion, error) {
	historial, err := s.repo.BuscarHistorialPorID(id)
	if err != nil {
		return models.HistorialMedicacion{}, service.ErrNoEncontrado
	}
	return historial, nil
}

func (s *HistorialService) Crear(h models.HistorialMedicacion) (models.HistorialMedicacion, error) {
	return s.repo.CrearHistorial(h)
}

func (s *HistorialService) Actualizar(id int, h models.HistorialMedicacion) (models.HistorialMedicacion, error) {
	return s.repo.ActualizarHistorial(id, h)
}

func (s *HistorialService) Eliminar(id int) error {
	return s.repo.EliminarHistorial(id)
}

type MedicacionHistorialService struct {
	repo storage.Medicacion_HistorialRepository
}

func NewMedicacionHistorialService(repo storage.Medicacion_HistorialRepository) *MedicacionHistorialService {
	return &MedicacionHistorialService{repo: repo}
}

func (s *MedicacionHistorialService) ListarMedicacionPorPaciente(pacienteID int) ([]models.Medicacion, error) {
	return s.repo.ListarMedicacionPorPaciente(pacienteID)
}

func (s *MedicacionHistorialService) ListarHistorialPorPaciente(pacienteID int) ([]models.HistorialMedicacion, error) {
	return s.repo.ListarHistorialPorPaciente(pacienteID)
}