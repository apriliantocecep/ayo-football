package usecase

import (
	"context"
	"errors"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/config"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/entity"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/gateway/messaging"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/repository"
	sharedmodel "github.com/apriliantocecep/posfin-blog/shared/model"
	"github.com/apriliantocecep/posfin-blog/shared/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log"
	"strings"
)

type UserUseCase struct {
	DB                   *gorm.DB
	UserRepository       *repository.UserRepository
	Jwt                  *config.JwtWrapper
	UserCreatedPublisher *messaging.UserPublisher
}

func (u *UserUseCase) Register(ctx context.Context, request *model.RegisterRequest) (*model.RegisterResponse, error) {
	password := utils.HashPassword(request.Password)

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
		Password: password,
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

func (u *UserUseCase) RegisterWithQueue(ctx context.Context, request *model.RegisterRequest) (*model.RegisterResponse, error) {
	password := utils.HashPassword(request.Password)

	// get username from email
	username := strings.Split(request.Email, "@")[0]

	existingUser := new(entity.User)
	err := u.UserRepository.FindByEmailOrUsername(u.DB, request.Email, username, existingUser)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.AlreadyExists, "email or username already exists")
	}

	if existingUser.Email == request.Email {
		return nil, status.Errorf(codes.AlreadyExists, "email already exists")
	}

	if existingUser.Username == username {
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	}

	// publish to broker
	event := sharedmodel.UserEvent{
		Name:     request.Name,
		Email:    request.Email,
		Username: username,
		Password: password,
	}
	err = u.UserCreatedPublisher.Publish(&event)
	if err != nil {
		log.Printf("failed publish user created event : %+v", err)
		return nil, status.Errorf(codes.Aborted, "failed to publish user data")
	}

	response := model.RegisterResponse{UserId: "<known after queued>", Username: username}
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

func NewUserUseCase(userRepository *repository.UserRepository, jwt *config.JwtWrapper, db *gorm.DB, userCreatedPublisher *messaging.UserPublisher) *UserUseCase {
	return &UserUseCase{
		UserRepository:       userRepository,
		Jwt:                  jwt,
		DB:                   db,
		UserCreatedPublisher: userCreatedPublisher,
	}
}
