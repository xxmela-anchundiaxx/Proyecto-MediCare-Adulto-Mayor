package farmacia

type BusquedaFarmaciaRequest struct {
	Latitud  float64 `json:"latitud"`
	Longitud float64 `json:"longitud"`
	RadioKM  float64 `json:"radio_km"`
}