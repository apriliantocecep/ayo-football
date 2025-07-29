package usecase

import (
	"context"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/config"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/entity"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/model"
	"github.com/apriliantocecep/ayo-football/services/auth/internal/repository"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserUseCase struct {
	DB             *gorm.DB
	UserRepository *repository.UserRepository
	Jwt            *config.JwtWrapper
}

func (u *UserUseCase) hashPassword(parentCtx context.Context, ctx context.Context, password string) (string, error) {
	var passwordResultChan = make(chan utils.PasswordResult)
	go utils.HashPasswordAsync(password, passwordResultChan)

	for {
		select {
		case <-ctx.Done():
			return "", status.Errorf(codes.DeadlineExceeded, "context timeout")
		case pass := <-passwordResultChan:
			if pass.Err != nil {
				return "", status.Errorf(codes.Internal, pass.Err.Error())
			}
			return pass.Password, nil
		}
	}
}

func (u *UserUseCase) Register(ctx context.Context, request *model.RegisterRequest) (*model.RegisterResponse, error) {
	userEntity, err := new(entity.User), *new(error)

	hashCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	password, err := u.hashPassword(ctx, hashCtx, request.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hashing password error: %v", err)
	}
	username := strings.Split(request.Email, "@")[0]

	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	userEntity, err = u.UserRepository.FindByEmailOrUsername(ctx, tx, request.Email, username)
	if err == nil {
		if userEntity.Email == request.Email {
			return nil, status.Errorf(codes.AlreadyExists, "email already exists")
		}

		if userEntity.Username == username {
			return nil, status.Errorf(codes.AlreadyExists, "username already exists")
		}
	}

	user := entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Username: username,
		Password: password,
		Role:     "admin",
	}

	userUuid, err := u.UserRepository.Create(ctx, tx, &user)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create user")
	}

	if err = tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create user")
	}

	response := model.RegisterResponse{UserId: userUuid.String(), Username: username}

	return &response, nil
}

func (u *UserUseCase) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error) {
	userEntity, err := new(entity.User), *new(error)

	if utils.ValidateEmail(request.Identity) {
		userEntity, err = u.UserRepository.FindByEmail(ctx, u.DB, request.Identity)
	} else {
		userEntity, err = u.UserRepository.FindByUsername(ctx, u.DB, request.Identity)
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
	return &UserUseCase{
		UserRepository: userRepository,
		Jwt:            jwt,
		DB:             db,
	}
}
