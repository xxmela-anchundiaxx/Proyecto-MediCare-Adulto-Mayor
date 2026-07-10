package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
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
	"proyecto-medicare-adulto-mayor/internal/models" // Asegúrate de tener Usuario aquí
)

func main() {
	var dialector gorm.Dialector
	
	// 1. Configuración de Base de Datos
	if os.Getenv("DB_ENV") == "production" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
		dialector = postgres.Open(dsn)
	} else {
		dsn := "db/medicare.db"
		_ = os.MkdirAll(filepath.Dir(dsn), 0755)
		dialector = sqlite.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("❌ No se pudo conectar a la base de datos:", err)
	}

	// 2. MIGRACIONES EN ORDEN ESTRICTO (Esto soluciona el error 42P01)
	// Primero tablas sin dependencias, luego las que tienen llaves foráneas.
	err = db.AutoMigrate(
		&models.Usuario{},
		&modelsFarmacia.Farmacia{},
		&modelsMedicacion.Paciente{}, // Migrar Paciente ANTES que Medicacion
	)
	if err == nil {
		err = db.AutoMigrate(
			&modelsMedicacion.Medicacion{},
			&modelsMedicacion.HistorialMedicacion{},
			&modelsMonitoreo.CuidadorPaciente{},
		)
	}
	
	if err != nil {
		log.Fatal("❌ Error al migrar la base de datos:", err)
	}

	// 3. Inyección de dependencias
	sqlDB, _ := db.DB()
	almacenMedicacion := storage.NuevoAlmacenSQLC(sqlDB)
	almacenFarmacia := storageFarmacia.NuevoStorageFarmaciaGORM(db)
	almacenMonitoreo := storageMonitoreo.NewMonitoreoSQLite(db)
	
	usuarioRepo := storage.NuevoAlmacenUsuario(db)
	authService := service.NewAuthService(usuarioRepo)

	// Servicios
	medicacionSvc := servicioMedicacion.NewMedicacionService(almacenMedicacion)
	pacienteSvc := servicioMedicacion.NewPacienteService(almacenMedicacion)
	historialSvc := servicioMedicacion.NewHistorialService(almacenMedicacion)
	medicacionHistorialSvc := servicioMedicacion.NewMedicacionHistorialService(almacenMedicacion)

	farmaciaSvc := servicioFarmacia.NuevoServicioFarmacia(almacenFarmacia)
	monitoreoSvc := servicioMonitoreo.NuevoServicioMonitoreo(almacenMonitoreo)

	// Handlers
	servidor := handlers.NewServer(
		handlersMedicacion.NewServer(medicacionSvc, pacienteSvc, historialSvc, medicacionHistorialSvc),
		handlersFarmacia.NuevoManejadorFarmacia(farmaciaSvc),
		handlersMonitoreo.NewManejadorMonitoreo(monitoreoSvc),
		authService,
	)

	// 4. Rutas
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
			r.Get("/farmacias/{id}", servidor.Farmacia.ObtenerPorID)    
			r.Put("/farmacias/{id}", servidor.Farmacia.ActualizarFarmacia)
			r.Delete("/farmacias/{id}", servidor.Farmacia.EliminarFarmacia)

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