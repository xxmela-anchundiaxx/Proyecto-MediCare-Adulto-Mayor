package medicacion

import (
    "proyecto-medicare-adulto-mayor/internal/models/medicacion"
    "time"
)


func (a *AlmacenSQLite) ListarMedicacion() ([]medicacion.Medicacion, error) {
    var meds []medicacion.Medicacion
    result := a.db.Find(&meds)
    return meds, result.Error
}

func (a *AlmacenSQLite) CrearMedicacion(m medicacion.Medicacion) (medicacion.Medicacion, error) {
    result := a.db.Create(&m)
    return m, result.Error
}

func (a *AlmacenSQLite) BuscarMedicacionPorID(id int) (medicacion.Medicacion, error) {
    var med medicacion.Medicacion
    if err := a.db.First(&med, id).Error; err != nil {
        return medicacion.Medicacion{}, err
    }
    return med, nil
}

func (a *AlmacenSQLite) ActualizarMedicacion(id int, datos medicacion.Medicacion) (medicacion.Medicacion, error) {
    var med medicacion.Medicacion
    if err := a.db.First(&med, id).Error; err != nil {
        return medicacion.Medicacion{}, err
    }
    datos.ID = med.ID
    if err := a.db.Save(&datos).Error; err != nil {
        return medicacion.Medicacion{}, err
    }
    return datos, nil
}

func (a *AlmacenSQLite) EliminarMedicacion(id int) (bool, error) {
    result := a.db.Delete(&medicacion.Medicacion{}, id)
    if result.Error != nil {
        return false, result.Error
    }
    if result.RowsAffected == 0 {
        return false, nil
    }
    return true, nil
}

func (a *AlmacenSQLite) SembrarVacioMedicacion() {
    var n int64
    a.db.Model(&medicacion.Medicacion{}).Count(&n)

    if n == 0 {
        meds := []medicacion.Medicacion{
            {
                Nombre:            "Paracetamol",
                Descripcion:       "Analgésico y antipirético",
                Dosis:             "500mg",
                Frecuencia:        "Cada 8 horas",
                Hora_programada:    "08:00",
                Inicio_tratamiento: time.Now(),
                Fecha_creacion:     time.Now(),
            },
            {
                Nombre:            "Ibuprofeno",
                Descripcion:       "Antiinflamatorio no esteroideo",
                Dosis:             "400mg",
                Frecuencia:        "Cada 12 horas",
                Hora_programada:    "09:00",
                Inicio_tratamiento: time.Now(),
                Fecha_creacion:     time.Now(),
            },
            {
                Nombre:            "Amoxicilina",
                Descripcion:       "Antibiótico de amplio espectro",
                Dosis:             "500mg",
                Frecuencia:        "Cada 8 horas",
                Hora_programada:    "07:00",
                Inicio_tratamiento: time.Now(),
                Fecha_creacion:     time.Now(),
            },
        }
        a.db.Create(&meds)
    }
}
