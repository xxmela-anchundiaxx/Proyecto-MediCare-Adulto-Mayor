package farmacia

type Farmacia struct {
	ID        string  `json:"id" gorm:"primaryKey"`
	Nombre    string  `json:"nombre"`
	Direccion string  `json:"direccion"`
	Telefono  string  `json:"telefono,omitempty"`
	Latitud   float64 `json:"latitud,omitempty"`
	Longitud  float64 `json:"longitud,omitempty"`
}

func (Farmacia) TableName() string {
	return "farmacias"
}


type BusquedaFarmaciaRequest struct {
	Latitud  float64 `json:"latitud"`
	Longitud float64 `json:"longitud"`
	RadioKM  float64 `json:"radio_km"`
}
