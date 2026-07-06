package service

import "errors"

//Lista de errores posibles
var (
	ErrNombreVacio           = errors.New("nombre es requerido")
	ErrPrecioNegativo        = errors.New("precio no puede ser negativo")
	ErrProductoNoEncontrado  = errors.New("producto no encontrado")
	ErrNoEncontrado          = errors.New("recurso no encontrado")
	ErrEmailEnUso            = errors.New("email ya en uso")
	ErrCredencialesInvalidas = errors.New("email o contraseña")
	ErrCategoriaNoEncontrada = errors.New("categoría no encontrada")
)
