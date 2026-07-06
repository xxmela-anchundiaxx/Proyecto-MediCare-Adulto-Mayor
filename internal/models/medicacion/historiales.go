package medicacion

import "time"

type HistorialMedicacion struct {
	ID           int       `json:"id"`
	MedicacionID int       `json:"medicacion_id"`
	FechaHora    time.Time `json:"fecha_hora"`
	Tomada       bool      `json:"tomada"`
	Observacion  string    `json:"observacion"`
}

func (HistorialMedicacion) TableName() string {
	return "historial_medicaciones"
}
