package medicacion

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	models "proyecto-medicare-adulto-mayor/internal/models/medicacion"
)

// Función auxiliar para inicializar la base de datos de pruebas en memoria
func setupTestSQLite(t *testing.T) *AlmacenSQLite {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error al abrir base de datos SQLite en memoria: %v", err)
	}

	// Migrar el modelo estructural de medicación
	err = db.AutoMigrate(&models.Medicacion{})
	if err != nil {
		t.Fatalf("Error al migrar tablas de pruebas: %v", err)
	}

	return &AlmacenSQLite{db: db}
}

// 1. Test Listar Medicación
func TestAlmacenSQLite_ListarMedicacion(t *testing.T) {
	almacen := setupTestSQLite(t)

	// Insertar datos dummy directo con GORM
	almacen.db.Create(&models.Medicacion{Nombre: "Losartán", Dosis: "50mg"})
	almacen.db.Create(&models.Medicacion{Nombre: "Metformina", Dosis: "850mg"})

	meds, err := almacen.ListarMedicacion()
	if err != nil {
		t.Fatalf("Error inesperado al listar medicación: %v", err)
	}

	if len(meds) != 2 {
		t.Errorf("Se esperaban 2 registros de medicación, se obtuvieron: %d", len(meds))
	}
}

// 2. Test Crear Medicación
func TestAlmacenSQLite_CrearMedicacion(t *testing.T) {
	almacen := setupTestSQLite(t)

	m := models.Medicacion{
		Nombre:     "Omeprazol",
		Dosis:      "20mg",
		Frecuencia: "Cada 24 horas",
	}

	nueva, err := almacen.CrearMedicacion(m)
	if err != nil {
		t.Fatalf("Error inesperado al crear medicación: %v", err)
	}

	if nueva.ID == 0 {
		t.Error("La medicación creada debería tener un ID autoincremental asignado")
	}

	// Verificar persistencia física en la base de datos
	var chequeo models.Medicacion
	if err := almacen.db.First(&chequeo, nueva.ID).Error; err != nil {
		t.Fatalf("No se encontró el registro persistido en SQLite: %v", err)
	}
}

// 3. Test Buscar por ID
func TestAlmacenSQLite_BuscarMedicacionPorID(t *testing.T) {
	almacen := setupTestSQLite(t)

	m := models.Medicacion{Nombre: "Atorvastatina", Dosis: "20mg"}
	almacen.db.Create(&m)

	med, err := almacen.BuscarMedicacionPorID(m.ID)
	if err != nil {
		t.Fatalf("Error inesperado al buscar por ID: %v", err)
	}

	if med.Nombre != "Atorvastatina" {
		t.Errorf("Se esperaba la medicación 'Atorvastatina', llegó: %s", med.Nombre)
	}

	// Test de error de registro inexistente
	_, err = almacen.BuscarMedicacionPorID(9999)
	if err == nil {
		t.Error("Se esperaba un error al buscar un ID aleatorio que no existe, pero retornó nil")
	}
}

// 4. Test Actualizar Medicación
func TestAlmacenSQLite_ActualizarMedicacion(t *testing.T) {
	almacen := setupTestSQLite(t)

	m := models.Medicacion{Nombre: "Clonazepam", Dosis: "0.5mg", Frecuencia: "Cada 24h"}
	almacen.db.Create(&m)

	datosNuevos := models.Medicacion{
		Nombre:     "Clonazepam Modificado",
		Dosis:      "1mg",
		Frecuencia: "Cada 12h",
	}

	actualizada, err := almacen.ActualizarMedicacion(m.ID, datosNuevos)
	if err != nil {
		t.Fatalf("Error inesperado al actualizar: %v", err)
	}

	if actualizada.Nombre != "Clonazepam Modificado" || actualizada.Dosis != "1mg" {
		t.Errorf("Los campos actualizados no coinciden con los guardados")
	}
}

// 5. Test Eliminar Medicación
func TestAlmacenSQLite_EliminarMedicacion(t *testing.T) {
	almacen := setupTestSQLite(t)

	m := models.Medicacion{Nombre: "Enalapril", Dosis: "10mg"}
	almacen.db.Create(&m)

	// Caso 1: Eliminar registro existente
	eliminado, err := almacen.EliminarMedicacion(m.ID)
	if err != nil {
		t.Fatalf("Error inesperado al eliminar: %v", err)
	}
	if !eliminado {
		t.Error("Se esperaba que 'eliminado' fuera true para un ID válido")
	}

	// Caso 2: Intentar eliminar un ID ya inexistente
	eliminadoFalso, err := almacen.EliminarMedicacion(m.ID)
	if err != nil {
		t.Fatalf("Error inesperado al intentar re-eliminar: %v", err)
	}
	if eliminadoFalso {
		t.Error("Se esperaba que 'eliminado' fuera false si RowsAffected es 0")
	}
}

// 6. Test Sembrar Datos en Vacío (Seeding)
func TestAlmacenSQLite_SembrarVacioMedicacion(t *testing.T) {
	almacen := setupTestSQLite(t)

	// Al estar vacía la base de datos recién creada, el Seeder debe activarse
	almacen.SembrarVacioMedicacion()

	var conteo int64
	almacen.db.Model(&models.Medicacion{}).Count(&conteo)

	if conteo != 3 {
		t.Errorf("El seeder debió registrar exactamente 3 medicaciones base, se encontraron: %d", conteo)
	}

	// Si llamamos de nuevo al Seeder, no debería duplicar datos ya que conteo ya no es 0
	almacen.SembrarVacioMedicacion()
	almacen.db.Model(&models.Medicacion{}).Count(&conteo)

	if conteo != 3 {
		t.Errorf("El seeder volvió a insertar duplicados incorrectamente. Conteo actual: %d", conteo)
	}
}