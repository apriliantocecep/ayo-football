package repository

import (
	"context"
	"errors"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/entity"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserRepository struct {
}

func (u *UserRepository) Create(ctx context.Context, db *gorm.DB, user *entity.User) (uuid.UUID, error) {
	_, span := otel.Tracer("UserRepository").Start(ctx, "UserRepository.Create")
	defer span.End()

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

func (u *UserRepository) FindByEmail(ctx context.Context, db *gorm.DB, email string) (*entity.User, error) {
	_, span := otel.Tracer("UserRepository").Start(ctx, "UserRepository.FindByEmail")
	defer span.End()

	var user entity.User
	if err := db.Where(&entity.User{Email: email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "identity not found")
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) FindByUsername(ctx context.Context, db *gorm.DB, username string) (*entity.User, error) {
	_, span := otel.Tracer("UserRepository").Start(ctx, "UserRepository.FindByUsername")
	defer span.End()

	var user entity.User
	if err := db.Where(&entity.User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "identity not found")
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) FindByEmailOrUsername(ctx context.Context, db *gorm.DB, email, username string, user *entity.User) error {
	_, span := otel.Tracer("UserRepository").Start(ctx, "UserRepository.FindByEmailOrUsername")
	defer span.End()

	return db.Where("email = ? OR username = ?", email, username).First(user).Error
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
