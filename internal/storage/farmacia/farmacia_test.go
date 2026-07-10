package farmacia

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	farmaciaModel "proyecto-medicare-adulto-mayor/internal/models/farmacia"
)


func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("No se pudo abrir la base de datos de prueba en memoria: %v", err)
	}

	// Migramos el modelo para que la tabla 'farmacias' exista en los tests
	err = db.AutoMigrate(&farmaciaModel.Farmacia{})
	if err != nil {
		t.Fatalf("Error al migrar base de datos de prueba: %v", err)
	}

	return db
}

// 1. Test Crear Farmacia
func TestStorageFarmaciaGORM_CrearFarmacia(t *testing.T) {
	db := setupTestDB(t)
	storage := NuevoStorageFarmaciaGORM(db)

	f := &farmaciaModel.Farmacia{
		ID:        "uuid-1",
		Nombre:    "Farmacia Sana Sana",
		Direccion: "Av. Principal 123",
		Telefono:  "0999999999",
	}

	err := storage.CrearFarmacia(f)
	if err != nil {
		t.Fatalf("No se esperaba error al crear la farmacia: %v", err)
	}

	// Verificar que realmente se guardó
	var resultado farmaciaModel.Farmacia
	if err := db.First(&resultado, "id = ?", "uuid-1").Error; err != nil {
		t.Fatalf("La farmacia no se guardó correctamente en la DB: %v", err)
	}
}

// 2. Test Buscar por ID
func TestStorageFarmaciaGORM_BuscarPorID(t *testing.T) {
	db := setupTestDB(t)
	storage := NuevoStorageFarmaciaGORM(db)

	f := farmaciaModel.Farmacia{ID: "uuid-2", Nombre: "Farmacia Cruz Azul", Direccion: "Calle B"}
	db.Create(&f)

	res, err := storage.BuscarPorID("uuid-2")
	if err != nil {
		t.Fatalf("Error inesperado al buscar por ID: %v", err)
	}

	if res.Nombre != "Farmacia Cruz Azul" {
		t.Errorf("Se esperaba encontrar 'Farmacia Cruz Azul', pero llegó: %s", res.Nombre)
	}

	// Test caso no encontrado
	_, err = storage.BuscarPorID("no-existe")
	if err == nil {
		t.Error("Se esperaba un error al buscar un ID inexistente, pero dio nil")
	}
}

// 3. Test Listar Todas
func TestStorageFarmaciaGORM_ListarTodas(t *testing.T) {
	db := setupTestDB(t)
	storage := NuevoStorageFarmaciaGORM(db)

	db.Create(&farmaciaModel.Farmacia{ID: "f1", Nombre: "F1"})
	db.Create(&farmaciaModel.Farmacia{ID: "f2", Nombre: "F2"})

	lista, err := storage.ListarTodas()
	if err != nil {
		t.Fatalf("Error inesperado al listar: %v", err)
	}

	if len(lista) != 2 {
		t.Errorf("Se esperaban 2 farmacias, se obtuvieron: %d", len(lista))
	}
}

// 4. Test Buscar Cercanas (Lógica de geolocalización por radio/grados)
func TestStorageFarmaciaGORM_BuscarCercanas(t *testing.T) {
	db := setupTestDB(t)
	storage := NuevoStorageFarmaciaGORM(db)

	// Farmacia muy cercana (Manta Centro)
	db.Create(&farmaciaModel.Farmacia{
		ID:        "cercana",
		Nombre:    "Farmacia Cerca",
		Latitud:   -0.950000,
		Longitud:  -80.720000,
	})

	// Farmacia muy lejana (En otra ciudad / país)
	db.Create(&farmaciaModel.Farmacia{
		ID:        "lejana",
		Nombre:    "Farmacia Lejos",
		Latitud:   40.7128,
		Longitud:  -74.0060,
	})

	// Buscamos con un radio de 5 KM en el mismo punto de la cercana
	resultados, err := storage.BuscarCercanas(-0.950000, -80.720000, 5.0)
	if err != nil {
		t.Fatalf("Error al buscar cercanas: %v", err)
	}

	if len(resultados) != 1 {
		t.Fatalf("Se esperaba encontrar solo 1 farmacia en el radio, se encontraron: %d", len(resultados))
	}

	if resultados[0].ID != "cercana" {
		t.Errorf("Se esperaba recuperar la farmacia 'cercana', se obtuvo: %s", resultados[0].ID)
	}
}

// 5. Test Actualizar Farmacia
func TestStorageFarmaciaGORM_ActualizarFarmacia(t *testing.T) {
	db := setupTestDB(t)
	storage := NuevoStorageFarmaciaGORM(db)

	f := farmaciaModel.Farmacia{ID: "uuid-update", Nombre: "Nombre Viejo", Direccion: "Dirección Vieja"}
	db.Create(&f)

	datosNuevos := &farmaciaModel.Farmacia{
		Nombre: "Nombre Nuevo",
	}

	err := storage.ActualizarFarmacia("uuid-update", datosNuevos)
	if err != nil {
		t.Fatalf("Error inesperado al actualizar: %v", err)
	}

	var resultado farmaciaModel.Farmacia
	db.First(&resultado, "id = ?", "uuid-update")

	if resultado.Nombre != "Nombre Nuevo" {
		t.Errorf("El nombre no se actualizó, sigue siendo: %s", resultado.Nombre)
	}
}

// 6. Test Eliminar Farmacia
func TestStorageFarmaciaGORM_EliminarFarmacia(t *testing.T) {
	db := setupTestDB(t)
	storage := NuevoStorageFarmaciaGORM(db)

	f := farmaciaModel.Farmacia{ID: "uuid-delete", Nombre: "A borrar"}
	db.Create(&f)

	err := storage.EliminarFarmacia("uuid-delete")
	if err != nil {
		t.Fatalf("Error inesperado al eliminar: %v", err)
	}

	var conteo int64
	db.Model(&farmaciaModel.Farmacia{}).Where("id = ?", "uuid-delete").Count(&conteo)

	if conteo != 0 {
		t.Error("La farmacia todavía existe en la base de datos y no fue eliminada")
	}
}