package monitoreo

import (
	"errors"
	"medicare-adulto-mayor/internal/models/monitoreo"
	monitoreoStorage "medicare-adulto-mayor/internal/storage/monitoreo"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ServicioMonitoreo struct {
	Repo monitoreoStorage.RepositorioMonitoreo
}

func NuevoServicioMonitoreo(repo monitoreoStorage.RepositorioMonitoreo) *ServicioMonitoreo {
	return &ServicioMonitoreo{Repo: repo}
}

func (s *ServicioMonitoreo) RegistrarMedicion(req monitoreo.RegistrarSignosRequest) (*monitoreo.MonitoreoSignos, error) {
	if req.PacienteID == "" {
		return nil, errors.New("paciente_id es obligatorio")
	}

	alerta := s.analizarUmbralesCriticos(req)

	m := &monitoreo.MonitoreoSignos{
		ID:             uuid.New().String(),
		PacienteID:     req.PacienteID,
		RitmoCardiaco:  req.RitmoCardiaco,
		PresionArterial: req.PresionArterial,
		NivelAzucar:    req.NivelAzucar,
		Temperatura:    req.Temperatura,
		FechaHora:      time.Now(),
		AlertaEnviada:  alerta,
		Observaciones:  req.Observaciones,
	}

	if alerta && m.Observaciones == "" {
		m.Observaciones = "¡ALERTA AUTOMÁTICA! Signos fuera del rango normal de salud."
	}

	if err := s.Repo.RegistrarSignos(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *ServicioMonitoreo) ListarPorPaciente(pacienteID string) ([]monitoreo.MonitoreoSignos, error) {
	return s.Repo.ListarPorPaciente(pacienteID)
}

func (s *ServicioMonitoreo) analizarUmbralesCriticos(req monitoreo.RegistrarSignosRequest) bool {
	// Analizar Ritmo Cardíaco (Normal: 60 - 100 ppm para adulto mayor en reposo)
	if req.RitmoCardiaco > 0 && (req.RitmoCardiaco < 50 || req.RitmoCardiaco > 105) {
		return true
	}

	// Analizar Temperatura (Fiebre > 38.0°C, Hipotermia < 35.5°C)
	if req.Temperatura > 0 && (req.Temperatura < 35.5 || req.Temperatura > 37.8) {
		return true
	}

	// Analizar Nivel de Azúcar (Normal en ayunas: 70 - 130 mg/dL, postprandial < 180 mg/dL)
	if req.NivelAzucar > 0 && (req.NivelAzucar < 70 || req.NivelAzucar > 150) {
		return true
	}

	// Analizar Presión Arterial (ej. "140/90")
	if req.PresionArterial != "" {
		partes := strings.Split(req.PresionArterial, "/")
		if len(partes) == 2 {
			sistolica, err1 := strconv.Atoi(strings.TrimSpace(partes[0]))
			diastolica, err2 := strconv.Atoi(strings.TrimSpace(partes[1]))
			if err1 == nil && err2 == nil {
				// Criterio de hipertensión sistólica > 140 o diastólica > 90, o hipotensión sistólica < 90
				if sistolica < 90 || sistolica > 140 || diastolica > 90 {
					return true
				}
			}
		}
	}

	return false
}
