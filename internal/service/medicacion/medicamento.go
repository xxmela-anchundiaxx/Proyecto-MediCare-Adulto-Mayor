package medicacion

import (
    "proyecto-medicare-adulto-mayor/internal/models/medicacion"
    "proyecto-medicare-adulto-mayor/internal/storage"
    "proyecto-medicare-adulto-mayor/internal/service"
)

type MedicacionService struct {
    repo storage.MedicacionRepository
}

func NewMedicacionService(repo storage.MedicacionRepository) *MedicacionService {
    return &MedicacionService{repo: repo}
}

func (s *MedicacionService) Listar() ([]medicacion.Medicacion, error) {
    return s.repo.ListarMedicacion()
}

func (s *MedicacionService) Obtener(id int) (medicacion.Medicacion, error) {
    med, err := s.repo.BuscarMedicacionPorID(id)
    if err != nil {
        return medicacion.Medicacion{}, service.ErrNoEncontrado
    }
    return med, nil
}

func (s *MedicacionService) Crear(m medicacion.Medicacion) (medicacion.Medicacion, error) {
    return s.repo.CrearMedicacion(m)
}

func (s *MedicacionService) Actualizar(id int, m medicacion.Medicacion) (medicacion.Medicacion, error) {
    return s.repo.ActualizarMedicacion(id, m)
}

func (s *MedicacionService) Eliminar(id int) error {
    return s.repo.EliminarMedicacion(id)
}
