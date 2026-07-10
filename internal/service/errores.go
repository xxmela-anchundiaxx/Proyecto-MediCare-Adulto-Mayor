package service

import "errors"

var (
	ErrNombreVacio           = errors.New("nombre es requerido")
	ErrPrecioNegativo        = errors.New("precio no puede ser negativo")
	ErrProductoNoEncontrado  = errors.New("producto no encontrado")
	ErrNoEncontrado          = errors.New("recurso no encontrado")
	ErrEmailEnUso            = errors.New("email ya en uso")
	ErrEmailVacio            = errors.New("email o password vacíos")
	ErrCredencialesInvalidas = errors.New("credenciales inválidas")
	ErrCategoriaNoEncontrada = errors.New("categoría no encontrada")
)