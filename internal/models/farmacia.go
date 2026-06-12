package models

type Medicamento_farmacia struct {
	ID         int     `json:"id"`
	Nombre     string  `json:"nombre"`
	Precio     float64 `json:"precio"`
	Disponible bool    `json:"disponible"`
	Farmacia   string  `json:"farmacia"`
	Horario    string  `json:"horario"`
}

