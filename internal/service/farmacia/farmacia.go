package service

import (
	"errors"
	"proyecto-medicare-adulto-mayor/internal/models/farmacia"
	farmaciaStorage "proyecto-medicare-adulto-mayor/internal/storage/farmacia"
	"github.com/google/uuid"
)

type ServicioFarmacia struct {
	Repo farmaciaStorage.RepositorioFarmacia
}

func NuevoServicioFarmacia(repo farmaciaStorage.RepositorioFarmacia) *ServicioFarmacia {
	return &ServicioFarmacia{Repo: repo}
}

func (s *ServicioFarmacia) RegistrarFarmacia(f *farmacia.Farmacia) error {
	if f.Nombre == "" || f.Direccion == "" {
		return errors.New("nombre y direccion de farmacia son obligatorios")
	}
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return s.Repo.CrearFarmacia(f)
}

func (s *ServicioFarmacia) ListarTodas() ([]farmacia.Farmacia, error) {
	return s.Repo.ListarTodas()
}

func (s *ServicioFarmacia) BuscarCercanas(lat, lon, radioKM float64) ([]farmacia.Farmacia, error) {
	if radioKM <= 0 {
		radioKM = 5.0 // Por defecto 5 km
	}
	return s.Repo.BuscarCercanas(lat, lon, radioKM)
}