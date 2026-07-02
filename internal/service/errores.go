package service

import "errors"

var (
	ErrUsuarioExistente     = errors.New("el usuario con este correo ya se encuentra registrado")
	ErrCredencialesInvalidas = errors.New("correo o contraseña incorrectos")
	ErrNoEncontrado          = errors.New("recurso no encontrado")
	ErrPermisoDenegado       = errors.New("no tienes permisos para realizar esta acción")
	ErrDatosInvalidos        = errors.New("los datos ingresados no son válidos")
)
