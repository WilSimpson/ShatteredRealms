package repository

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"gorm.io/gorm"
)

type CharacterRepository interface {
	Create(character *model.Character) (*model.Character, error)
	Save(character *model.Character) (*model.Character, error)
	Delete(character *model.Character) error
	FindById(id uint64) (*model.Character, error)
	FindAll() ([]*model.Character, error)

	FindAllByOwner(id uint64) ([]*model.Character, error)

	WithTrx(trx *gorm.DB) CharacterRepository
	Migrate() error
}

type characterRepository struct {
	DB *gorm.DB
}

func NewCharacterRepository(db *gorm.DB) CharacterRepository {
	return characterRepository{
		DB: db,
	}
}

func (r characterRepository) Create(character *model.Character) (*model.Character, error) {
	// Set the ID to zero so the database can generate the value
	character.ID = 0

	err := r.DB.Create(&character).Error
	if err != nil {
		return nil, err
	}

	return character, nil
}

func (r characterRepository) Save(character *model.Character) (*model.Character, error) {
	err := r.DB.Save(&character).Error
	if err != nil {
		return nil, err
	}

	return character, nil
}

func (r characterRepository) Delete(character *model.Character) error {
	return r.DB.Delete(&character).Error
}

func (r characterRepository) FindById(id uint64) (*model.Character, error) {
	var character *model.Character = nil
	err := r.DB.First(&character, id).Error
	if err != nil {
		return nil, err
	}

	return character, nil
}

func (r characterRepository) FindAll() ([]*model.Character, error) {
	var characters []*model.Character
	return characters, r.DB.Find(&characters).Error
}

func (r characterRepository) FindAllByOwner(id uint64) ([]*model.Character, error) {
	var characters []*model.Character
	return characters, r.DB.Where("owner_id = ?", id).Find(&characters).Error
}

func (r characterRepository) WithTrx(trx *gorm.DB) CharacterRepository {
	if trx == nil {
		return r
	}

	r.DB = trx
	return r
}

func (r characterRepository) Migrate() error {
	return r.DB.AutoMigrate(&model.Character{})
}
