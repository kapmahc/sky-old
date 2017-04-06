package orm

import (
	"github.com/jinzhu/gorm"
	"github.com/kapmahc/sky/settings"
)

// New new gorm store
func New(db *gorm.DB) settings.Store {
	db.AutoMigrate(&Model{})
	return &Store{
		db: db,
	}
}

// Store gorm store
type Store struct {
	db *gorm.DB
}

// Set save
func (p *Store) Set(key string, val []byte, enc bool) error {
	var m Model
	if p.db.Where("key = ?", key).First(&m).RecordNotFound() {
		return p.db.Create(&Model{
			Key:    key,
			Val:    val,
			Encode: enc,
		}).Error
	}

	return p.db.Model(&m).Updates(map[string]interface{}{
		"encode": enc,
		"val":    val,
	}).Error

}

// Get get
func (p *Store) Get(key string) ([]byte, bool, error) {
	var m Model
	err := p.db.Where("key = ?", key).First(&m).Error
	if err != nil {
		return nil, false, err
	}
	return m.Val, m.Encode, nil
}
