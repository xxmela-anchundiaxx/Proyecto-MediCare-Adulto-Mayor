# Medicare Adulto Mayor - Backend REST API (Go)

Este es el backend del proyecto **Medicare Adulto Mayor**, un sistema integral diseñado para facilitar el cuidado, monitoreo de signos vitales, adherencia a medicamentos, y búsqueda de farmacias cercanas para adultos mayores, integrando a familiares/cuidadores y médicos.

## 🏗️ Arquitectura Corregida y Modular

El proyecto ha sido completamente reestructurado y corregido bajo una **Arquitectura en Capas (Layered Clean Architecture)** limpia, escalable y mantenible. Se ha removido el acoplamiento y se ha distribuido el código según sus responsabilidades:

```text
Proyecto-MediCare-Adulto-Mayor/
├── cmd/
│   └── api/
│       └── main.go         # Punto de entrada de la aplicación. Configura DB e inicializa el servidor.
├── db/
│   ├── schema.sql          # Esquema relacional DDL (SQLite) para usuarios, pacientes, medicamentos, etc.
│   └── queries.sql         # Consultas SQL CRUD preparadas para Sqlc.
├── internal/
│   ├── models/             # Entidades del dominio (Usuario, Paciente, Medicamento, Signos).
│   ├── storage/            # Capa de Persistencia (SQLite para datos críticos, Memoria para Farmacias).
│   ├── service/            # Capa de Negocio / Casos de Uso (Validación, control de umbrales críticos, lógica).
│   ├── handlers/           # Capa de Presentación (Controladores HTTP, serialización JSON).
│   └── middleware/         # Filtros HTTP (CORS, Autenticación Bearer Token por Contexto).
├── go.mod                  # Declaración de dependencias (Go 1.22+).
├── sqlc.yaml               # Configuración del compilador de queries SQL a Go (Sqlc).
└── README.md               # Este archivo guía.
```

---

## 🚀 Requisitos para Correr en Local

Asegúrate de tener instalado en tu máquina local:
- **Go 1.22.0** o superior.
- **GCC** (necesario para compilar el driver SQLite CGO `go-sqlite3` en local).

### Instalación y Compilación

1. Entra al directorio del proyecto:
   ```bash
   cd Proyecto-MediCare-Adulto-Mayor
   ```

2. Descarga y sincroniza las dependencias de Go:
   ```bash
   go mod tidy
   ```

3. Compila el ejecutable del backend:
   ```bash
   go build -o medicare-api cmd/api/main.go
   ```

4. Ejecuta el servidor:
   ```bash
   ./medicare-api
   ```

Por defecto, la API levantará un servidor web en `http://localhost:8080` y creará automáticamente un archivo de base de datos relacional ligero `medicare.db` sincronizando todas las tablas declaradas en `db/schema.sql`.

---

## 🛡️ Endpoints y Guía de Uso de la API

La API cuenta con endpoints públicos para la autenticación y endpoints protegidos por el middleware de seguridad que validan un encabezado `Authorization: Bearer <usuario_id>`.

### 1. Autenticación (Público)

#### 📝 Registro de Usuario
Permite registrar usuarios con roles específicos: `paciente` (adulto mayor), `familiar` (cuidador), `medico` o `administrador`.

- **Endpoint**: `POST /api/v1/auth/registro`
- **Cuerpo (JSON)**:
```json
{
  "nombre": "Juan Pérez",
  "email": "juan.perez@email.com",
  "password": "mi_password_segura",
  "rol": "paciente"
}
```

#### 🔑 Iniciar Sesión (Login)
Valida credenciales y retorna los datos del usuario junto a un **token de autenticación** (para simplificar el flujo, se usa el ID del usuario como token Bearer).

- **Endpoint**: `POST /api/v1/auth/login`
- **Cuerpo (JSON)**:
```json
{
  "email": "juan.perez@email.com",
  "password": "mi_password_segura"
}
```
- **Respuesta (JSON)**:
```json
{
  "id": "e2a4f4bd-44c1-4545-9856-7871e8bf28a4",
  "nombre": "Juan Pérez",
  "email": "juan.perez@email.com",
  "rol": "paciente",
  "token": "e2a4f4bd-44c1-4545-9856-7871e8bf28a4"
}
```

---

### 2. Pacientes (Protegido)
Requiere añadir en los headers: `Authorization: Bearer <usuario_id>`

#### ➕ Registrar Perfil Clínico de un Paciente
Registra la edad, grupo sanguíneo, alergias, condiciones médicas y contactos de emergencia del adulto mayor.

