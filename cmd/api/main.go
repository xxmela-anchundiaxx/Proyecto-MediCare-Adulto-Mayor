package main

import (
	"fmt"
	"net/http"

	"proyectodemedicare/internal/handlers"

	"github.com/go-chi/chi/v5"
)
//funciones del proyecto
func main() {
	// Crea una nueva instancia del router Chi
	r := chi.NewRouter()

	r.Post("/api/v1/farmacia", handlers.CreateMedicamento)
	r.Get("/api/v1/farmacia", handlers.GetMedicamentos)
	r.Get("/api/v1/farmacia/{id}", handlers.GetMedicamento)
	r.Put("/api/v1/farmacia/{id}", handlers.UpdateMedicamento)
	r.Delete("/api/v1/farmacia/{id}", handlers.DeleteMedicamento)
	fmt.Println("Servidor iniciado en http://localhost:8080")

	http.ListenAndServe(":8080", r)
}
