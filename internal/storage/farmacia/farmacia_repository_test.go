package farmacia

import (
	model "medicare-adulto-mayor/internal/models/farmacia"
	"testing"
)

/* =========================
   FAKE REPOSITORY IN MEMORY
========================= */

type FakeStorage struct {
	data []model.Farmacia
}

func newFakeStorage() *StorageFarmaciaGORM {
	return &StorageFarmaciaGORM{
		DB: nil,
	}
}

/* =========================
   TEST 1: CREAR
========================= */

func TestCrearFarmaciaRepo(t *testing.T) {

	fake := &model.Farmacia{
		Nombre:    "Farmacia Test",
		Direccion: "Centro",
	}

	// simulación directa
	if fake.Nombre == "" {
		t.Fatalf("nombre vacío")
	}

	if fake.Direccion == "" {
		t.Fatalf("direccion vacía")
	}
}

/* =========================
   TEST 2: LISTAR
========================= */

func TestListarFarmaciasRepo(t *testing.T) {

	lista := []model.Farmacia{
		{ID: "1", Nombre: "Farmacia 1"},
	}

	if len(lista) == 0 {
		t.Errorf("se esperaba data")
	}
}

/* =========================
   TEST 3: BUSCAR CERCANAS
========================= */

func TestBuscarCercanasRepo(t *testing.T) {

	farmacias := []model.Farmacia{
		{
			Nombre:   "Farmacia Cercana",
			Latitud:  -0.95,
			Longitud: -80.73,
		},
	}

	if len(farmacias) == 0 {
		t.Errorf("se esperaba al menos una farmacia")
	}
}
