package usecase

import (
	"context"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/entity"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/model"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type RegisterUseCase struct {
	DB             *gorm.DB
	UserRepository *repository.UserRepository
}

func (u *RegisterUseCase) CreateUser(ctx context.Context, request *model.CreateUserRequest) (*model.CreateUserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Username: request.Username,
		Password: request.Password,
	}

	userUuid, err := u.UserRepository.Create(ctx, tx, &user)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create user")
	}

	if err = tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create user")
	}

	response := model.CreateUserResponse{UserId: userUuid.String()}
	return &response, nil
}

func NewRegisterUseCase(DB *gorm.DB, userRepository *repository.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{DB: DB, UserRepository: userRepository}
}
