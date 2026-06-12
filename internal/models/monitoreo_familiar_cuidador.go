package models

type CuidadorPaciente struct {
	ID         int    `json:"id"`
	CuidadorID int    `json:"cuidador_id"`
	PacienteID int    `json:"paciente_id"`
	Relacion   string `json:"relacion"`
}