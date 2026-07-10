package medicacion

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	models "proyecto-medicare-adulto-mayor/internal/models/medicacion"
	svc "proyecto-medicare-adulto-mayor/internal/service/medicacion"
	"proyecto-medicare-adulto-mayor/internal/storage"
)

// ---------- mocks de repositorio (interfaces reales de storage) ----------

type mockMedRepo struct {
	listar     func() ([]models.Medicacion, error)
	obtener    func(int) (models.Medicacion, error)
	crear      func(models.Medicacion) (models.Medicacion, error)
	actualizar func(int, models.Medicacion) (models.Medicacion, error)
	eliminar   func(int) error
}

func (r *mockMedRepo) ListarMedicacion() ([]models.Medicacion, error)        { return r.listar() }
func (r *mockMedRepo) BuscarMedicacionPorID(id int) (models.Medicacion, error) { return r.obtener(id) }
func (r *mockMedRepo) CrearMedicacion(m models.Medicacion) (models.Medicacion, error) { return r.crear(m) }
func (r *mockMedRepo) ActualizarMedicacion(id int, m models.Medicacion) (models.Medicacion, error) {
	return r.actualizar(id, m)
}
func (r *mockMedRepo) EliminarMedicacion(id int) error { return r.eliminar(id) }

var _ storage.MedicacionRepository = (*mockMedRepo)(nil)

type mockHistRepo struct {
	listar     func() ([]models.HistorialMedicacion, error)
	obtener    func(int) (models.HistorialMedicacion, error)
	crear      func(models.HistorialMedicacion) (models.HistorialMedicacion, error)
	actualizar func(int, models.HistorialMedicacion) (models.HistorialMedicacion, error)
	eliminar   func(int) error
}

func (r *mockHistRepo) ListarHistorial() ([]models.HistorialMedicacion, error) { return r.listar() }
func (r *mockHistRepo) BuscarHistorialPorID(id int) (models.HistorialMedicacion, error) {
	return r.obtener(id)
}
func (r *mockHistRepo) CrearHistorial(h models.HistorialMedicacion) (models.HistorialMedicacion, error) {
	return r.crear(h)
}
func (r *mockHistRepo) ActualizarHistorial(id int, h models.HistorialMedicacion) (models.HistorialMedicacion, error) {
	return r.actualizar(id, h)
}
func (r *mockHistRepo) EliminarHistorial(id int) error { return r.eliminar(id) }

var _ storage.HistorialRepository = (*mockHistRepo)(nil)

type mockMedHistRepo struct {
	medPorPaciente  func(int) ([]models.Medicacion, error)
	histPorPaciente func(int) ([]models.HistorialMedicacion, error)
}

func (r *mockMedHistRepo) ListarMedicacionPorPaciente(id int) ([]models.Medicacion, error) {
	return r.medPorPaciente(id)
}
func (r *mockMedHistRepo) ListarHistorialPorPaciente(id int) ([]models.HistorialMedicacion, error) {
	return r.histPorPaciente(id)
}

var _ storage.MedicacionHistorialRepository = (*mockMedHistRepo)(nil)

// ---------- helpers ----------

