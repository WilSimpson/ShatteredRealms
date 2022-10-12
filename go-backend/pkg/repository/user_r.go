package repository

import (
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

type UserRepository interface {
	Create(*model.User) (*model.User, error)
	Save(*model.User) (*model.User, error)
	AddToRole(*model.User, *model.Role) error
	RemFromRole(*model.User, *model.Role) error
	WithTrx(*gorm.DB) UserRepository
	FindById(id uint) *model.User
	FindByEmail(email string) *model.User
	FindByUsername(username string) *model.User
	Migrate() error
	All() []*model.User
	Ban(user *model.User) error
	UnBan(user *model.User) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return userRepository{
		DB: db,
	}
}

func (u userRepository) Create(user *model.User) (*model.User, error) {
	err := user.Validate()
	if err != nil {
		return user, err
	}

	conflict := u.FindByEmail(user.Email)
	if conflict.Exists() {
		return user, fmt.Errorf("email is already taken")
	}

	conflict = u.FindByUsername(user.Username)
	if conflict.Exists() {
		return user, fmt.Errorf("username is already taken")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
	if err != nil {
		return user, fmt.Errorf("password: %w", err)
	}

	user.Password = string(hashedPass)
	err = u.DB.Create(&user).Error

	return user, err
}

func (u userRepository) Save(user *model.User) (*model.User, error) {
	conflict := u.FindByEmail(user.Email)
	if conflict.Exists() && user.ID != conflict.ID {
		return user, fmt.Errorf("email is already taken")
	}

	conflict = u.FindByUsername(user.Username)
	if conflict.Exists() && user.ID != conflict.ID {
		return user, fmt.Errorf("username is already taken")
	}

	return user, u.DB.Save(&user).Error
}

func (u userRepository) AddToRole(user *model.User, role *model.Role) error {
	return u.DB.Model(&user).Association("Roles").Append([]model.Role{*role})
}

func (u userRepository) RemFromRole(user *model.User, role *model.Role) error {
	return u.DB.Model(&user).Association("Roles").Delete([]model.Role{*role})
}

func (u userRepository) WithTrx(trx *gorm.DB) UserRepository {
	if trx == nil {
		return u
	}

	u.DB = trx
	return u
}

func (u userRepository) FindById(id uint) *model.User {
	var user *model.User
	u.DB.Where("id=?", id).Preload("Roles").Find(&user)
	return user
}

func (u userRepository) FindByEmail(email string) *model.User {
	var user *model.User
	u.DB.Where("email=?", email).Preload("Roles").Find(&user)
	return user
}

func (u userRepository) FindByUsername(username string) *model.User {
	var user *model.User
	u.DB.Where("username=?", username).Preload("Roles").Find(&user)
	return user
}

func (u userRepository) Migrate() error {
	return u.DB.AutoMigrate(&model.User{})
}

func (u userRepository) All() []*model.User {
	var users []*model.User
	u.DB.Preload("Roles").Find(&users)
	return users
}

func (u userRepository) Ban(user *model.User) error {
	user.BannedAt = time.Now()
	return u.DB.Save(&user).Error
}

func (u userRepository) UnBan(user *model.User) error {
	user.BannedAt = time.Time{}
	return u.DB.Save(&user).Error
}
