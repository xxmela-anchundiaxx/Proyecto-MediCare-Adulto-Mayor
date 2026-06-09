package storage

import (
	//"proyecto-medicare-adulto-mayor/internal/models"
	"proyecto-medicare-adulto-mayor/internal/models"
)

type Almacen interface {
	ListarMedicacion() ([]models.Medicacion, error)
	BuscarMedicacionPorID(id int) (models.Medicacion, error)
	CrearMedicacion(medicacion models.Medicacion) (models.Medicacion, error)
	ActualizarMedicacion(id int, medicacion models.Medicacion) (models.Medicacion,bool, error)
	//EliminarMedicacion(id int) (bool, error)
}
