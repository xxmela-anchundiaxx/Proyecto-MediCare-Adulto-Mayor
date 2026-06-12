# Proyecto MediCare Adulto Mayor

API REST desarrollada en **Go** para la gestión de dos módulos:
- **Medicaciones** (almacenadas en SQLite con GORM).
- **Medicamentos de Farmacia** (almacenados en memoria).
- **Monitorio familiar o cuidador (almacenados en memoria)

---

## Instalación
1. Clonar el repositorio:
   git clone https://github.com/xxmela-anchundiaxx/Proyecto-MediCare-Adulto-Mayor.git
  


Descripción general del proyecto
Proyecto MediCare Adulto Mayor es una API REST desarrollada en Go que busca apoyar a personas adultas mayores en la gestión de sus medicaciones, considerando que muchas veces olvidan tomarlas. El sistema permite registrar medicamentos, controlar horarios de tratamiento y facilitar la participación de cuidadores en el monitoreo de los pacientes.

El proyecto está dividido en tres módulos principales:

--MÓDULO MEDICACION
Este módulo gestiona las medicaciones que debe tomar un paciente.
Funcionalidad:
Registrar nuevas medicaciones.
Consultar medicaciones existentes.
Actualizar información de dosis y frecuencia.
Eliminar medicaciones.


MÓDULO MEDICAMENTO FARMACIA
Este módulo gestiona los medicamentos disponibles en farmacias, almacenados en memoria.
Funcionalidad:
Registrar medicamentos disponibles en farmacias.
Consultar lista de medicamentos.
Filtrar por disponibilidad.
Obtener el medicamento más económico.


MÓDULO MONITOREO FAMILIAR-CUIDADOR
Este módulo gestiona la relación entre cuidadores y pacientes, permitiendo que los cuidadores participen activamente en el control de las medicaciones.
Funcionalidad:
Crear relaciones entre cuidadores y pacientes.
Consultar todas las relaciones.
Obtener relación por ID.
Actualizar o eliminar relaciones.


.Objetivo del proyecto
El sistema busca ser una herramienta práctica para:
- Adultos mayores, que necesitan recordar sus medicaciones.
- Cuidadores y familiares, que pueden monitorear y apoyar en el cumplimiento de los tratamientos.


