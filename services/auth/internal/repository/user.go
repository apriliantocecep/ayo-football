package repository

import (
	"errors"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/entity"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserRepository struct {
}

func (u *UserRepository) Create(db *gorm.DB, user *entity.User) (uuid.UUID, error) {
	result := db.Create(&user)
	return user.ID, result.Error
}

func (u *UserRepository) FindById(db *gorm.DB, id string) (*entity.User, error) {
	var user entity.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) FindByEmail(db *gorm.DB, email string) (*entity.User, error) {
	var user entity.User
	if err := db.Where(&entity.User{Email: email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "identity not found")
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) FindByUsername(db *gorm.DB, username string) (*entity.User, error) {
	var user entity.User
	if err := db.Where(&entity.User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "identity not found")
		}
		return nil, err
	}
	return &user, nil
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
