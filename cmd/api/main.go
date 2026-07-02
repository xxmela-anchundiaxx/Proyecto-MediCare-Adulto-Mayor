package main

import (
	"log"
	"medicare-adulto-mayor/internal/handlers"
	"medicare-adulto-mayor/internal/storage"
	"os"
)

func main() {
	log.Println("Iniciando aplicación Medicare Adulto Mayor...")

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "medicare.db"
	}

	schemaPath := os.Getenv("SCHEMA_PATH")
	if schemaPath == "" {
		schemaPath = "db/schema.sql"
	}

	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	// 1. Inicializar almacenamiento SQLite
	almacen, err := storage.NuevoAlmacen(dbPath)
	if err != nil {
		log.Fatalf("Error crítico al inicializar almacenamiento: %v", err)
	}
	defer almacen.Cerrar()

	// 2. Cargar esquema inicial (Tablas)
	log.Println("Cargando esquema de base de datos...")
	if err := almacen.InicializarCargarEsquema(schemaPath); err != nil {
		log.Fatalf("Error crítico al cargar esquema de base de datos: %v", err)
	}
	log.Println("Base de datos SQLite sincronizada exitosamente.")

	// 3. Lanzar servidor REST API
	servidor := handlers.NuevoServidorMedicare(almacen)
	servidor.Iniciar(puerto)
}
