package medicacion

import "time"

type HistorialMedicacion struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	MedicamentoID string    `json:"medicamento_id"`
	PacienteID    string    `json:"paciente_id"`
	FechaHora     time.Time `json:"fecha_hora"`
	Tomado        bool      `json:"tomado"`
	Observaciones string    `json:"observaciones,omitempty"`
}

func (HistorialMedicacion) TableName() string {
	return "historial_medicacion"
}


type RecordAdherenceRequest struct {
	MedicamentoID string `json:"medicamento_id"`
	PacienteID    string `json:"paciente_id"`
	Tomado        bool   `json:"tomado"`
	Observaciones string `json:"observaciones,omitempty"`
}
