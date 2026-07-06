package storage

import (
    "context"
    "database/sql"

    "proyecto-medicare-adulto-mayor/internal/models/farmacia"
    "proyecto-medicare-adulto-mayor/internal/models/medicacion"
    "proyecto-medicare-adulto-mayor/internal/models/monitoreo"
    "proyecto-medicare-adulto-mayor/internal/storage/sqlcdb"
)

type AlmacenSQLC struct {
    q *sqlcdb.Queries
}

func NuevoAlmacenSQLC(db *sql.DB) *AlmacenSQLC {
    return &AlmacenSQLC{q: sqlcdb.New(db)}
}

// ── MAPEOS ───────────────────────────────────────────────────────────────────

func aPacienteDominio(p sqlcdb.Pacientes) medicacion.Paciente {
    return medicacion.Paciente{ID: int(p.ID), Nombre: p.Nombre, Edad: int(p.Edad)}
}

func aMedicacionDominio(m sqlcdb.Medicaciones) medicacion.Medicacion {
    return medicacion.Medicacion{
        ID: int(m.ID), Nombre: m.Nombre, Descripcion: m.Descripcion,
        Dosis: m.Dosis, Frecuencia: m.Frecuencia,
        Hora_programada: m.HoraProgramada, Inicio_tratamiento: m.InicioTratamiento,
        Fecha_creacion: m.FechaCreacion, PacienteID: int(m.PacienteID),
    }
}

func aHistorialDominio(h sqlcdb.HistorialMedicaciones) medicacion.HistorialMedicacion {
    return medicacion.HistorialMedicacion{
        ID: int(h.ID), MedicacionID: int(h.MedicacionID),
        FechaHora: h.FechaHora, Tomada: h.Tomada, Observacion: h.Observacion,
    }
}

func aFarmaciaDominio(f sqlcdb.Farmacias) farmacia.Farmacia {
    return farmacia.Farmacia{
        ID: f.ID, Nombre: f.Nombre, Direccion: f.Direccion,
        Telefono: f.Telefono, Latitud: f.Latitud, Longitud: f.Longitud,
    }
}

func aCuidadorPacienteDominio(c sqlcdb.CuidadoresPacientes) monitoreo.CuidadorPaciente {
    return monitoreo.CuidadorPaciente{
        ID: int(c.ID), CuidadorID: int(c.CuidadorID),
        PacienteID: int(c.PacienteID), Relacion: c.Relacion,
    }
}

// ── PACIENTES ────────────────────────────────────────────────────────────────

func (a *AlmacenSQLC) ListarPacientes() ([]medicacion.Paciente, error) {
    filas, err := a.q.ListarPacientes(context.Background())
    if err != nil {
        return nil, err
    }
    out := make([]medicacion.Paciente, 0, len(filas))
    for _, f := range filas {
        out = append(out, aPacienteDominio(f))
    }
    return out, nil
}

func (a *AlmacenSQLC) BuscarPacientePorID(id int) (medicacion.Paciente, error) {
    f, err := a.q.ObtenerPacientePorID(context.Background(), int64(id))
    if err != nil {
        return medicacion.Paciente{}, err
    }
    return aPacienteDominio(f), nil
}

func (a *AlmacenSQLC) CrearPaciente(p medicacion.Paciente) (medicacion.Paciente, error) {
    f, err := a.q.CrearPaciente(context.Background(), sqlcdb.CrearPacienteParams{
        Nombre: p.Nombre, Edad: int64(p.Edad),
    })
    if err != nil {
        return medicacion.Paciente{}, err
    }
    return aPacienteDominio(f), nil
}

func (a *AlmacenSQLC) ActualizarPaciente(id int, datos medicacion.Paciente) (medicacion.Paciente, error) {
    f, err := a.q.ActualizarPaciente(context.Background(), sqlcdb.ActualizarPacienteParams{
        Nombre: datos.Nombre, Edad: int64(datos.Edad), ID: int64(id),
    })
    if err != nil {
        return medicacion.Paciente{}, err
    }
    return aPacienteDominio(f), nil
}

func (a *AlmacenSQLC) EliminarPaciente(id int) error {
    return a.q.EliminarPaciente(context.Background(), int64(id))
}

// ── MEDICACIONES ─────────────────────────────────────────────────────────────

func (a *AlmacenSQLC) ListarMedicacion() ([]medicacion.Medicacion, error) {
    filas, err := a.q.ListarMedicaciones(context.Background())
    if err != nil {
        return nil, err
    }
    out := make([]medicacion.Medicacion, 0, len(filas))
    for _, f := range filas {
        out = append(out, aMedicacionDominio(f))
    }
    return out, nil
}

