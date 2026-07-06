package medicacion

type Paciente struct {
    ID     int    `json:"id"`
    Nombre string `json:"nombre"`
    Edad   int    `json:"edad"`

	Medicaciones []Medicacion `gorm:"foreignKey:PacienteID;references:ID"`
}