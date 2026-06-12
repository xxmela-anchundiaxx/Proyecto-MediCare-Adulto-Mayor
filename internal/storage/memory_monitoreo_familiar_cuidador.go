package storage

import "proyecto-medicare-adulto-mayor/internal/models"

var Relaciones []models.CuidadorPaciente
var IDRelacion int = 1

func CrearRelacion(relacion models.CuidadorPaciente) models.CuidadorPaciente {
	relacion.ID = IDRelacion
	IDRelacion++
	Relaciones = append(Relaciones, relacion)
	return relacion
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

func ActualizarRelacion(id int, nuevaRelacion models.CuidadorPaciente) bool {
	for i, relacion := range Relaciones {
		if relacion.ID == id {
			nuevaRelacion.ID = id
			Relaciones[i] = nuevaRelacion
			return true
		}
	}

	return false
}

func EliminarRelacion(id int) bool {
	for i, relacion := range Relaciones {
		if relacion.ID == id {
			Relaciones = append(Relaciones[:i], Relaciones[i+1:]...)
			return true
		}
	}

	return false
}