func (a *AlmacenSQLC) BuscarMedicacionPorID(id int) (medicacion.Medicacion, error) {
    f, err := a.q.ObtenerMedicacionPorID(context.Background(), int64(id))
    if err != nil {
        return medicacion.Medicacion{}, err
    }
    return aMedicacionDominio(f), nil
}

func (a *AlmacenSQLC) CrearMedicacion(m medicacion.Medicacion) (medicacion.Medicacion, error) {
    f, err := a.q.CrearMedicacion(context.Background(), sqlcdb.CrearMedicacionParams{
        Nombre: m.Nombre, Descripcion: m.Descripcion, Dosis: m.Dosis,
        Frecuencia: m.Frecuencia, HoraProgramada: m.Hora_programada,
        InicioTratamiento: m.Inicio_tratamiento, FechaCreacion: m.Fecha_creacion,
        PacienteID: int64(m.PacienteID),
    })
    if err != nil {
        return medicacion.Medicacion{}, err
    }
    return aMedicacionDominio(f), nil
}

func (a *AlmacenSQLC) ActualizarMedicacion(id int, datos medicacion.Medicacion) (medicacion.Medicacion, error) {
    f, err := a.q.ActualizarMedicacion(context.Background(), sqlcdb.ActualizarMedicacionParams{
        Nombre: datos.Nombre, Descripcion: datos.Descripcion, Dosis: datos.Dosis,
        Frecuencia: datos.Frecuencia, HoraProgramada: datos.Hora_programada,
        InicioTratamiento: datos.Inicio_tratamiento, FechaCreacion: datos.Fecha_creacion,
        PacienteID: int64(datos.PacienteID), ID: int64(id),
    })
    if err != nil {
        return medicacion.Medicacion{}, err
    }
    return aMedicacionDominio(f), nil
}

func (a *AlmacenSQLC) EliminarMedicacion(id int) error {
    return a.q.EliminarMedicacion(context.Background(), int64(id))
}

func (a *AlmacenSQLC) ListarMedicacionPorPaciente(pacienteID int) ([]medicacion.Medicacion, error) {
    filas, err := a.q.ListarMedicacionesPorPaciente(context.Background(), int64(pacienteID))
    if err != nil {
        return nil, err
    }
    out := make([]medicacion.Medicacion, 0, len(filas))
    for _, f := range filas {
        out = append(out, aMedicacionDominio(f))
    }
    return out, nil
}

// ── HISTORIAL ────────────────────────────────────────────────────────────────

func (a *AlmacenSQLC) ListarHistorial() ([]medicacion.HistorialMedicacion, error) {
    filas, err := a.q.ListarHistorial(context.Background())
    if err != nil {
        return nil, err
    }
    out := make([]medicacion.HistorialMedicacion, 0, len(filas))
    for _, f := range filas {
        out = append(out, aHistorialDominio(f))
    }
    return out, nil
}

func (a *AlmacenSQLC) BuscarHistorialPorID(id int) (medicacion.HistorialMedicacion, error) {
    f, err := a.q.ObtenerHistorialPorID(context.Background(), int64(id))
    if err != nil {
        return medicacion.HistorialMedicacion{}, err
    }
    return aHistorialDominio(f), nil
}

func (a *AlmacenSQLC) CrearHistorial(h medicacion.HistorialMedicacion) (medicacion.HistorialMedicacion, error) {
    f, err := a.q.CrearHistorial(context.Background(), sqlcdb.CrearHistorialParams{
        MedicacionID: int64(h.MedicacionID), FechaHora: h.FechaHora,
        Tomada: h.Tomada, Observacion: h.Observacion,
    })
    if err != nil {
        return medicacion.HistorialMedicacion{}, err
    }
    return aHistorialDominio(f), nil
}

func (a *AlmacenSQLC) ActualizarHistorial(id int, datos medicacion.HistorialMedicacion) (medicacion.HistorialMedicacion, error) {
    f, err := a.q.ActualizarHistorial(context.Background(), sqlcdb.ActualizarHistorialParams{
        MedicacionID: int64(datos.MedicacionID), FechaHora: datos.FechaHora,
        Tomada: datos.Tomada, Observacion: datos.Observacion, ID: int64(id),
    })
    if err != nil {
        return medicacion.HistorialMedicacion{}, err
    }
    return aHistorialDominio(f), nil
}

func (a *AlmacenSQLC) EliminarHistorial(id int) error {
    return a.q.EliminarHistorial(context.Background(), int64(id))
}

func (a *AlmacenSQLC) ListarHistorialPorPaciente(pacienteID int) ([]medicacion.HistorialMedicacion, error) {
    filas, err := a.q.ListarHistorialPorPaciente(context.Background(), int64(pacienteID))
    if err != nil {
        return nil, err
    }
    out := make([]medicacion.HistorialMedicacion, 0, len(filas))
    for _, f := range filas {
        out = append(out, aHistorialDominio(f))
    }
    return out, nil
}