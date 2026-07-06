package medicacion

import (
	"proyecto-medicare-adulto-mayor/internal/models/medicacion"
	"proyecto-medicare-adulto-mayor/internal/models"
)

// historial
func (a *AlmacenSQLite) ListarHistorial() ([]medicacion.HistorialMedicacion, error) {
    var historiales []medicacion.HistorialMedicacion
    result := a.db.Find(&historiales)
    return historiales, result.Error
}

func (a *AlmacenSQLite) BuscarHistorialPorID(id int) (medicacion.HistorialMedicacion, error) {
    var h medicacion.HistorialMedicacion
    if err := a.db.First(&h, id).Error; err != nil {
        return medicacion.HistorialMedicacion{}, err
    }
    return h, nil
}

func (a *AlmacenSQLite) CrearHistorial(h medicacion.HistorialMedicacion) (medicacion.HistorialMedicacion, error) {
    result := a.db.Create(&h)
    return h, result.Error
}

func (a *AlmacenSQLite) ActualizarHistorial(id int, datos medicacion.HistorialMedicacion) (medicacion.HistorialMedicacion, error) {
    var h medicacion.HistorialMedicacion
    if err := a.db.First(&h, id).Error; err != nil {
        return medicacion.HistorialMedicacion{}, err
    }
    datos.ID = h.ID
    if err := a.db.Save(&datos).Error; err != nil {
        return medicacion.HistorialMedicacion{}, err
    }
    return datos, nil
}

func (a *AlmacenSQLite) EliminarHistorial(id int) (bool, error) {
    result := a.db.Delete(&medicacion.HistorialMedicacion{}, id)

    if result.Error != nil {
        return false, result.Error
    }
    if result.RowsAffected == 0 {
        return false, nil
    }
    return true, nil
}


// Listar todas las medicaciones de un paciente
func (a *AlmacenSQLite) ListarMedicacionPorPaciente(pacienteID int) ([]medicacion.Medicacion, error) {
    var medicaciones []medicacion.Medicacion
    result := a.db.Where("paciente_id = ?", pacienteID).Find(&medicaciones)
    return medicaciones, result.Error
}

// Listar todo el historial de un paciente
func (a *AlmacenSQLite) ListarHistorialPorPaciente(pacienteID int) ([]medicacion.HistorialMedicacion, error) {
    var historiales []medicacion.HistorialMedicacion
    result := a.db.Joins("JOIN medicacions ON medicacions.id = historial_medicacions.medicacion_id").
        Where("medicacions.paciente_id = ?", pacienteID).
        Find(&historiales)
    return historiales, result.Error
}


func (a *AlmacenSQLite) CrearUsuario(email, passwordHash string) (models.Usuario, error) {
    usuario := models.Usuario{
        Email:        email,
        PasswordHash: passwordHash,
    }
    result := a.db.Create(&usuario)
    return usuario, result.Error
}

func (a *AlmacenSQLite) BuscarUsuarioPorEmail(email string) (models.Usuario, error) {
    var usuario models.Usuario
    if err := a.db.Where("email = ?", email).First(&usuario).Error; err != nil {
        return models.Usuario{}, err
    }
    return usuario, nil
}