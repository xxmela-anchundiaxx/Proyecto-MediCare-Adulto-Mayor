package models

import "time"

type Medicacion struct {
	ID                  int         `json:"id"`
	PacienteID          int         `json:"paciente_id"`
	Nombre              string      `json:"nombre"`
	Descripcion         string      `json:"descripcion"`
	Dosis               string      `json:"dosis"`
	Frecuencia          string      `json:"frecuencia"`
	Hora_programada     string      `json:"hora_programada"`
	Inicio_tratamiento  time.Time   `json:"inicio_tratamiento"`
	Fecha_creacion      time.Time   `json:"fecha_creacion"`
}