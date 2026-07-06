package storage

import (
    "proyecto-medicare-adulto-mayor/internal/models/monitoreo"
    "gorm.io/gorm"
)

type MonitoreoSQLite struct {
    db *gorm.DB
}

func NewMonitoreoSQLite(db *gorm.DB) *MonitoreoSQLite {
    return &MonitoreoSQLite{db: db}
}


func (r *MonitoreoSQLite) ListarRelaciones() ([]monitoreo.CuidadorPaciente, error) {
    var relaciones []monitoreo.CuidadorPaciente
    result := r.db.Find(&relaciones)
    return relaciones, result.Error
}

func (r *MonitoreoSQLite) BuscarRelacionPorID(id int) (monitoreo.CuidadorPaciente, error) {
    var rel monitoreo.CuidadorPaciente
    if err := r.db.First(&rel, id).Error; err != nil {
        return monitoreo.CuidadorPaciente{}, err
    }
    return rel, nil
}

func (r *MonitoreoSQLite) CrearRelacion(rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error) {
    if err := r.db.Create(&rel).Error; err != nil {
        return monitoreo.CuidadorPaciente{}, err
    }
    return rel, nil
}

func (r *MonitoreoSQLite) ActualizarRelacion(id int, rel monitoreo.CuidadorPaciente) (monitoreo.CuidadorPaciente, error) {
    var existente monitoreo.CuidadorPaciente
    if err := r.db.First(&existente, id).Error; err != nil {
        return monitoreo.CuidadorPaciente{}, err
    }
    rel.ID = existente.ID
    if err := r.db.Save(&rel).Error; err != nil {
        return monitoreo.CuidadorPaciente{}, err
    }
    return rel, nil
}

func (r *MonitoreoSQLite) EliminarRelacion(id int) (bool, error) {
    result := r.db.Delete(&monitoreo.CuidadorPaciente{}, id)
    if result.Error != nil {
        return false, result.Error
    }
    if result.RowsAffected == 0 {
        return false, nil
    }
    return true, nil
}