package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres" // <-- Driver necesario para producción en Docker
	"gorm.io/gorm"

	"proyecto-medicare-adulto-mayor/internal/handlers"
	handlersFarmacia "proyecto-medicare-adulto-mayor/internal/handlers/farmacia"
	handlersMedicacion "proyecto-medicare-adulto-mayor/internal/handlers/medicacion"
	handlersMonitoreo "proyecto-medicare-adulto-mayor/internal/handlers/monitoreo"
	"proyecto-medicare-adulto-mayor/internal/middleware"
	modelsFarmacia "proyecto-medicare-adulto-mayor/internal/models/farmacia"
	modelsMedicacion "proyecto-medicare-adulto-mayor/internal/models/medicacion"
	modelsMonitoreo "proyecto-medicare-adulto-mayor/internal/models/monitoreo"
	"proyecto-medicare-adulto-mayor/internal/service"
	servicioFarmacia "proyecto-medicare-adulto-mayor/internal/service/farmacia"
	servicioMedicacion "proyecto-medicare-adulto-mayor/internal/service/medicacion"
	servicioMonitoreo "proyecto-medicare-adulto-mayor/internal/service/monitoreo"
	"proyecto-medicare-adulto-mayor/internal/storage"
	storageFarmacia "proyecto-medicare-adulto-mayor/internal/storage/farmacia"
	storageMonitoreo "proyecto-medicare-adulto-mayor/internal/storage/monitoreo"
)

func main() {
	var dialector gorm.Dialector

	// 1. Comprobamos si estamos dentro del entorno Docker
	dbEnv := os.Getenv("DB_ENV")

	if dbEnv == "production" {
		// --- CONFIGURACIÓN PARA POSTGRESQL (DOCKER) ---
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port)

		dialector = postgres.Open(dsn)
		log.Println("🚀 Conectado con éxito a PostgreSQL en contenedor Docker")
	} else {
		// --- CONFIGURACIÓN PARA SQLITE (LOCAL TRADICIONAL) ---
		dialector = sqlite.Open("medicare.db")
		log.Println("💻 Conectado localmente a SQLite (medicare.db)")
	}

	// Abrir la conexión con GORM utilizando el dialector seleccionado
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
	}

	// 2. Migrar tablas dinámicamente en el motor activo
	if err := db.AutoMigrate(
		&modelsMedicacion.Medicacion{},
		&modelsMedicacion.Paciente{},
		&modelsMedicacion.HistorialMedicacion{},
		&modelsMonitoreo.CuidadorPaciente{},
		&modelsFarmacia.Farmacia{},
	); err != nil {
		log.Fatal("Error al migrar la base de datos:", err)
	}

	// 3. Obtener *sql.DB para sqlc (módulo medicación)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error obteniendo sql.DB:", err)
	}

	// 4. Repositorios
	almacenMedicacion := storage.NuevoAlmacenSQLC(sqlDB)
	almacenFarmacia := storageFarmacia.NuevoStorageFarmaciaGORM(db)
	almacenMonitoreo := storageMonitoreo.NewMonitoreoSQLite(db)

	// 5. Auth
	usuarioRepo := storage.NuevoAlmacenUsuario()
	authService := service.NewAuthService(usuarioRepo)

	// 6. Servicios medicación
	medicacionSvc := servicioMedicacion.NewMedicacionService(almacenMedicacion)
	pacienteSvc := servicioMedicacion.NewPacienteService(almacenMedicacion)
	historialSvc := servicioMedicacion.NewHistorialService(almacenMedicacion)
	medicacionHistorialSvc := servicioMedicacion.NewMedicacionHistorialService(almacenMedicacion)

	// 7. Servicios farmacia y monitoreo
	farmaciaSvc := servicioFarmacia.NuevoServicioFarmacia(almacenFarmacia)
	monitoreoSvc := servicioMonitoreo.NewMonitoreoService(almacenMonitoreo)

	// 8. Handlers por módulo
	medicacionHandler := handlersMedicacion.NewServer(
		medicacionSvc,
		pacienteSvc,
		historialSvc,
		medicacionHistorialSvc,
	)
	farmaciaHandler := handlersFarmacia.NuevoManejadorFarmacia(farmaciaSvc)
	monitoreoHandler := handlersMonitoreo.NewManejadorMonitoreo(monitoreoSvc)

	// 9. Servidor principal
	servidor := handlers.NewServer(
		medicacionHandler,
		farmaciaHandler,
		monitoreoHandler,
		authService,
	)

	// 10. Router Chi
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	r.Route("/api/v1", func(r chi.Router) {
		// Rutas públicas
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Rutas protegidas
		r.Group(func(r chi.Router) {
			r.Use(middleware.Autenticacion(authService))

			// Medicaciones
			r.Get("/medicaciones", servidor.Medicacion.ListarMedicacion)
			r.Post("/medicaciones", servidor.Medicacion.CrearMedicacion)
			r.Get("/medicaciones/{id}", servidor.Medicacion.ObtenerMedicacion)
			r.Put("/medicaciones/{id}", servidor.Medicacion.ActualizarMedicacion)
			r.Delete("/medicaciones/{id}", servidor.Medicacion.EliminarMedicacion)

			// Pacientes
			r.Get("/pacientes", servidor.Medicacion.ListarPacientes)
			r.Post("/pacientes", servidor.Medicacion.CrearPaciente)
			r.Get("/pacientes/{id}", servidor.Medicacion.BuscarPacientePorID)
			r.Put("/pacientes/{id}", servidor.Medicacion.ActualizarPaciente)
			r.Delete("/pacientes/{id}", servidor.Medicacion.EliminarPaciente)

			// Historial
			r.Get("/historial", servidor.Medicacion.ListarHistorial)
			r.Get("/historial/{id}", servidor.Medicacion.BuscarPorID)
			r.Post("/historial", servidor.Medicacion.CrearHistorial)
			r.Put("/historial/{id}", servidor.Medicacion.ActualizarHistorial)
			r.Delete("/historial/{id}", servidor.Medicacion.EliminarHistorial)

			// Rutas cruzadas
			r.Get("/pacientes/{id}/medicaciones", servidor.Medicacion.ListarMedicacionPorPaciente)
			r.Get("/pacientes/{id}/historial", servidor.Medicacion.ListarHistorialPorPaciente)

			// Farmacia
			r.Post("/farmacias", servidor.Farmacia.RegistrarFarmacia)
			r.Get("/farmacias", servidor.Farmacia.BuscarCercanas)

			// Monitoreo
			r.Get("/relaciones", servidor.Monitoreo.ObtenerRelacionesHandler)
			r.Post("/relaciones", servidor.Monitoreo.CrearRelacionHandler)
			r.Get("/relaciones/{id}", servidor.Monitoreo.ObtenerRelacionPorIDHandler)
			r.Put("/relaciones/{id}", servidor.Monitoreo.ActualizarRelacionHandler)
			r.Delete("/relaciones/{id}", servidor.Monitoreo.EliminarRelacionHandler)
		})
	})

	chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("[%s] %s\n", method, route)
		return nil
	})

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}