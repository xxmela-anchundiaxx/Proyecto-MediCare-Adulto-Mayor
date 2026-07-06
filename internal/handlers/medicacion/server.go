package medicacion

import (
	servicioMedicacion "proyecto-medicare-adulto-mayor/internal/service/medicacion"
)

type Server struct {
	Medicacion          *servicioMedicacion.MedicacionService
	Paciente            *servicioMedicacion.PacienteService
	Historial           *servicioMedicacion.HistorialService
	MedicacionHistorial *servicioMedicacion.MedicacionHistorialService
}

// NewServer solo recibe los 4 servicios principales, sin Auth
func NewServer(
	medicacion *servicioMedicacion.MedicacionService,
	paciente *servicioMedicacion.PacienteService,
	historial *servicioMedicacion.HistorialService,
	medicacionHistorial *servicioMedicacion.MedicacionHistorialService,
) *Server {
	return &Server{
		Medicacion:          medicacion,
		Paciente:            paciente,
		Historial:           historial,
		MedicacionHistorial: medicacionHistorial,
	}
}