// run construye la request (con param de ruta chi opcional), ejecuta el
// handler y devuelve el recorder.
func run(handler http.HandlerFunc, method, url, body, paramKey, paramVal string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, url, bytes.NewBufferString(body))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	if paramKey != "" {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add(paramKey, paramVal)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

func check(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rr.Code != want {
		t.Fatalf("esperado status %d, obtenido %d: %s", want, rr.Code, rr.Body.String())
	}
}

func medValida() models.Medicacion {
	return models.Medicacion{
		Nombre:             "Losartán",
		Dosis:              "50mg",
		Frecuencia:         "cada 12 horas",
		Hora_programada:    "08:00",
		Inicio_tratamiento: time.Now(),
	}
}

// ---------- Medicación ----------

func TestMedicacionHandlers(t *testing.T) {
	t.Run("Listar", func(t *testing.T) {
		for _, c := range []struct {
			name string
			fn   func() ([]models.Medicacion, error)
			want int
		}{
			{"exito", func() ([]models.Medicacion, error) { return []models.Medicacion{medValida()}, nil }, http.StatusOK},
			{"error", func() ([]models.Medicacion, error) { return nil, errors.New("db") }, http.StatusInternalServerError},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{Medicacion: svc.NewMedicacionService(&mockMedRepo{listar: c.fn})}
				check(t, run(s.ListarMedicacion, http.MethodGet, "/medicaciones", "", "", ""), c.want)
			})
		}
	})

	t.Run("Obtener", func(t *testing.T) {
		for _, c := range []struct {
			name, id string
			fn       func(int) (models.Medicacion, error)
			want     int
		}{
			{"exito", "1", func(int) (models.Medicacion, error) { return medValida(), nil }, http.StatusOK},
			{"id_invalido", "abc", nil, http.StatusBadRequest},
			{"no_encontrada", "99", func(int) (models.Medicacion, error) { return models.Medicacion{}, gorm.ErrRecordNotFound }, http.StatusNotFound},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{Medicacion: svc.NewMedicacionService(&mockMedRepo{obtener: c.fn})}
				check(t, run(s.ObtenerMedicacion, http.MethodGet, "/medicaciones/"+c.id, "", "id", c.id), c.want)
			})
		}
	})

	t.Run("Crear", func(t *testing.T) {
		exito := func(m models.Medicacion) (models.Medicacion, error) { return m, nil }
		validBody, _ := json.Marshal(medValida())

		for _, c := range []struct {
			name, body string
			fn         func(models.Medicacion) (models.Medicacion, error)
			want       int
		}{
			{"exito", string(validBody), exito, http.StatusCreated},
			{"json_invalido", "{invalido", nil, http.StatusBadRequest},
			{"error_servicio", string(validBody), func(models.Medicacion) (models.Medicacion, error) {
				return models.Medicacion{}, errors.New("db")
			}, http.StatusInternalServerError},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{Medicacion: svc.NewMedicacionService(&mockMedRepo{crear: c.fn})}
				check(t, run(s.CrearMedicacion, http.MethodPost, "/medicaciones", c.body, "", ""), c.want)
			})
		}

		// campos requeridos: cada mutación debe devolver 400
		mutaciones := map[string]func(*models.Medicacion){
			"nombre":              func(m *models.Medicacion) { m.Nombre = "" },
			"dosis":               func(m *models.Medicacion) { m.Dosis = "" },
			"frecuencia":          func(m *models.Medicacion) { m.Frecuencia = "" },
			"hora_programada":     func(m *models.Medicacion) { m.Hora_programada = "" },
			"inicio_tratamiento":  func(m *models.Medicacion) { m.Inicio_tratamiento = time.Time{} },
		}
		for campo, mutar := range mutaciones {
			t.Run("falta_"+campo, func(t *testing.T) {
				m := medValida()
				mutar(&m)
				body, _ := json.Marshal(m)
				s := &Server{Medicacion: svc.NewMedicacionService(&mockMedRepo{})}
				check(t, run(s.CrearMedicacion, http.MethodPost, "/medicaciones", string(body), "", ""), http.StatusBadRequest)
			})
		}
	})

	t.Run("Actualizar", func(t *testing.T) {
		validBody, _ := json.Marshal(medValida())
		exito := func(id int, m models.Medicacion) (models.Medicacion, error) { return m, nil }

		for _, c := range []struct {
			name, id, body string
			fn             func(int, models.Medicacion) (models.Medicacion, error)
			want           int
		}{
			{"exito", "1", string(validBody), exito, http.StatusOK},
			{"id_invalido", "abc", "{}", nil, http.StatusBadRequest},
			{"json_invalido", "1", "{invalido", nil, http.StatusBadRequest},
			{"no_encontrada", "99", string(validBody), func(int, models.Medicacion) (models.Medicacion, error) {
				return models.Medicacion{}, gorm.ErrRecordNotFound
			}, http.StatusNotFound},
			{"error_servicio", "1", string(validBody), func(int, models.Medicacion) (models.Medicacion, error) {
				return models.Medicacion{}, errors.New("db")
			}, http.StatusInternalServerError},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{Medicacion: svc.NewMedicacionService(&mockMedRepo{actualizar: c.fn})}
				check(t, run(s.ActualizarMedicacion, http.MethodPut, "/medicaciones/"+c.id, c.body, "id", c.id), c.want)
			})
		}
	})

	t.Run("Eliminar", func(t *testing.T) {
		for _, c := range []struct {
			name, id string
			fn       func(int) error
			want     int
		}{
			{"exito", "1", func(int) error { return nil }, http.StatusOK},
			{"id_invalido", "abc", nil, http.StatusBadRequest},
			{"error_servicio", "1", func(int) error { return errors.New("db") }, http.StatusInternalServerError},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{Medicacion: svc.NewMedicacionService(&mockMedRepo{eliminar: c.fn})}
				check(t, run(s.EliminarMedicacion, http.MethodDelete, "/medicaciones/"+c.id, "", "id", c.id), c.want)
			})
		}
	})
}

// ---------- Historial ----------

