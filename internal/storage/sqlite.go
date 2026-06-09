package storage

import (
	"gorm.io/gorm"
	"proyecto-medicare-adulto-mayor/internal/models"
	"time"
)

type AlmacenSQLite struct {
	db *gorm.DB
}

func NewAlmacenSQLite(db *gorm.DB) *AlmacenSQLite {
	return &AlmacenSQLite{db: db}
}

func (a *AlmacenSQLite) ListarMedicacion() ([]models.Medicacion, error) {
	var medicacion []models.Medicacion
	result := a.db.Find(&medicacion)
	return medicacion, result.Error
}

func (a *AlmacenSQLite) CrearMedicacion(m models.Medicacion) (models.Medicacion, error) {
    result := a.db.Create(&m)
    return m, result.Error
}

func (a *AlmacenSQLite) BuscarMedicacionPorID(id int) (models.Medicacion, error) {
	var medicacion models.Medicacion
	if err := a.db.First(&medicacion, id).Error; err != nil {
        return models.Medicacion{}, err
    }
    return medicacion, nil
}

func (a *AlmacenSQLite) SembrarVacioMedicacion() {
    var n int64
    a.db.Model(&models.Medicacion{}).Count(&n)

    if n == 0 {
        medicaciones := []models.Medicacion{
            {
                PacienteID:         1,
                Nombre:             "Paracetamol",
                Descripcion:        "Analgésico y antipirético",
                Dosis:              "500mg",
                Frecuencia:         "Cada 8 horas",
                Hora_programada:    "08:00",
                Inicio_tratamiento: time.Now(),
                Fecha_creacion:     time.Now(),
            },
            {
                PacienteID:         2,
                Nombre:             "Ibuprofeno",
                Descripcion:        "Antiinflamatorio no esteroideo",
                Dosis:              "400mg",
                Frecuencia:         "Cada 12 horas",
                Hora_programada:    "09:00",
                Inicio_tratamiento: time.Now(),
                Fecha_creacion:     time.Now(),
            },
            {
                PacienteID:         3,
                Nombre:             "Amoxicilina",
                Descripcion:        "Antibiótico de amplio espectro",
                Dosis:              "500mg",
                Frecuencia:         "Cada 8 horas",
                Hora_programada:    "07:00",
                Inicio_tratamiento: time.Now(),
                Fecha_creacion:     time.Now(),
            },
        }

        a.db.Create(&medicaciones)
    }
}

var _ Almacen = (*AlmacenSQLite)(nil)







