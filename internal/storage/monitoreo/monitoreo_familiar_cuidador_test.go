package storage

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"proyecto-medicare-adulto-mayor/internal/models/monitoreo"
)

// Función auxiliar para configurar la base de datos SQLite de prueba en memoria
func setupTestMonitoreoDB(t *testing.T) *MonitoreoSQLite {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error al abrir base de datos SQLite de pruebas en memoria: %v", err)
	}

	// estructura del modelo de monitoreo para crear las tablas físicas en memoria
	err = db.AutoMigrate(&monitoreo.CuidadorPaciente{})
	if err != nil {
		t.Fatalf("Error al realizar la migración del modelo de monitoreo: %v", err)
	}

	return NewMonitoreoSQLite(db)
}

// 1. Test Listar Relaciones
func TestMonitoreoSQLite_ListarRelaciones(t *testing.T) {
	repo := setupTestMonitoreoDB(t)

	repo.db.Create(&monitoreo.CuidadorPaciente{CuidadorID: 1, PacienteID: 10, Relacion: "Hijo"})
	repo.db.Create(&monitoreo.CuidadorPaciente{CuidadorID: 2, PacienteID: 20, Relacion: "Enfermero"})

	resultados, err := repo.ListarRelaciones()
	if err != nil {
		t.Fatalf("No se esperaba error al listar relaciones: %v", err)
	}

	if len(resultados) != 2 {
		t.Errorf("Se esperaban 2 relaciones en la lista, se obtuvieron: %d", len(resultados))
	}
}

// 2. Test Buscar Relación Por ID
func TestMonitoreoSQLite_BuscarRelacionPorID(t *testing.T) {
	repo := setupTestMonitoreoDB(t)

	nueva := monitoreo.CuidadorPaciente{CuidadorID: 5, PacienteID: 15, Relacion: "Esposa"}
	repo.db.Create(&nueva)

	// Caso: Encontrar registro existente
	rel, err := repo.BuscarRelacionPorID(nueva.ID)
	if err != nil {
		t.Fatalf("Error inesperado al buscar relación por ID: %v", err)
	}
	if rel.Relacion != "Esposa" {
		t.Errorf("Se esperaba encontrar la relación 'Esposa', se obtuvo: %s", rel.Relacion)
	}

	// Caso: Registro no existente
	_, err = repo.BuscarRelacionPorID(9999)
	if err == nil {
		t.Error("Se esperaba un error al buscar un ID inexistente, pero se obtuvo nil")
	}
}

// 3. Test Crear Relación
func TestMonitoreoSQLite_CrearRelacion(t *testing.T) {
	repo := setupTestMonitoreoDB(t)

	rel := monitoreo.CuidadorPaciente{
		CuidadorID: 3,
		PacienteID: 30,
		Relacion:   "Voluntario",
	}

	creada, err := repo.CrearRelacion(rel)
	if err != nil {
		t.Fatalf("No se esperaba error al crear la relación: %v", err)
	}

	if creada.ID == 0 {
		t.Error("La relación guardada debería retornar un ID autoincremental válido")
	}

	// Validamos que se encuentre efectivamente en la base de datos
	var chequeo monitoreo.CuidadorPaciente
	if err := repo.db.First(&chequeo, creada.ID).Error; err != nil {
		t.Fatalf("La relación no se guardó físicamente en la DB: %v", err)
	}
}

// 4. Test Actualizar Relación
func TestMonitoreoSQLite_ActualizarRelacion(t *testing.T) {
	repo := setupTestMonitoreoDB(t)

	existente := monitoreo.CuidadorPaciente{CuidadorID: 4, PacienteID: 40, Relacion: "Amigo"}
	repo.db.Create(&existente)

	datosNuevos := monitoreo.CuidadorPaciente{
		CuidadorID: 4,
		PacienteID: 40,
		Relacion:   "Hijo Adoptivo",
	}

	actualizada, err := repo.ActualizarRelacion(existente.ID, datosNuevos)
	if err != nil {
		t.Fatalf("Error inesperado al actualizar relación: %v", err)
	}

	if actualizada.Relacion != "Hijo Adoptivo" {
		t.Errorf("El campo Relacion no se actualizó correctamente, se obtuvo: %s", actualizada.Relacion)
	}
}

// 5. Test Eliminar Relación
func TestMonitoreoSQLite_EliminarRelacion(t *testing.T) {
	repo := setupTestMonitoreoDB(t)

	rel := monitoreo.CuidadorPaciente{CuidadorID: 9, PacienteID: 90, Relacion: "Vecino"}
	repo.db.Create(&rel)

	// Caso: Eliminar registro real existente
	eliminado, err := repo.EliminarRelacion(rel.ID)
	if err != nil {
		t.Fatalf("Error inesperado al ejecutar eliminación: %v", err)
	}
	if !eliminado {
		t.Error("Se esperaba que la función retornara true al eliminar un registro real")
	}

	// Caso: Intentar eliminar un ID que ya no existe (RowsAffected == 0)
	eliminadoFalso, err := repo.EliminarRelacion(rel.ID)
	if err != nil {
		t.Fatalf("Error inesperado en segunda eliminación: %v", err)
	}
	if eliminadoFalso {
		t.Error("Se esperaba que retornara false si el registro a eliminar ya no existe")
	}
}