package medicacion

import (
	"errors"
	"medicare-adulto-mayor/internal/models/medicacion"
	medicacionStorage "medicare-adulto-mayor/internal/storage/medicacion"
	"time"

	"github.com/google/uuid"
)

type ServicioMedicamento struct {
	Repo medicacionStorage.RepositorioMedicamento
}

func NuevoServicioMedicamento(repo medicacionStorage.RepositorioMedicamento) *ServicioMedicamento {
	return &ServicioMedicamento{Repo: repo}
}

func (s *ServicioMedicamento) RegistrarMedicamiento(req medicacion.CreateMedicamentoRequest) (*medicacion.Medicamento, error) {
	if req.PacienteID == "" || req.Nombre == "" || req.Dosis == "" || req.Frecuencia == "" {
		return nil, errors.New("campos obligatorios incompletos")
	}

	m := &medicacion.Medicamento{
		ID:                uuid.New().String(),
		PacienteID:        req.PacienteID,
		Nombre:            req.Nombre,
		Descripcion:       req.Descripcion,
		Dosis:             req.Dosis,
		Frecuencia:        req.Frecuencia,
		ViaAdministracion: req.ViaAdministracion,
		Stock:             req.Stock,
		FechaRegistro:     time.Now(),
	}

	if err := s.Repo.CrearMedicamento(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *ServicioMedicamento) ListarPorPaciente(pacienteID string) ([]medicacion.Medicamento, error) {
	return s.Repo.ListarPorPaciente(pacienteID)
}

func (s *ServicioMedicamento) ObtenerPorID(id string) (*medicacion.Medicamento, error) {
	m, err := s.Repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("medicamento no encontrado")
	}
	return m, nil
}

func (s *ServicioMedicamento) ActualizarStock(id string, cantidadConsumida int) error {
	m, err := s.Repo.BuscarPorID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("medicamento no encontrado")
	}

	nuevoStock := m.Stock - cantidadConsumida
	if nuevoStock < 0 {
		nuevoStock = 0 // Evita stock negativo
	}

	return s.Repo.ActualizarStock(id, nuevoStock)
}