func TestHistorialHandlers(t *testing.T) {
	histValida := func() models.HistorialMedicacion { return models.HistorialMedicacion{} }
	validBody, _ := json.Marshal(histValida())

	t.Run("Listar", func(t *testing.T) {
		for _, c := range []struct {
			name string
			fn   func() ([]models.HistorialMedicacion, error)
			want int
		}{
			{"exito", func() ([]models.HistorialMedicacion, error) { return []models.HistorialMedicacion{histValida()}, nil }, http.StatusOK},
			{"error", func() ([]models.HistorialMedicacion, error) { return nil, errors.New("db") }, http.StatusInternalServerError},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{Historial: svc.NewHistorialService(&mockHistRepo{listar: c.fn})}
				check(t, run(s.ListarHistorial, http.MethodGet, "/historial", "", "", ""), c.want)
			})
		}
	})

	t.Run("BuscarPorID", func(t *testing.T) {
		for _, c := range []struct {
			name, id string
			fn       func(int) (models.HistorialMedicacion, error)
			want     int
		}{
			{"exito", "1", func(int) (models.HistorialMedicacion, error) { return histValida(), nil }, http.StatusOK},
			{"id_invalido", "abc", nil, http.StatusBadRequest},
			{"no_encontrado", "99", func(int) (models.HistorialMedicacion, error) {
				return models.HistorialMedicacion{}, gorm.ErrRecordNotFound
			}, http.StatusNotFound},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{Historial: svc.NewHistorialService(&mockHistRepo{obtener: c.fn})}
				check(t, run(s.BuscarPorID, http.MethodGet, "/historial/"+c.id, "", "id", c.id), c.want)
			})
		}
	})

	t.Run("Crear", func(t *testing.T) {
		for _, c := range []struct {
			name, body string
			fn         func(models.HistorialMedicacion) (models.HistorialMedicacion, error)
			want       int
		}{
			{"exito", string(validBody), func(h models.HistorialMedicacion) (models.HistorialMedicacion, error) { return h, nil }, http.StatusCreated},
			{"json_invalido", "{invalido", nil, http.StatusBadRequest},
			{"error_servicio", string(validBody), func(models.HistorialMedicacion) (models.HistorialMedicacion, error) {
				return models.HistorialMedicacion{}, errors.New("db")
			}, http.StatusInternalServerError},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{Historial: svc.NewHistorialService(&mockHistRepo{crear: c.fn})}
				check(t, run(s.CrearHistorial, http.MethodPost, "/historial", c.body, "", ""), c.want)
			})
		}
	})

	t.Run("Actualizar", func(t *testing.T) {
		for _, c := range []struct {
			name, id, body string
			fn             func(int, models.HistorialMedicacion) (models.HistorialMedicacion, error)
			want           int
		}{
			{"exito", "1", string(validBody), func(id int, h models.HistorialMedicacion) (models.HistorialMedicacion, error) { return h, nil }, http.StatusOK},
			{"id_invalido", "abc", "{}", nil, http.StatusBadRequest},
			{"json_invalido", "1", "{invalido", nil, http.StatusBadRequest},
			{"no_encontrado", "99", string(validBody), func(int, models.HistorialMedicacion) (models.HistorialMedicacion, error) {
				return models.HistorialMedicacion{}, gorm.ErrRecordNotFound
			}, http.StatusNotFound},
			{"error_servicio", "1", string(validBody), func(int, models.HistorialMedicacion) (models.HistorialMedicacion, error) {
				return models.HistorialMedicacion{}, errors.New("db")
			}, http.StatusInternalServerError},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{Historial: svc.NewHistorialService(&mockHistRepo{actualizar: c.fn})}
				check(t, run(s.ActualizarHistorial, http.MethodPut, "/historial/"+c.id, c.body, "id", c.id), c.want)
			})
		}
	})

	t.Run("Eliminar", func(t *testing.T) {
		for _, c := range []struct {
			name, id string
			fn       func(int) error
			want     int
		}{
			{"exito", "1", func(int) error { return nil }, http.StatusNoContent},
			{"id_invalido", "abc", nil, http.StatusBadRequest},
			{"error_servicio", "1", func(int) error { return errors.New("db") }, http.StatusInternalServerError},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{Historial: svc.NewHistorialService(&mockHistRepo{eliminar: c.fn})}
				check(t, run(s.EliminarHistorial, http.MethodDelete, "/historial/"+c.id, "", "id", c.id), c.want)
			})
		}
	})
}

// ---------- Consultas cruzadas por paciente ----------

func TestMedicacionHistorialHandlers(t *testing.T) {
	t.Run("ListarMedicacionPorPaciente", func(t *testing.T) {
		for _, c := range []struct {
			name string
			fn   func(int) ([]models.Medicacion, error)
			want int
		}{
			{"exito", func(int) ([]models.Medicacion, error) { return []models.Medicacion{medValida()}, nil }, http.StatusOK},
			{"error", func(int) ([]models.Medicacion, error) { return nil, errors.New("db") }, http.StatusInternalServerError},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{MedicacionHistorial: svc.NewMedicacionHistorialService(&mockMedHistRepo{medPorPaciente: c.fn})}
				check(t, run(s.ListarMedicacionPorPaciente, http.MethodGet, "/pacientes/5/medicaciones", "", "id", "5"), c.want)
			})
		}
	})

	t.Run("ListarHistorialPorPaciente", func(t *testing.T) {
		for _, c := range []struct {
			name string
			fn   func(int) ([]models.HistorialMedicacion, error)
			want int
		}{
			{"exito", func(int) ([]models.HistorialMedicacion, error) { return []models.HistorialMedicacion{{}}, nil }, http.StatusOK},
			{"error", func(int) ([]models.HistorialMedicacion, error) { return nil, errors.New("db") }, http.StatusInternalServerError},
		} {
			t.Run(c.name, func(t *testing.T) {
				s := &Server{MedicacionHistorial: svc.NewMedicacionHistorialService(&mockMedHistRepo{histPorPaciente: c.fn})}
				check(t, run(s.ListarHistorialPorPaciente, http.MethodGet, "/pacientes/5/historial", "", "id", "5"), c.want)
			})
		}
	})
}