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
    // 1. Conexión a SQLite medicacion.db
    db, err := gorm.Open(sqlite.Open("medicare.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("No se pudo conectar a la base de datos:", err)
    }

    // 2. Migrar tablas según los modelos
    if err := db.AutoMigrate(
        &models.Medicacion{},
        &models.Paciente{},
        &models.HistorialMedicacion{},
    ); err != nil {
        log.Fatal("Error al migrar la base de datos:", err)
    }

    // 3. Crear almacenamiento y sembrar datos iniciales
    almacen := storage.NewAlmacenSQLite(db)
    almacen.SembrarVacioMedicacion()
    
    // 4. Crear el handler inyectando el almacenamiento
    medicamentoHandler := handlers.NewMedicamentoHandler(almacen)
    pacienteHandler := handlers.NewPacienteHandler(almacen)
    historialHandler := handlers.NewHistorialHandler(almacen)

    // 5. Configurar router con Chi
    r := chi.NewRouter()

    // 6. Middleware global
    r.Use(chimw.Logger)
    r.Use(chimw.Recoverer)
    r.Use(middleware.CORS)

    // 7. Rutas versionadas /api/v1/
    r.Route("/api/v1", func(r chi.Router) {
        // Rutas para medicacion
        r.Get("/medicaciones", medicamentoHandler.ListarMedicacion)
        r.Post("/medicaciones", medicamentoHandler.CrearMedicacion)
        r.Get("/medicaciones/{id}", medicamentoHandler.ObtenerMedicacion)
        r.Put("/medicaciones/{id}", medicamentoHandler.ActualizarMedicacion)
        r.Delete("/medicaciones/{id}", medicamentoHandler.EliminarMedicacion)

        r.Get("/pacientes", pacienteHandler.ListarPacientes)
        r.Post("/pacientes", pacienteHandler.CrearPaciente)
        r.Get("/pacientes/{id}", pacienteHandler.BuscarPacientePorID)
        r.Put("/pacientes/{id}", pacienteHandler.ActualizarPaciente)
        r.Delete("/pacientes/{id}", pacienteHandler.EliminarPaciente)

        r.Get("/historial", historialHandler.ListarHistorial) 
        r.Get("/historial/{id}", historialHandler.BuscarPorID)                 
        r.Post("/historial", historialHandler.Crear)                           
        r.Put("/historial/{id}", historialHandler.Actualizar)                  
        r.Delete("/historial/{id}", historialHandler.Eliminar) 

        r.Get("/pacientes/{id}/medicaciones", medicamentoHandler.ListarPorPaciente)
        r.Get("/pacientes/{id}/historial", historialHandler.ListarPorPaciente)

        r.Post("/medicamentos_farmacia", handlers.CreateMedicamento)
	    r.Get("/medicamentos_farmacia", handlers.GetMedicamentos)
	    r.Get("/medicamentos_farmacia/{id}", handlers.GetMedicamento)
	    r.Put("/medicamentos_farmacia/{id}", handlers.UpdateMedicamento)
	    r.Delete("/medicamentos_farmacia/{id}", handlers.DeleteMedicamento)

        // Monitoreo familiar/cuidador (memoria)
        r.Post("/relaciones", handlers.CrearRelacionHandler)
        r.Get("/relaciones", handlers.ObtenerRelacionesHandler)
        r.Get("/relaciones/{id}", handlers.ObtenerRelacionPorIDHandler)
        r.Put("/relaciones/{id}", handlers.ActualizarRelacionHandler)
        r.Delete("/relaciones/{id}", handlers.EliminarRelacionHandler)

    }) 

    log.Println("Servidor escuchando en http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
