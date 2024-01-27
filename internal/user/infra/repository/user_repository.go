package repository

import (
	"p2p/internal/user/entity"

	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepositoryPostgres struct {
	DB *gorm.DB
}

func NewUserRepositoryPostgres(db *gorm.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{DB: db}
}

func (r *UserRepositoryPostgres) Create(user *entity.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepositoryPostgres) FindAll() ([]*entity.User, error) {
	var users []*entity.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryPostgres) Update(user *entity.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepositoryPostgres) DeleteByID(id uint) error {
	return r.DB.Where("id = ?", id).Delete(entity.User{}).Error
}

func (r *UserRepositoryPostgres) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
