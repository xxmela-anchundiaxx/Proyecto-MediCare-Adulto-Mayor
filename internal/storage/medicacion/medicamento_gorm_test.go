package medicacion

import (
	"medicare-adulto-mayor/internal/models/medicacion"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestStorageMedicamentoGORM_RealDatabase_CrearYBuscar(t *testing.T) {
	// 1. Inicializar base de datos SQLite en memoria con GORM
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("error al abrir base de datos en memoria para el test: %v", err)
	}

	// 2. Ejecutar AutoMigrate de la entidad medicamento
	err = db.AutoMigrate(&medicacion.Medicamento{})
	if err != nil {
		t.Fatalf("error al auto-migrar tabla de medicamentos: %v", err)
	}

	// 3. Crear el repositorio GORM real
	repo := NuevoStorageMedicamentoGORM(db)

	m := &medicacion.Medicamento{
		ID:                "m-uniq-id-789",
		PacienteID:        "paciente-abc",
		Nombre:            "Atorvastatina",
		Descripcion:       "Tomar por las noches para el colesterol",
		Dosis:             "20mg",
		Frecuencia:        "Cada 24 horas",
		ViaAdministracion: "Oral",
		Stock:             45,
		FechaRegistro:     time.Now(),
	}

	// 4. Test: CrearMedicamento
	err = repo.CrearMedicamento(m)
	if err != nil {
		t.Fatalf("error al crear medicamento en BD real: %v", err)
	}

	// 5. Test: BuscarPorID (debe reflejar lo guardado)
	guardado, err := repo.BuscarPorID("m-uniq-id-789")
	if err != nil {
		t.Fatalf("error al buscar medicamento guardado: %v", err)
	}

	if guardado == nil {
		t.Fatal("se esperaba encontrar el medicamento guardado, pero retornó nil")
	}

	if guardado.Nombre != "Atorvastatina" {
		t.Errorf("se esperaba el nombre 'Atorvastatina', pero se obtuvo: %s", guardado.Nombre)
	}

	if guardado.Stock != 45 {
		t.Errorf("se esperaba stock 45, pero se obtuvo: %d", guardado.Stock)
	}

	// 6. Test: ListarPorPaciente (debe contener el medicamento)
	lista, err := repo.ListarPorPaciente("paciente-abc")
	if err != nil {
		t.Fatalf("error al listar medicamentos por paciente: %v", err)
	}

	if len(lista) != 1 {
		t.Errorf("se esperaba 1 medicamento en la lista, pero se obtuvo: %d", len(lista))
	} else if lista[0].Nombre != "Atorvastatina" {
		t.Errorf("se esperaba 'Atorvastatina', pero se obtuvo: %s", lista[0].Nombre)
	}

	// 7. Test: ActualizarStock -> BuscarPorID refleja el nuevo stock
	err = repo.ActualizarStock("m-uniq-id-789", 44)
	if err != nil {
		t.Fatalf("error al actualizar stock: %v", err)
	}

	modificado, err := repo.BuscarPorID("m-uniq-id-789")
	if err != nil {
		t.Fatalf("error al recuperar medicamento modificado: %v", err)
	}

	if modificado.Stock != 44 {
		t.Errorf("se esperaba stock actualizado de 44, pero se obtuvo: %d", modificado.Stock)
	}
}
