package models

import "time"

//módulo Medicación

type Medicacion struct {
	ID                  int         `json:"id"`
	Nombre              string      `json:"nombre"`
	Descripcion         string      `json:"descripcion"`
	Dosis               string      `json:"dosis"`
	Frecuencia          string      `json:"frecuencia"`
	Hora_programada     string      `json:"hora_programada"`
	Inicio_tratamiento  time.Time   `json:"inicio_tratamiento"`
	Fecha_creacion      time.Time   `json:"fecha_creacion"`

	PacienteID int      `json:"paciente_id"`
    Historial  []HistorialMedicacion `gorm:"foreignKey:MedicacionID;references:ID"`
}

//paciente (1) ----- (N) medicacion
//medicacion(1) ----- (N) historial_medicaciones

type Paciente struct {
    ID     int    `json:"id"`
    Nombre string `json:"nombre"`
    Edad   int    `json:"edad"`

	Medicaciones []Medicacion `gorm:"foreignKey:PacienteID;references:ID"`
}

type HistorialMedicacion struct {
    ID           int       `json:"id"`
    MedicacionID int       `json:"medicacion_id"`
    FechaHora    time.Time `json:"fecha_hora"`
    Tomada       bool      `json:"tomada"`
    Observacion  string    `json:"observacion"`
}
