package storage

import "github.com/xxmela-anchundiaxx/Proyecto-MediCare-Adulto-Mayor/internal/models"

var Relaciones []models.CuidadorPaciente

func CrearRelacion(relacion models.CuidadorPaciente) {
	Relaciones = append(Relaciones, relacion)
}

func ObtenerRelaciones() []models.CuidadorPaciente {
	return Relaciones
}

func ObtenerRelacionPorID(id int) (models.CuidadorPaciente, bool) {
	for _, relacion := range Relaciones {
		if relacion.ID == id {
			return relacion, true
		}
	}

	return models.CuidadorPaciente{}, false
}
