package medicacion

import (
	"errors"
	"medicare-adulto-mayor/internal/models/medicacion"
	medicacionStorage "medicare-adulto-mayor/internal/storage/medicacion"
	"github.com/google/uuid"
)

type ServicioPaciente struct {
	Repo medicacionStorage.RepositorioPaciente
}

func NuevoServicioPaciente(repo medicacionStorage.RepositorioPaciente) *ServicioPaciente {
	return &ServicioPaciente{Repo: repo}
}

func (s *ServicioPaciente) CrearPaciente(req medicacion.CreatePacienteRequest) (*medicacion.Paciente, error) {
	if req.UsuarioID == "" || req.ContactoEmergencia == "" {
		return nil, errors.New("datos de paciente obligatorios incompletos")
	}

	existente, err := s.Repo.BuscarPorUsuarioID(req.UsuarioID)
	if err != nil {
		return nil, err
	}
	if existente != nil {
		return existente, nil // Retorna el paciente existente si ya está registrado
	}

	p := &medicacion.Paciente{
		ID:                 uuid.New().String(),
		UsuarioID:          req.UsuarioID,
		CuidadorID:         req.CuidadorID,
		Edad:               req.Edad,
		GrupoSanguineo:     req.GrupoSanguineo,
		Alergias:           req.Alergias,
		CondicionesMedicas: req.CondicionesMedicas,
		ContactoEmergencia: req.ContactoEmergencia,
	}

	if err := s.Repo.CrearPaciente(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ServicioPaciente) ObtenerPaciente(id string) (*medicacion.Paciente, error) {
	p, err := s.Repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("paciente no encontrado")
	}
	return p, nil
}

func (s *ServicioPaciente) ListarPorCuidador(cuidadorID string) ([]medicacion.Paciente, error) {
	return s.Repo.ListarPorCuidador(cuidadorID)
}
