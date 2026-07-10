package medicacion

import (
	"testing"
	//"time"

	// CORRECCIÓN: Usa el driver Pure Go
	"github.com/glebarez/sqlite" 
	"gorm.io/gorm"

	"proyecto-medicare-adulto-mayor/internal/models"
	"proyecto-medicare-adulto-mayor/internal/models/medicacion"
)

// setupTestDB ahora utiliza el driver Pure Go
func setupTestDB(t *testing.T) *AlmacenSQLite {
	// sqlite.Open de glebarez/sqlite no requiere CGO
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("No se pudo iniciar base de datos en memoria: %v", err)
	}

	// Migramos los modelos
	err = db.AutoMigrate(&medicacion.HistorialMedicacion{}, &medicacion.Medicacion{}, &models.Usuario{})
	if err != nil {
		t.Fatalf("Error ejecutando automigrate de prueba: %v", err)
	}

	return &AlmacenSQLite{db: db}
}

// ... El resto de tus funciones Test... se mantienen igual.