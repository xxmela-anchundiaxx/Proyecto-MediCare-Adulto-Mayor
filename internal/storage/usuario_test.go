package storage

import (
	"proyecto-medicare-adulto-mayor/internal/models"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestAlmacenUsuario(t *testing.T) {
	// 1. Configurar DB en memoria para el test
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("No se pudo abrir la base de datos en memoria: %v", err)
	}

	// Migrar el modelo
	db.AutoMigrate(&models.Usuario{})

	// Inicializar el almacen
	almacen := NuevoAlmacenUsuario(db)

	// 2. Definir un usuario de prueba
	nuevoUsuario := models.Usuario{
		Email:        "test@usuario.com",
		PasswordHash: "hashed_password",
		CreadoEn:     time.Now(),
	}

	// 3. Probar la creación
	t.Run("CrearUsuario", func(t *testing.T) {
		resultado, err := almacen.CrearUsuario(nuevoUsuario)
		if err != nil {
			t.Errorf("Error al crear usuario: %v", err)
		}
		if resultado.Email != nuevoUsuario.Email {
			t.Errorf("Esperaba email %s, got %s", nuevoUsuario.Email, resultado.Email)
		}
	})

	// 4. Probar la búsqueda
	t.Run("BuscarUsuarioPorEmail", func(t *testing.T) {
		usuario, encontrado := almacen.BuscarUsuarioPorEmail("test@usuario.com")
		if !encontrado {
			t.Error("Debería haber encontrado al usuario")
		}
		if usuario.Email != "test@usuario.com" {
			t.Errorf("Email no coincide: %s", usuario.Email)
		}
	})

	// 5. Probar búsqueda de usuario inexistente
	t.Run("UsuarioNoExiste", func(t *testing.T) {
		_, encontrado := almacen.BuscarUsuarioPorEmail("noexiste@correo.com")
		if encontrado {
			t.Error("No debería haber encontrado un usuario que no existe")
		}
	})
}