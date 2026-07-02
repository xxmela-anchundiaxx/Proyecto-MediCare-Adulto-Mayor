-- Esquema de base de datos SQLite para Medicare Adulto Mayor

CREATE TABLE IF NOT EXISTS usuarios (
    id TEXT PRIMARY KEY,
    nombre TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    rol TEXT NOT NULL, -- 'paciente', 'familiar', 'medico', 'administrador'
    fecha_creacion DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pacientes (
    id TEXT PRIMARY KEY,
    usuario_id TEXT NOT NULL,
    cuidador_id TEXT, -- usuario_id del familiar o cuidador a cargo
    edad INTEGER NOT NULL,
    grupo_sanguineo TEXT,
    alergias TEXT,
    condiciones_medicas TEXT,
    contacto_emergencia TEXT NOT NULL,
    FOREIGN KEY(usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    FOREIGN KEY(cuidador_id) REFERENCES usuarios(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS medicamentos (
    id TEXT PRIMARY KEY,
    paciente_id TEXT NOT NULL,
    nombre TEXT NOT NULL,
    descripcion TEXT,
    dosis TEXT NOT NULL, -- ej. "500mg", "1 tableta"
    frecuencia TEXT NOT NULL, -- ej. "Cada 8 horas"
    via_administracion TEXT, -- ej. "Oral", "Intravenosa"
    stock INTEGER NOT NULL DEFAULT 0,
    fecha_registro DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(paciente_id) REFERENCES pacientes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS historial_medicacion (
    id TEXT PRIMARY KEY,
    medicamento_id TEXT NOT NULL,
    paciente_id TEXT NOT NULL,
    fecha_hora DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    tomado INTEGER NOT NULL CHECK (tomado IN (0, 1)), -- 0 = No tomado, 1 = Tomado
    observaciones TEXT,
    FOREIGN KEY(medicamento_id) REFERENCES medicamentos(id) ON DELETE CASCADE,
    FOREIGN KEY(paciente_id) REFERENCES pacientes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS farmacias (
    id TEXT PRIMARY KEY,
    nombre TEXT NOT NULL,
    direccion TEXT NOT NULL,
    telefono TEXT,
    latitud REAL,
    longitud REAL
);

CREATE TABLE IF NOT EXISTS monitoreo_signos (
    id TEXT PRIMARY KEY,
    paciente_id TEXT NOT NULL,
    ritmo_cardiaco INTEGER, -- ppm
    presion_arterial TEXT, -- ej: "120/80"
    nivel_azucar REAL, -- mg/dL
    temperatura REAL, -- Celsius
    fecha_hora DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    alerta_enviada INTEGER NOT NULL CHECK (alerta_enviada IN (0, 1)) DEFAULT 0,
    observaciones TEXT,
    FOREIGN KEY(paciente_id) REFERENCES pacientes(id) ON DELETE CASCADE
);
