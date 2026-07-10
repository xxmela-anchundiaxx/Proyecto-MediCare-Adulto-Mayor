package farmacia

import (
	"gorm.io/gorm"
	farmaciaModel "proyecto-medicare-adulto-mayor/internal/models/farmacia"
)

// Este es el struct que necesita main.go
type StorageFarmaciaGORM struct {
	DB *gorm.DB
}

// Este es el constructor exacto que llama tu main.go
func NuevoStorageFarmaciaGORM(db *gorm.DB) *StorageFarmaciaGORM {
	return &StorageFarmaciaGORM{DB: db}
}

// 1. Crear
func (s *StorageFarmaciaGORM) CrearFarmacia(f *farmaciaModel.Farmacia) error {
	return s.DB.Create(f).Error
}

// 2. Buscar por ID
func (s *StorageFarmaciaGORM) BuscarPorID(id string) (*farmaciaModel.Farmacia, error) {
	var f farmaciaModel.Farmacia
	if err := s.DB.First(&f, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

// 3. Listar Todas
func (s *StorageFarmaciaGORM) ListarTodas() ([]farmaciaModel.Farmacia, error) {
	var lista []farmaciaModel.Farmacia
	err := s.DB.Find(&lista).Error
	return lista, err
}

// 4. Buscar Cercanas
func (s *StorageFarmaciaGORM) BuscarCercanas(lat, lon, radioKM float64) ([]farmaciaModel.Farmacia, error) {
	var lista []farmaciaModel.Farmacia
	grados := radioKM / 111.0
	err := s.DB.Where("latitud BETWEEN ? AND ? AND longitud BETWEEN ? AND ?", 
		lat-grados, lat+grados, lon-grados, lon+grados).Find(&lista).Error
	return lista, err
}

// 5. Actualizar Farmacia (¡El que requería la interfaz!)
func (s *StorageFarmaciaGORM) ActualizarFarmacia(id string, f *farmaciaModel.Farmacia) error {
	return s.DB.Model(&farmaciaModel.Farmacia{}).Where("id = ?", id).Updates(f).Error
}

// 6. Eliminar Farmacia
func (s *StorageFarmaciaGORM) EliminarFarmacia(id string) error {
	return s.DB.Delete(&farmaciaModel.Farmacia{}, "id = ?", id).Error
}