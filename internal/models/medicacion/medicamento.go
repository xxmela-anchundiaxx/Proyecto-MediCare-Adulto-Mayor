package medicacion

import "time"

//módulo Medicación

type Medicacion struct {
	ID                 int       `json:"id"`
	Nombre             string    `json:"nombre"`
	Descripcion        string    `json:"descripcion"`
	Dosis              string    `json:"dosis"`
	Frecuencia         string    `json:"frecuencia"`
	Hora_programada    string    `json:"hora_programada"`
	Inicio_tratamiento time.Time `json:"inicio_tratamiento"`
	Fecha_creacion     time.Time `json:"fecha_creacion"`

	PacienteID int                   `json:"paciente_id"`
	Historial  []HistorialMedicacion `gorm:"foreignKey:MedicacionID;references:ID"`
}

func (Medicacion) TableName() string {
	return "medicaciones"
}

//paciente (1) ----- (N) medicacion
//medicacion(1) ----- (N) historial_medicaciones
