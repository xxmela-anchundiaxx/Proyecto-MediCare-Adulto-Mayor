package storage

import (
	"gorm.io/gorm"
	"proyecto-medicare-adulto-mayor/internal/models"
	"time"
)

type AlmacenSQLite struct {
	db *gorm.DB
}

// medicacion
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

func (a *AlmacenSQLite) ActualizarMedicacion(id int, datos models.Medicacion) (models.Medicacion, error) {
    var medicacion models.Medicacion
    if err := a.db.First(&medicacion, id).Error; err != nil {
        return models.Medicacion{}, err
    }
    datos.ID = medicacion.ID
    if err := a.db.Save(&datos).Error; err != nil {
        return models.Medicacion{}, err
    }
    return datos, nil
}


func (a *AlmacenSQLite) EliminarMedicacion(id int) (bool, error) {
    result := a.db.Delete(&models.Medicacion{}, id)

    if result.Error != nil {
        // Error al ejecutar el DELETE
        return false, result.Error
    }

    if result.RowsAffected == 0 {
        return false, nil
    }
    return true, nil
}

// pacientes
func (a *AlmacenSQLite) ListarPacientes() ([]models.Paciente, error) {
    var pacientes []models.Paciente
    result := a.db.Find(&pacientes)
    return pacientes, result.Error
}

func (a *AlmacenSQLite) CrearPaciente(p models.Paciente) (models.Paciente, error) {
    result := a.db.Create(&p)
    return p, result.Error
}

func (a *AlmacenSQLite) BuscarPacientePorID(id int) (models.Paciente, error) {
    var paciente models.Paciente
    if err := a.db.First(&paciente, id).Error; err != nil {
        return models.Paciente{}, err
    }
    return paciente, nil
}

func (a *AlmacenSQLite) ActualizarPaciente(id int, p models.Paciente) (models.Paciente, error) {
    var paciente models.Paciente
    if err := a.db.First(&paciente, id).Error; err != nil {
        return models.Paciente{}, err
    }
    p.ID = paciente.ID
    if err := a.db.Save(&p).Error; err != nil {
        return models.Paciente{}, err
    }
    return p, nil
}

func (a *AlmacenSQLite) EliminarPaciente(id int) (bool, error) {
    result := a.db.Delete(&models.Paciente{}, id)

    if result.Error != nil {
        return false, result.Error
    }
    if result.RowsAffected == 0 {
        return false, nil
    }
    return true, nil
}


// historial
func (a *AlmacenSQLite) ListarHistorial() ([]models.HistorialMedicacion, error) {
    var historiales []models.HistorialMedicacion
    result := a.db.Find(&historiales)
    return historiales, result.Error
}

func (a *AlmacenSQLite) BuscarHistorialPorID(id int) (models.HistorialMedicacion, error) {
    var h models.HistorialMedicacion
    if err := a.db.First(&h, id).Error; err != nil {
        return models.HistorialMedicacion{}, err
    }
    return h, nil
}

func (a *AlmacenSQLite) CrearHistorial(h models.HistorialMedicacion) (models.HistorialMedicacion, error) {
    result := a.db.Create(&h)
    return h, result.Error
}

func (a *AlmacenSQLite) ActualizarHistorial(id int, datos models.HistorialMedicacion) (models.HistorialMedicacion, error) {
    var h models.HistorialMedicacion
    if err := a.db.First(&h, id).Error; err != nil {
        return models.HistorialMedicacion{}, err
    }
    datos.ID = h.ID
    if err := a.db.Save(&datos).Error; err != nil {
        return models.HistorialMedicacion{}, err
    }
    return datos, nil
}

func (a *AlmacenSQLite) EliminarHistorial(id int) (bool, error) {
    result := a.db.Delete(&models.HistorialMedicacion{}, id)

    if result.Error != nil {
        return false, result.Error
    }
    if result.RowsAffected == 0 {
        return false, nil
    }
    return true, nil
}


// Listar todas las medicaciones de un paciente
func (a *AlmacenSQLite) ListarMedicacionPorPaciente(pacienteID int) ([]models.Medicacion, error) {
    var medicaciones []models.Medicacion
    result := a.db.Where("paciente_id = ?", pacienteID).Find(&medicaciones)
    return medicaciones, result.Error
}

// Listar todo el historial de un paciente
func (a *AlmacenSQLite) ListarHistorialPorPaciente(pacienteID int) ([]models.HistorialMedicacion, error) {
    var historiales []models.HistorialMedicacion
    result := a.db.Joins("JOIN medicacions ON medicacions.id = historial_medicacions.medicacion_id").
        Where("medicacions.paciente_id = ?", pacienteID).
        Find(&historiales)
    return historiales, result.Error
}



func (a *AlmacenSQLite) SembrarVacioMedicacion() {
    var n int64
    a.db.Model(&models.Medicacion{}).Count(&n)

    if n == 0 {
        medicaciones := []models.Medicacion{
            {
                Nombre:             "Paracetamol",
                Descripcion:        "Analgésico y antipirético",
                Dosis:              "500mg",
                Frecuencia:         "Cada 8 horas",
                Hora_programada:    "08:00",
                Inicio_tratamiento: time.Now(),
                Fecha_creacion:     time.Now(),
            },
            {
                Nombre:             "Ibuprofeno",
                Descripcion:        "Antiinflamatorio no esteroideo",
                Dosis:              "400mg",
                Frecuencia:         "Cada 12 horas",
                Hora_programada:    "09:00",
                Inicio_tratamiento: time.Now(),
                Fecha_creacion:     time.Now(),
            },
            {
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







