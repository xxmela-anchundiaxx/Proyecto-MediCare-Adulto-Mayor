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
    // Ruta para registrar un nuevo medicamento
	r.Post("/api/v1/farmacia", handlers.CreateMedicamento)
	// Ruta para obtener la lista de todos los medicamentos
	r.Get("/api/v1/farmacia", handlers.GetMedicamentos)
	// Ruta para obtener un medicamento específico mediante su ID
	r.Get("/api/v1/farmacia/{id}", handlers.GetMedicamento)
	// Ruta para actualizar la información de un medicamento por ID
	r.Put("/api/v1/farmacia/{id}", handlers.UpdateMedicamento)
	// Ruta para eliminar un medicamento por ID
	r.Delete("/api/v1/farmacia/{id}", handlers.DeleteMedicamento)
	fmt.Println("Servidor iniciado en http://localhost:8080")

	http.ListenAndServe(":8080", r)
}
