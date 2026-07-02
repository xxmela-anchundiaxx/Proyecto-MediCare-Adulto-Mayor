package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"medicare-adulto-mayor/internal/models"
	farmaciaModels "medicare-adulto-mayor/internal/models/farmacia"
	medicacionModels "medicare-adulto-mayor/internal/models/medicacion"
	monitoreoModels "medicare-adulto-mayor/internal/models/monitoreo"
)

type Almacen struct {
	DB   *sql.DB
	GORM *gorm.DB
}

func NuevoAlmacen(dbPath string) (*Almacen, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir base de datos: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error al conectar con base de datos: %w", err)
	}

	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error al inicializar GORM: %w", err)
	}

	// Ejecutar AutoMigrate de GORM
	err = gormDB.AutoMigrate(
		&models.Usuario{},
		&medicacionModels.Paciente{},
		&medicacionModels.Medicamento{},
		&medicacionModels.HistorialMedicacion{},
		&farmaciaModels.Farmacia{},
		&monitoreoModels.MonitoreoSignos{},
	)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error al ejecutar AutoMigrate: %w", err)
	}

	// Sembrar algunas farmacias por defecto si la tabla está vacía
	var count int64
	gormDB.Model(&farmaciaModels.Farmacia{}).Count(&count)
	if count == 0 {
		farmaciasSeed := []farmaciaModels.Farmacia{
			{
				ID:        "f1",
				Nombre:    "Farmacia Medicare Centro",
				Direccion: "Av. Principal 123, Centro",
				Telefono:  "+56 9 1234 5678",
				Latitud:   -33.456,
				Longitud:  -70.648,
			},
			{
				ID:        "f2",
				Nombre:    "Farmacia Cruz Salud Adulto Mayor",
				Direccion: "Calle Las Rosas 456, Providencia",
				Telefono:  "+56 9 8765 4321",
				Latitud:   -33.430,
				Longitud:  -70.612,
			},
		}
		for _, f := range farmaciasSeed {
			gormDB.Create(&f)
		}
	}

	return &Almacen{
		DB:   db,
		GORM: gormDB,
	}, nil
}

// InicializarCargarEsquema lee db/schema.sql y ejecuta las tablas iniciales si es necesario
func (a *Almacen) InicializarCargarEsquema(schemaPath string) error {
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo de esquema: %w", err)
	}

	_, err = a.DB.Exec(string(schemaBytes))
	if err != nil {
		// Ignorar errores de tabla ya existente si ya se auto-migraron
		return nil
	}

	return nil
}

func (a *Almacen) Cerrar() error {
	return a.DB.Close()
}