- **Endpoint**: `POST /api/v1/pacientes`
- **Headers**: `Authorization: Bearer <usuario_id>`
- **Cuerpo (JSON)**:
```json
{
  "usuario_id": "e2a4f4bd-44c1-4545-9856-7871e8bf28a4",
  "edad": 78,
  "grupo_sanguineo": "O+",
  "alergias": "Penicilina",
  "condiciones_medicas": "Hipertensión leve, Diabetes Tipo 2",
  "contacto_emergencia": "+56911112222 (Hija Maria)"
}
```

---

### 3. Medicamentos y Adherencia (Protegido)

#### 💊 Registrar un Medicamento (Prescripción)
Añade un medicamento recetado especificando la dosis, frecuencia, vía de administración y el stock actual.

- **Endpoint**: `POST /api/v1/medicamentos`
- **Cuerpo (JSON)**:
```json
{
  "paciente_id": "id-del-paciente-creado",
  "nombre": "Metformina",
  "descripcion": "Tomar junto con el desayuno",
  "dosis": "850mg",
  "frecuencia": "Cada 24 horas",
  "via_administracion": "Oral",
  "stock": 30
}
```

#### 🗓️ Registrar Toma de Medicamento (Adherencia)
Registra si el adulto mayor tomó su dosis correspondiente. **¡Logica Corregida!** Al registrar que la dosis fue tomada (`tomado: true`), el backend descuenta automáticamente `1` unidad del stock disponible del medicamento.

- **Endpoint**: `POST /api/v1/historial`
- **Cuerpo (JSON)**:
```json
{
  "medicamento_id": "id-del-medicamento-creado",
  "paciente_id": "id-del-paciente-creado",
  "tomado": true,
  "observaciones": "Tomada a la hora correcta sin malestares"
}
```

---

### 4. Monitoreo Inteligente de Signos Vitales (Protegido)

#### 💓 Registrar Mediciones y Alertas de Vía
Permite registrar mediciones de ritmo cardíaco, presión arterial, azúcar o temperatura. El servicio analiza automáticamente si las mediciones representan un riesgo (fiebre, hipertensión, hipoglucemia, arritmias) y, en caso afirmativo, marca la medición con `alerta_enviada: true` y genera observaciones de alerta automáticas para avisar inmediatamente al cuidador.

- **Endpoint**: `POST /api/v1/monitoreo`
- **Cuerpo de Medición Normal (Sin Alerta)**:
```json
{
  "paciente_id": "id-del-paciente-creado",
  "ritmo_cardiaco": 72,
  "presion_arterial": "120/80",
  "nivel_azucar": 95,
  "temperatura": 36.6,
  "observaciones": "Control matutino ordinario"
}
```

- **Cuerpo de Medición de Riesgo (Desencadena Alerta)**:
```json
{
  "paciente_id": "id-del-paciente-creado",
  "ritmo_cardiaco": 115, // Arritmia / Taquicardia (>105 ppm)
  "presion_arterial": "150/95", // Hipertensión (>140/90)
  "nivel_azucar": 62, // Hipoglucemia (<70 mg/dL)
  "temperatura": 38.5, // Fiebre (>37.8 °C)
  "observaciones": ""
}
```
*La respuesta retornará el objeto guardado con la bandera `alerta_enviada: true` y observaciones automáticas cargadas.*

---

### 5. Farmacias y Coordenadas (Protegido)

#### 📍 Buscar Farmacias Cercanas por Geolocalización
Permite buscar farmacias registradas en un radio específico de kilómetros utilizando cálculo de distancia de Haversine.

- **Endpoint**: `GET /api/v1/farmacias?lat=-33.456&lon=-70.648&radio_km=5`
- **Respuesta (JSON)**:
```json
[
  {
    "id": "f1",
    "nombre": "Farmacia Medicare Centro",
    "direccion": "Av. Principal 123, Centro",
    "telefono": "+56 9 1234 5678",
    "latitud": -33.456,
    "longitud": -70.648
  }
]
```

---

## 🛠️ Buenas Prácticas de Go Aplicadas

1. **Uso Exclusivo de Biblioteca Estándar para Ruteo**: Aprovecha las mejoras de ruteo nativas de `net/http` en Go 1.22+, evitando cargar dependencias innecesarias o vulnerables.
2. **Contextos Seguros**: Uso del tipo `contextKey` privado y seguro para evitar colisiones de datos al propagar el `usuario_id` en peticiones autenticadas.
3. **Control Concurrente**: Implementación de `sync.RWMutex` (bloqueos de lectura/escritura) en los repositorios de memoria para evitar carreras de datos (*race conditions*).
4. **Manejo de Errores Tipificados**: Centralización de los errores de negocio para responder con códigos de estado HTTP semánticos y coherentes (`StatusUnauthorized`, `StatusNotFound`, `StatusConflict`).
