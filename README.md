
# Proyecto: API MediCare Adulto Mayor

**Proyecto MediCare Adulto Mayor** es una API REST desarrollada en **Go** que busca apoyar a personas adultas mayores en la gestión de sus medicaciones, considerando que muchas veces olvidan tomarlas. El sistema permite registrar medicamentos, controlar horarios de tratamiento y facilitar la participación de cuidadores en el monitoreo de los pacientes.

## Objetivo del proyecto
El sistema busca ser una herramienta práctica para:
- **Adultos mayores:** Que necesitan una forma estructurada de recordar sus medicaciones.
- **Cuidadores y familiares:** Que pueden monitorear el sistema a distancia y apoyar activamente en el cumplimiento de los tratamientos médicos.

---

## Módulos del Sistema

El proyecto está dividido en tres módulos principales para separar las responsabilidades del negocio:

### 1. Módulo de Medicación
Este módulo es el núcleo del sistema, encargado de gestionar los tratamientos que debe tomar un paciente y registrar su cumplimiento.
**Funcionalidades:**
- Registrar nuevas medicaciones asignadas a un paciente.
- Consultar medicaciones existentes y el historial médico.
- Actualizar información de dosis, frecuencia y horarios.
- Eliminar medicaciones.

### 2. Módulo de Medicamento - Farmacia
Este módulo gestiona la información de los medicamentos en el mercado y las farmacias donde pueden adquirirse.
**Funcionalidades:**
- Registrar medicamentos y farmacias cercanas.
- Consultar lista de medicamentos.
- Filtrar información por disponibilidad.
- Obtener el medicamento más económico.

### 3. Módulo de Monitoreo Familiar - Cuidador
Este módulo gestiona la relación entre cuidadores (familiares, enfermeros) y pacientes, permitiendo una red de apoyo funcional.
**Funcionalidades:**
- Crear relaciones formales entre cuidadores y pacientes.
- Consultar todas las relaciones registradas.
- Obtener detalles de una relación por su ID.
- Actualizar o eliminar relaciones.

---

## Tecnologías Utilizadas
- **Lenguaje:** Go (Golang)
- **Base de Datos:** PostgreSQL
- **Generación de Código SQL:** sqlc
- **Infraestructura:** Docker y Docker Compose

---

## Instalación y Ejecución Local

1. Clonar el repositorio:
   ```bash
   git clone [https://github.com/xxmela-anchundiaxx/Proyecto-MediCare-Adulto-Mayor.git](https://github.com/xxmela-anchundiaxx/Proyecto-MediCare-Adulto-Mayor.git)



## Equipo de Desarrollo
- **[ARAY GAÓN LISBETH DOLORES]** - Módulo Medicacion 
- **[ACOSTA ANCHUNDIA MELANIE ARIANA]** -' Módulo Farmacia
- **[]** - Módulo de Monitoreo Familiar

## Stack Tecnológico
- **Lenguaje:** Go 1.x
- **Arquitectura:** En capas (Handler → Service → Repository)
- **Persistencia:** GORM con SQLite (para Medicación/Pacientes) y almacenamiento en memoria (para Farmacias/Monitoreo).
- **Seguridad:** Autenticación centralizada mediante JWT con roles (paciente, cuidador, admin).
- **Contenedores:** Docker Multi-stage y `docker-compose`.

## Ejecución
Para levantar el sistema completo con todos sus servicios, ejecuta el siguiente comando en la raíz del proyecto:
```bash
docker-compose up --build



## Endpoints de la API

La API está protegida mediante autenticación JWT. Todas las rutas dentro del grupo protegido requieren un token válido en el header `Authorization: Bearer <token>`.

### 1. Autenticación (Público)
- `POST /auth/registrar`: Registro de nuevos usuarios (paciente/cuidador).
- `POST /auth/login`: Autenticación para obtener token de acceso.

### 2. Módulo de Medicación y Pacientes (Protegido)
**Medicaciones:**
- `GET /medicaciones`: Lista todas las medicaciones.
- `POST /medicaciones`: Registra nueva medicación.
- `GET /medicaciones/{id}`: Obtiene detalle de una medicación.
- `PUT /medicaciones/{id}`: Actualiza medicación.
- `DELETE /medicaciones/{id}`: Elimina medicación.

**Pacientes:**
- `GET /pacientes`: Lista todos los pacientes.
- `POST /pacientes`: Registra un nuevo paciente.
- `GET /pacientes/{id}`: Obtiene datos de un paciente.
- `PUT /pacientes/{id}`: Actualiza datos del paciente.
- `DELETE /pacientes/{id}`: Elimina un paciente.

**Historial:**
- `GET /historial`: Lista todo el historial de tomas.
- `POST /historial`: Registra una entrada en el historial.
- `GET /historial/{id}`: Busca registro de historial por ID.
- `PUT /historial/{id}`: Actualiza registro de historial.
- `DELETE /historial/{id}`: Elimina registro de historial.

**Consultas Cruzadas:**
- `GET /pacientes/{id}/medicaciones`: Lista medicaciones de un paciente específico.
- `GET /pacientes/{id}/historial`: Lista historial de un paciente específico.

### 3. Módulo de Farmacias (Protegido)
- `POST /farmacias`: Registra una nueva farmacia.
- `GET /farmacias`: Busca farmacias (cercanas).
- `GET /farmacias/{id}`: Obtiene detalles de una farmacia.
- `PUT /farmacias/{id}`: Actualiza datos de farmacia.
- `DELETE /farmacias/{id}`: Elimina una farmacia.

### 4. Módulo de Monitoreo (Protegido)
- `POST /relaciones`: Crea relación cuidador-paciente.
- `GET /relaciones`: Lista todas las relaciones.
- `GET /relaciones/{id}`: Obtiene relación por ID.
- `PUT /relaciones/{id}`: Actualiza relación.
- `DELETE /relaciones/{id}`: Elimina relación.
