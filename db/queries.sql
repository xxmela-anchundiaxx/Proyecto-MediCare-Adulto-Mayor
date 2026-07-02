-- name: GetUserByUsername :one
SELECT * FROM usuarios
WHERE email = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO usuarios (id, nombre, email, password_hash, rol)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetPacienteByUsuarioID :one
SELECT * FROM pacientes
WHERE usuario_id = ? LIMIT 1;

-- name: CreatePaciente :one
INSERT INTO pacientes (id, usuario_id, cuidador_id, edad, grupo_sanguineo, alergias, condiciones_medicas, contacto_emergencia)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListPacientesByCuidador :many
SELECT * FROM pacientes
WHERE cuidador_id = ?;

-- name: GetMedicamentoByID :one
SELECT * FROM medicamentos
WHERE id = ? LIMIT 1;

-- name: ListMedicamentosByPaciente :many
SELECT * FROM medicamentos
WHERE paciente_id = ?;

-- name: CreateMedicamento :one
INSERT INTO medicamentos (id, paciente_id, nombre, descripcion, dosis, frecuencia, via_administracion, stock)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateMedicamentoStock :exec
UPDATE medicamentos
SET stock = ?
WHERE id = ?;

-- name: CreateHistorialMedicacion :one
INSERT INTO historial_medicacion (id, medicamento_id, paciente_id, fecha_hora, tomado, observaciones)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListHistorialMedicacion :many
SELECT h.*, m.nombre AS medicamento_nombre
FROM historial_medicacion h
JOIN medicamentos m ON h.medicamento_id = m.id
WHERE h.paciente_id = ?
ORDER BY h.fecha_hora DESC;

-- name: CreateMonitoreoSignos :one
INSERT INTO monitoreo_signos (id, paciente_id, ritmo_cardiaco, presion_arterial, nivel_azucar, temperatura, alerta_enviada, observaciones)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListMonitoreoSignos :many
SELECT * FROM monitoreo_signos
WHERE paciente_id = ?
ORDER BY fecha_hora DESC;
