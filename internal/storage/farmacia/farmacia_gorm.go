package farmacia

import (
	"math"
	"medicare-adulto-mayor/internal/models/farmacia"
	"gorm.io/gorm"
)

type RepositorioFarmacia interface {
	ListarTodas() ([]farmacia.Farmacia, error)
	BuscarCercanas(lat, lon, radioKM float64) ([]farmacia.Farmacia, error)
	CrearFarmacia(f *farmacia.Farmacia) error
}

type StorageFarmaciaGORM struct {
	DB *gorm.DB
}

func NuevoStorageFarmaciaGORM(db *gorm.DB) *StorageFarmaciaGORM {
	return &StorageFarmaciaGORM{DB: db}
}

func (s *StorageFarmaciaGORM) ListarTodas() ([]farmacia.Farmacia, error) {
	var lista []farmacia.Farmacia
	err := s.DB.Find(&lista).Error
	return lista, err
}

func (s *StorageFarmaciaGORM) BuscarCercanas(lat, lon, radioKM float64) ([]farmacia.Farmacia, error) {
	var todas []farmacia.Farmacia
	if err := s.DB.Find(&todas).Error; err != nil {
		return nil, err
	}

	var cercanas []farmacia.Farmacia
	for _, f := range todas {
		dist := calcularDistanciaKM(lat, lon, f.Latitud, f.Longitud)
		if dist <= radioKM {
			cercanas = append(cercanas, f)
		}
	}
	return cercanas, nil
}

func (s *StorageFarmaciaGORM) CrearFarmacia(f *farmacia.Farmacia) error {
	return s.DB.Create(f).Error
}

func calcularDistanciaKM(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	rLat1 := lat1 * math.Pi / 180.0
	rLat2 := lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(rLat1)*math.Cos(rLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
