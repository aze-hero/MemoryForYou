package repository

import (
	"github.com/zhaozeguang/timecapsule-api/internal/model"
	"github.com/zhaozeguang/timecapsule-api/pkg/database"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Create(user *model.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepo) FindByID(id string) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByProvider(provider, providerID string) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, "provider = ? AND provider_id = ?", provider, providerID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
