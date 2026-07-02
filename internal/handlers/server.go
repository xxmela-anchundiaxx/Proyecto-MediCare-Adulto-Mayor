package handlers

import (
	"fmt"
	"log"
	"medicare-adulto-mayor/internal/middleware"
	farmaciaHandler "medicare-adulto-mayor/internal/handlers/farmacia"
	medicacionHandler "medicare-adulto-mayor/internal/handlers/medicacion"
	monitoreoHandler "medicare-adulto-mayor/internal/handlers/monitoreo"
	farmaciaService "medicare-adulto-mayor/internal/service/farmacia"
	medicacionService "medicare-adulto-mayor/internal/service/medicacion"
	monitoreoService "medicare-adulto-mayor/internal/service/monitoreo"
	farmaciaStorage "medicare-adulto-mayor/internal/storage/farmacia"
	medicacionStorage "medicare-adulto-mayor/internal/storage/medicacion"
	monitoreoStorage "medicare-adulto-mayor/internal/storage/monitoreo"
	authService "medicare-adulto-mayor/internal/service"
	authStorage "medicare-adulto-mayor/internal/storage"
	"net/http"
)

type ServidorMedicare struct {
	Router  *http.ServeMux
	Almacen *authStorage.Almacen
}

func NuevoServidorMedicare(almacen *authStorage.Almacen) *ServidorMedicare {
	s := &ServidorMedicare{
		Router:  http.NewServeMux(),
		Almacen: almacen,
	}
	s.configurarRutas()
	return s
}

func (s *ServidorMedicare) configurarRutas() {
	// 1. Inicializar Capa de Almacenamiento (Storages)
	storageUsuario := authStorage.NuevoStorageUsuario(s.Almacen.DB)
	storagePaciente := medicacionStorage.NuevoStoragePacienteGORM(s.Almacen.GORM)
	storageMedicamento := medicacionStorage.NuevoStorageMedicamentoGORM(s.Almacen.GORM)
	storageHistorial := medicacionStorage.NuevoStorageHistorialGORM(s.Almacen.GORM)
	storageFarmacia := farmaciaStorage.NuevoStorageFarmaciaGORM(s.Almacen.GORM)
	storageMonitoreo := monitoreoStorage.NuevoStorageMonitoreoGORM(s.Almacen.GORM)


	// 2. Inicializar Capa de Negocio (Servicios)
	servicioAuth := authService.NuevoServicioAuth(storageUsuario)
	servicioPaciente := medicacionService.NuevoServicioPaciente(storagePaciente)
	servicioMedicamento := medicacionService.NuevoServicioMedicamento(storageMedicamento)
	servicioHistorial := medicacionService.NuevoServicioHistorial(storageHistorial, storageMedicamento)
	servicioFarmacia := farmaciaService.NuevoServicioFarmacia(storageFarmacia)
	servicioMonitoreo := monitoreoService.NuevoServicioMonitoreo(storageMonitoreo)

	// 3. Inicializar Capa de Presentación (Controladores / Handlers)
	manejadorAuth := NuevoManejadorAuth(servicioAuth)
	manejadorPaciente := medicacionHandler.NuevoManejadorPaciente(servicioPaciente)
	manejadorMedicamento := medicacionHandler.NuevoManejadorMedicamento(servicioMedicamento)
	manejadorHistorial := medicacionHandler.NuevoManejadorHistorial(servicioHistorial)
	manejadorFarmacia := farmaciaHandler.NuevoManejadorFarmacia(servicioFarmacia)
	manejadorMonitoreo := monitoreoHandler.NuevoManejadorMonitoreo(servicioMonitoreo)

	// 4. Registrar Rutas Públicas (Sin Autenticación)
	s.Router.HandleFunc("POST /api/v1/auth/registro", manejadorAuth.Registro)
	s.Router.HandleFunc("POST /api/v1/auth/login", manejadorAuth.Login)

	// Rutas Protegidas (Requieren cabecera 'Authorization: Bearer <usuario_id>')
	// Creamos un sub-mux o handler para aplicar el Middleware
	apiMux := http.NewServeMux()

	// Pacientes
	apiMux.HandleFunc("POST /pacientes", manejadorPaciente.RegistrarPaciente)
	apiMux.HandleFunc("GET /pacientes", manejadorPaciente.ObtenerPaciente)
	apiMux.HandleFunc("GET /pacientes/cuidador", manejadorPaciente.ListarPorCuidador)

	// Medicamentos
	apiMux.HandleFunc("POST /medicamentos", manejadorMedicamento.RegistrarMedicamento)
	apiMux.HandleFunc("GET /medicamentos", manejadorMedicamento.ListarPorPaciente)

	// Historial / Adherencia
	apiMux.HandleFunc("POST /historial", manejadorHistorial.RegistrarToma)
	apiMux.HandleFunc("GET /historial", manejadorHistorial.ListarPorPaciente)

	// Farmacias
	apiMux.HandleFunc("POST /farmacias", manejadorFarmacia.RegistrarFarmacia)
	apiMux.HandleFunc("GET /farmacias", manejadorFarmacia.BuscarCercanas)

	// Monitoreo de Signos
	apiMux.HandleFunc("POST /monitoreo", manejadorMonitoreo.RegistrarSignos)
	apiMux.HandleFunc("GET /monitoreo", manejadorMonitoreo.ListarPorPaciente)

	// Integración de Mux Protegido bajo AuthMiddleware
	s.Router.Handle("/api/v1/", http.StripPrefix("/api/v1", middleware.AuthMiddleware(apiMux)))
}

// Iniciar arranca el servidor HTTP en el puerto indicado con soporte para CORS
func (s *ServidorMedicare) Iniciar(puerto string) {
	// Aplicar CORS a toda la aplicación
	handlerConCORS := middleware.CORSMiddleware(s.Router)

	log.Printf("Servidor REST de Medicare Adulto Mayor corriendo en http://localhost:%s", puerto)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", puerto), handlerConCORS); err != nil {
		log.Fatalf("Error al arrancar el servidor: %v", err)
	}
}
