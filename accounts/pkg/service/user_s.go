package service

import (
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/repository"
	"gorm.io/gorm"
)

type UserService interface {
	Create(*model.User) (*model.User, error)
	Save(*model.User) (*model.User, error)
	AddToRole(*model.User, *model.Role) error
	RemFromRole(*model.User, *model.Role) error
	WithTrx(*gorm.DB) UserService
	FindById(id uint) *model.User
	FindByEmail(email string) *model.User
	FindAll() []*model.User
	Ban(user *model.User) error
	UnBan(user *model.User) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return userService{
		userRepository: r,
	}
}

func (u userService) Create(user *model.User) (*model.User, error) {
	return u.userRepository.Create(user)
}

func (u userService) Save(user *model.User) (*model.User, error) {
	return u.userRepository.Save(user)
}

func (u userService) AddToRole(user *model.User, role *model.Role) error {
	return u.userRepository.AddToRole(user, role)
}

func (u userService) RemFromRole(user *model.User, role *model.Role) error {
	return u.userRepository.RemFromRole(user, role)
}

func (u userService) WithTrx(trx *gorm.DB) UserService {
	u.userRepository = u.userRepository.WithTrx(trx)
	return u
}

func (u userService) FindById(id uint) *model.User {
	return u.userRepository.FindById(id)
}

func (u userService) FindByEmail(email string) *model.User {
	return u.userRepository.FindByEmail(email)
}

func (u userService) FindAll() []*model.User {
	return u.userRepository.All()
}

func (u userService) Ban(user *model.User) error {
	return u.userRepository.Ban(user)

}

func (u userService) UnBan(user *model.User) error {
	return u.userRepository.UnBan(user)
}
