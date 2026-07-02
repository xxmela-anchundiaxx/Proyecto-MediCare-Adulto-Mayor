package medicacion

import (
	"errors"
	"medicare-adulto-mayor/internal/models/medicacion"
	medicacionStorage "medicare-adulto-mayor/internal/storage/medicacion"
	"time"

	"github.com/google/uuid"
)

type ServicioHistorial struct {
	RepoHistorial   medicacionStorage.RepositorioHistorial
	RepoMedicamento medicacionStorage.RepositorioMedicamento
}

func NuevoServicioHistorial(rh medicacionStorage.RepositorioHistorial, rm medicacionStorage.RepositorioMedicamento) *ServicioHistorial {
	return &ServicioHistorial{
		RepoHistorial:   rh,
		RepoMedicamento: rm,
	}
}

func (s *ServicioHistorial) RegistrarToma(req medicacion.RecordAdherenceRequest) (*medicacion.HistorialMedicacion, error) {
	if req.MedicamentoID == "" || req.PacienteID == "" {
		return nil, errors.New("medicamento_id y paciente_id son obligatorios")
	}

	// 1. Obtener medicamento para validar existencia y descontar stock
	m, err := s.RepoMedicamento.BuscarPorID(req.MedicamentoID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("medicamento no encontrado")
	}

	h := &medicacion.HistorialMedicacion{
		ID:            uuid.New().String(),
		MedicamentoID: req.MedicamentoID,
		PacienteID:    req.PacienteID,
		FechaHora:     time.Now(),
		Tomado:        req.Tomado,
		Observaciones: req.Observaciones,
	}

	// 2. Registrar la toma en el historial
	if err := s.RepoHistorial.CrearHistorial(h); err != nil {
		return nil, err
	}

	// 3. Si fue tomado, descontar 1 del stock de medicamentos
	if req.Tomado && m.Stock > 0 {
		if err := s.RepoMedicamento.ActualizarStock(m.ID, m.Stock-1); err != nil {
			// No bloqueamos la respuesta, pero idealmente se maneja el error de consistencia
		}
	}

	return h, nil
}

func (s *ServicioHistorial) ListarPorPaciente(pacienteID string) ([]medicacion.HistorialMedicacion, error) {
	return s.RepoHistorial.ListarPorPaciente(pacienteID)
}
