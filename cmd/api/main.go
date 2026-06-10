package main

import (
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    chimw "github.com/go-chi/chi/v5/middleware"
    "github.com/glebarez/sqlite"
    "gorm.io/gorm"

    "proyecto-medicare-adulto-mayor/internal/handlers"
    "proyecto-medicare-adulto-mayor/internal/middleware"
    "proyecto-medicare-adulto-mayor/internal/models"
    "proyecto-medicare-adulto-mayor/internal/storage"
)

func main() {
    // 1. Conexión a SQLite
    db, err := gorm.Open(sqlite.Open("medicare.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("No se pudo conectar a la base de datos:", err)
    }

    // 2. Migrar tablas según los modelos
    if err := db.AutoMigrate(&models.Medicacion{}); err != nil {
        log.Fatal("Error al migrar la base de datos:", err)
    }

    // 3. Crear almacenamiento y sembrar datos iniciales
    almacen := storage.NewAlmacenSQLite(db)
    almacen.SembrarVacioMedicacion()

    // 4. Crear el handler inyectando el almacenamiento
    servidor := handlers.NewMedicamentoHandler(almacen)

    // 5. Configurar router con Chi
    r := chi.NewRouter()

    // 6. Middleware global
    r.Use(chimw.Logger)
    r.Use(chimw.Recoverer)
    r.Use(middleware.CORS)

    // 7. Rutas versionadas /api/v1/
    r.Route("/api/v1", func(r chi.Router) {
        r.Get("/medicaciones", servidor.ListarMedicacion)
        r.Post("/medicaciones", servidor.CrearMedicacion)
        r.Get("/medicaciones/{id}", servidor.ObtenerMedicacion)
        r.Put("/medicaciones/{id}", servidor.ActualizarMedicacion)
        // r.Delete("/medicaciones/{id}", servidor.BorrarMedicacion)
    }) 

    log.Println("Servidor escuchando en http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
