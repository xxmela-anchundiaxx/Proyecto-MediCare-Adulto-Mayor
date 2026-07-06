package medicacion

import "gorm.io/gorm"

// AlmacenSQLite implementa MedicacionRepository, PacientesRepository
// y HistorialRepository usando GORM sobre SQLite.
type AlmacenSQLite struct {
	db *gorm.DB
}

func NewAlmacenSQLite(db *gorm.DB) *AlmacenSQLite {
	return &AlmacenSQLite{db: db}
}