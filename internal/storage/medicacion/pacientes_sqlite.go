package medicacion

import (
	"proyecto-medicare-adulto-mayor/internal/models/medicacion"
)


// pacientes
func (a *AlmacenSQLite) ListarPacientes() ([]medicacion.Paciente, error) {
    var pacientes []medicacion.Paciente
    result := a.db.Find(&pacientes)
    return pacientes, result.Error
}

func (a *AlmacenSQLite) CrearPaciente(p medicacion.Paciente) (medicacion.Paciente, error) {
    result := a.db.Create(&p)
    return p, result.Error
}

func (a *AlmacenSQLite) BuscarPacientePorID(id int) (medicacion.Paciente, error) {
    var paciente medicacion.Paciente
    if err := a.db.First(&paciente, id).Error; err != nil {
        return medicacion.Paciente{}, err
    }
    return paciente, nil
}

func (a *AlmacenSQLite) ActualizarPaciente(id int, p medicacion.Paciente) (medicacion.Paciente, error) {
    var paciente medicacion.Paciente
    if err := a.db.First(&paciente, id).Error; err != nil {
        return medicacion.Paciente{}, err
    }
    p.ID = paciente.ID
    if err := a.db.Save(&p).Error; err != nil {
        return medicacion.Paciente{}, err
    }
    return p, nil
}

func (a *AlmacenSQLite) EliminarPaciente(id int) (bool, error) {
    result := a.db.Delete(&medicacion.Paciente{}, id)

    if result.Error != nil {
        return false, result.Error
    }
    if result.RowsAffected == 0 {
        return false, nil
    }
    return true, nil
}