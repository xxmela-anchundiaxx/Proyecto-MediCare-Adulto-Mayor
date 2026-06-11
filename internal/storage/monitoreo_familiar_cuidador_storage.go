package storage

import "github.com/xxmela-anchundiaxx/Proyecto-MediCare-Adulto-Mayor/internal/models"

var Relaciones []models.CuidadorPaciente

func CrearRelacion(relacion models.CuidadorPaciente) {
	Relaciones = append(Relaciones, relacion)
}
