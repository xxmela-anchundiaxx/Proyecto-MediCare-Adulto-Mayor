package medicacion

import "time"

type Medicamento struct {
	ID                string    `json:"id" gorm:"primaryKey"`
	PacienteID        string    `json:"paciente_id"`
	Nombre            string    `json:"nombre"`
	Descripcion       string    `json:"descripcion,omitempty"`
	Dosis             string    `json:"dosis"`
	Frecuencia        string    `json:"frecuencia"`
	ViaAdministracion string    `json:"via_administracion,omitempty"`
	Stock             int       `json:"stock"`
	FechaRegistro     time.Time `json:"fecha_registro"`
}

func (Medicamento) TableName() string {
	return "medicamentos"
}


type CreateMedicamentoRequest struct {
	PacienteID        string `json:"paciente_id"`
	Nombre            string `json:"nombre"`
	Descripcion       string `json:"descripcion,omitempty"`
	Dosis             string `json:"dosis"`
	Frecuencia        string `json:"frecuencia"`
	ViaAdministracion string `json:"via_administracion,omitempty"`
	Stock             int    `json:"stock"`
}
