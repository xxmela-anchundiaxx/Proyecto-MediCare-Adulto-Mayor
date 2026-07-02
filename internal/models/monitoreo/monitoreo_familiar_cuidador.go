package monitoreo

import "time"

type MonitoreoSignos struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	PacienteID     string    `json:"paciente_id"`
	RitmoCardiaco  int       `json:"ritmo_cardiaco,omitempty"`  // ppm
	PresionArterial string    `json:"presion_arterial,omitempty"` // ej. "120/80"
	NivelAzucar    float64   `json:"nivel_azucar,omitempty"`    // mg/dL
	Temperatura    float64   `json:"temperatura,omitempty"`     // Celsius
	FechaHora      time.Time `json:"fecha_hora"`
	AlertaEnviada  bool      `json:"alerta_enviada"`
	Observaciones  string    `json:"observaciones,omitempty"`
}

func (MonitoreoSignos) TableName() string {
	return "monitoreo_signos"
}


type RegistrarSignosRequest struct {
	PacienteID     string  `json:"paciente_id"`
	RitmoCardiaco  int     `json:"ritmo_cardiaco,omitempty"`
	PresionArterial string  `json:"presion_arterial,omitempty"`
	NivelAzucar    float64 `json:"nivel_azucar,omitempty"`
	Temperatura    float64 `json:"temperatura,omitempty"`
	Observaciones  string  `json:"observaciones,omitempty"`
}
