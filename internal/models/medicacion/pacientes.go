package medicacion

type Paciente struct {
	ID                 string `json:"id" gorm:"primaryKey"`
	UsuarioID          string `json:"usuario_id"`
	CuidadorID         string `json:"cuidador_id,omitempty"`
	Edad               int    `json:"edad"`
	GrupoSanguineo     string `json:"grupo_sanguineo,omitempty"`
	Alergias           string `json:"alergias,omitempty"`
	CondicionesMedicas string `json:"condiciones_medicas,omitempty"`
	ContactoEmergencia string `json:"contacto_emergencia"`
}

func (Paciente) TableName() string {
	return "pacientes"
}


type CreatePacienteRequest struct {
	UsuarioID          string `json:"usuario_id"`
	CuidadorID         string `json:"cuidador_id,omitempty"`
	Edad               int    `json:"edad"`
	GrupoSanguineo     string `json:"grupo_sanguineo,omitempty"`
	Alergias           string `json:"alergias,omitempty"`
	CondicionesMedicas string `json:"condiciones_medicas,omitempty"`
	ContactoEmergencia string `json:"contacto_emergencia"`
}
