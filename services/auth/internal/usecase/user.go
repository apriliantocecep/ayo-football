package usecase

import (
	"context"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/config"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/entity"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/repository"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
)

type UserUseCase struct {
	DB             *gorm.DB
	UserRepository *repository.UserRepository
	Jwt            *config.JwtWrapper
}

func (u *UserUseCase) Register(ctx context.Context, request *model.RegisterRequest) (*model.RegisterResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if _, err := u.UserRepository.FindByEmail(tx, request.Email); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "email already exists")
	}

	// get username from email
	username := strings.Split(request.Email, "@")

	if _, err := u.UserRepository.FindByUsername(tx, username[0]); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	}

	user := entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Username: username[0],
		Password: utils.HashPassword(request.Password),
	}

	userUuid, err := u.UserRepository.Create(tx, &user)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create user")
	}

	if err = tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create user")
	}

	response := model.RegisterResponse{UserId: userUuid.String(), Username: username[0]}

	return &response, nil
}

func (u *UserUseCase) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error) {
	userEntity, err := new(entity.User), *new(error)

	if utils.ValidateEmail(request.Identity) {
		userEntity, err = u.UserRepository.FindByEmail(u.DB, request.Identity)
	} else {
		userEntity, err = u.UserRepository.FindByUsername(u.DB, request.Identity)
	}

	if err != nil {
		return nil, err
	}

	match := utils.ComparePasswordHash(request.Password, userEntity.Password)
	if !match {
		return nil, status.Errorf(codes.InvalidArgument, "invalid identity or password")
	}

	accessToken, accessClaim, err := u.Jwt.GenerateAccessToken(userEntity)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create access token")
	}

	response := &model.LoginResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaim.ExpiresAt.Time,
	}

	return response, nil
}

func (u *UserUseCase) GetUserById(ctx context.Context, id string) (*model.UserResource, error) {
	userEntity, err := u.UserRepository.FindById(u.DB, id)
	if err != nil {
		return nil, err
	}

	response := model.UserResource{
		ID:        userEntity.ID.String(),
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		Username:  userEntity.Username,
		CreatedAt: userEntity.CreatedAt,
	}

	return &response, nil
}

func (u *UserUseCase) ValidateToken(ctx context.Context, token string) (string, error) {
	claims, err := u.Jwt.ValidateToken(token)
	if err != nil {
		return "", status.Errorf(codes.Aborted, "can not validated access token")
	}

	return claims.RegisteredClaims.Subject, nil
}

func NewUserUseCase(userRepository *repository.UserRepository, jwt *config.JwtWrapper, db *gorm.DB) *UserUseCase {
	return &UserUseCase{UserRepository: userRepository, Jwt: jwt, DB: db}
}
